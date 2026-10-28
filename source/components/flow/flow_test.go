package flow

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components/gesture"
	"github.com/oneo-me/onevra/source/theme"
)

func sample() Data {
	card := func(title string) func(*ui.Context) {
		return func(c *ui.Context) { Card(c, CardOptions{Title: title, Description: title + " details"}) }
	}
	return Data{
		Nodes: []Node{
			{ID: "plan", Label: "Plan", Content: card("Plan")},
			{ID: "build", Label: "Build", Content: card("Build")},
			{ID: "test", Label: "Test", Content: card("Test")},
			{ID: "ship", Label: "Ship", Muted: true, Content: card("Ship")},
			{ID: "fixed", Label: "Fixed", Width: 90, Height: 40},
		},
		Edges: []Edge{
			{Source: "plan", Target: "build"},
			{Source: "plan", Target: "test", Dashed: true, Label: "depends on"},
			{Source: "build", Target: "ship"},
			{Source: "test", Target: "ship"},
			{Source: "ship", Target: "ship"},
			{Source: "ship", Target: "missing"},
		},
	}
}

type harness struct {
	tt               *ui.Tester
	state            State
	data             Data
	opts             Options
	selected, opened []string
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	h := &harness{data: sample()}
	h.opts = Options{
		Label:  "Pipeline",
		Select: func(id string) { h.selected = append(h.selected, id); h.opts.Selected = id },
		Open:   func(id string) { h.opened = append(h.opened, id) },
	}
	h.tt = ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		View(c, &h.state, h.data, h.opts)
	}, 900, 700)
	h.tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	h.settle()
	if !h.state.Ready() {
		t.Fatal("diagram did not become ready")
	}
	return h
}

// settle builds a frame for changed data and lets measurement and motion
// finish.
func (h *harness) settle() {
	h.tt.Frame()
	for range 4 {
		h.tt.Advance(20 * time.Millisecond)
	}
}

func TestViewMeasuresAndLaysOut(t *testing.T) {
	h := newHarness(t)
	v := &h.state.view
	if len(v.nodes) != 5 || len(v.edges) != 4 {
		t.Fatalf("view kept %d nodes and %d edges", len(v.nodes), len(v.edges))
	}
	for _, n := range v.nodes {
		if !n.measured {
			t.Fatalf("node %s is unmeasured", n.id)
		}
	}
	if n := v.nodes[v.index["plan"]]; n.width != 240 || n.height <= 0 {
		t.Fatalf("card measured %vx%v", n.width, n.height)
	}
	if n := v.nodes[v.index["fixed"]]; n.width != 90 || n.height != 40 {
		t.Fatalf("fixed node measured %vx%v", n.width, n.height)
	}
	_, planY, _ := h.state.Position("plan")
	_, buildY, _ := h.state.Position("build")
	_, shipY, _ := h.state.Position("ship")
	if !(planY < buildY && buildY < shipY) {
		t.Fatalf("ranks at %v, %v, %v", planY, buildY, shipY)
	}
	if _, _, ok := h.state.Position("missing"); ok {
		t.Fatal("unknown node is laid out")
	}
	// The fitted diagram shows whole, at no more than its natural size.
	if z := h.state.Zoom(); z <= 0 || z > 1 {
		t.Fatalf("fitted zoom %v", z)
	}
	r, ok := h.tt.Find("Plan")
	if !ok || r.X < 0 || r.Y < 0 || r.X+r.W > 900 || r.Y+r.H > 700 {
		t.Fatalf("plan shows at %v", r)
	}
	if _, ok := h.tt.Find("Pipeline"); !ok {
		t.Fatal("diagram has no accessible label")
	}
}

func TestViewSelectsAndOpens(t *testing.T) {
	h := newHarness(t)
	if err := h.tt.Click("Build"); err != nil {
		t.Fatal(err)
	}
	h.settle()
	if len(h.selected) != 1 || h.selected[0] != "build" {
		t.Fatalf("selected %v", h.selected)
	}
	if !h.state.emphasized(h.state.edges[0]) || h.state.emphasized(h.state.edges[1]) {
		t.Fatal("selection does not emphasize its edges alone")
	}
	r, _ := h.tt.Find("Test")
	h.tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	h.tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
	h.settle()
	if len(h.opened) != 1 || h.opened[0] != "test" {
		t.Fatalf("opened %v", h.opened)
	}
	// Clicking empty space clears the selection without opening anything.
	h.tt.ClickAt(4, 690)
	h.settle()
	if last := h.selected[len(h.selected)-1]; last != "" || len(h.opened) != 1 {
		t.Fatalf("empty click selected %q and opened %v", last, h.opened)
	}
}

func TestViewPansAndZooms(t *testing.T) {
	h := newHarness(t)
	v := &h.state.view
	before := v.shown
	h.tt.Press(4, 690)
	h.tt.Move(40, 650)
	h.tt.Move(64, 640)
	h.tt.Release(64, 640)
	h.settle()
	if v.shown.x-before.x != 60 || v.shown.y-before.y != -50 {
		t.Fatalf("pan moved the camera from %v to %v", before, v.shown)
	}
	if len(h.selected) != 0 {
		t.Fatal("panning selected")
	}
	zoom := h.state.Zoom()
	h.state.ZoomIn()
	h.settle()
	if h.state.Zoom() <= zoom {
		t.Fatal("zoom in did not zoom")
	}
	h.state.Fit()
	h.settle()
	if h.state.Zoom() != zoom || v.userMoved {
		t.Fatalf("fit zoom %v, want %v", h.state.Zoom(), zoom)
	}
}

func TestViewRelaysOutWhenDataChanges(t *testing.T) {
	h := newHarness(t)
	h.data.Nodes = append(h.data.Nodes, Node{ID: "docs", Label: "Docs", Content: func(c *ui.Context) { Card(c, CardOptions{Title: "Docs"}) }})
	h.data.Edges = append(h.data.Edges, Edge{Source: "ship", Target: "docs"})
	h.settle()
	_, shipY, _ := h.state.Position("ship")
	_, docsY, ok := h.state.Position("docs")
	if !ok || docsY <= shipY {
		t.Fatalf("new node at %v below %v", docsY, shipY)
	}
	if !h.state.view.nodes[h.state.index["docs"]].measured {
		t.Fatal("new node is unmeasured")
	}
	// Content that grows is measured and laid out again.
	h.data.Nodes[0].Content = func(c *ui.Context) { Card(c, CardOptions{Width: 320, Title: "Plan"}) }
	h.settle()
	if n := h.state.view.nodes[0]; n.width != 320 {
		t.Fatalf("resized card measured %v wide", n.width)
	}
	if x, _, _ := h.state.Position("plan"); h.state.view.result.width < float64(x)+320 {
		t.Fatal("layout did not follow the resized card")
	}
	h.data = Data{}
	h.opts.EmptyText = "Nothing planned"
	h.settle()
	if !h.tt.HasText("Nothing planned") {
		t.Fatal("empty state is missing")
	}
}

func TestMapViewPansFromNodesAndKeepsClicks(t *testing.T) {
	h := newHarness(t)
	h.opts.Navigation = NavigationMap
	h.settle()
	v := &h.state.view
	if v.navigation.ScrollMode() != gesture.ScrollZoom {
		t.Fatal("NavigationMap did not reach the navigation")
	}
	// A press-drag that starts on a node pans the map, and stays a drag.
	r, ok := h.tt.Find("Build")
	if !ok {
		t.Fatal("Build node missing")
	}
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	before := v.shown
	h.tt.Press(cx, cy)
	h.tt.Move(cx+30, cy+20)
	h.tt.Move(cx+60, cy+40)
	h.tt.Release(cx+60, cy+40)
	h.settle()
	if v.shown.x-before.x != 60 || v.shown.y-before.y != 40 {
		t.Fatalf("map pan from a node moved the camera from %v to %v", before, v.shown)
	}
	if len(h.selected) != 0 {
		t.Fatalf("a map pan selected %v", h.selected)
	}
	// Clicks and double-clicks on nodes keep their meaning: select, open,
	// and no zoom — the zoom-on-double-click belongs to empty space. The
	// pan moved the node on screen, so its center is recomputed.
	if err := h.tt.Click("Build"); err != nil {
		t.Fatal(err)
	}
	h.settle()
	if len(h.selected) != 1 || h.selected[0] != "build" {
		t.Fatalf("map mode click selected %v", h.selected)
	}
	r2, ok := h.tt.Find("Build")
	if !ok {
		t.Fatal("Build node missing after pan")
	}
	cx2, cy2 := r2.X+r2.W/2, r2.Y+r2.H/2
	zoom := h.state.Zoom()
	h.tt.ClickAt(cx2, cy2)
	h.tt.ClickAt(cx2, cy2)
	h.settle()
	if len(h.opened) != 1 || h.opened[0] != "build" {
		t.Fatalf("map mode double-click opened %v", h.opened)
	}
	if h.state.Zoom() != zoom {
		t.Fatalf("double-click on a node zoomed from %v", zoom)
	}
	// Clicking empty space still clears.
	h.tt.ClickAt(4, 690)
	h.settle()
	if last := h.selected[len(h.selected)-1]; last != "" {
		t.Fatalf("empty click selected %q", last)
	}
}

// TestOnZoomTakesZoomGestures settles the hook contract: with OnZoom set,
// unmodified zoom gestures (wheel and scroll reach navigate as
// gesture.Zoom) call the hook and never move the camera's scale, and the
// map-mode double-click on empty space routes to the hook instead of
// zoomAt. Pinches keep scaling; see TestPinchScalesUnderOnZoom.
func TestOnZoomTakesZoomGestures(t *testing.T) {
	h := newHarness(t)
	var got []bool
	h.opts.Navigation = NavigationMap
	h.opts.OnZoom = func(zoomIn bool) { got = append(got, zoomIn) }
	h.settle()
	v := &h.state.view
	if v.onZoomFn == nil {
		t.Fatal("OnZoom did not reach the view")
	}
	v.navigate(gesture.Event{Kind: gesture.Zoom, Scale: 1.25, X: 10, Y: 10})
	v.navigate(gesture.Event{Kind: gesture.Zoom, Scale: .8, X: 10, Y: 10})
	if len(got) != 2 || !got[0] || got[1] {
		t.Fatalf("OnZoom saw %v, want [true false]", got)
	}
	if v.shown.zoom != 1 || v.camera.zoom != 1 {
		t.Fatalf("camera zoomed to shown %v camera %v under OnZoom", v.shown.zoom, v.camera.zoom)
	}
	// Empty-space double-click routes to the hook (plain = deeper,
	// Shift = shallower) and still leaves the scale alone.
	h.tt.ClickAt(4, 690)
	h.tt.ClickAt(4, 690)
	h.settle()
	if len(got) != 3 || !got[2] {
		t.Fatalf("empty double-click reached OnZoom as %v, want a third true", got)
	}
	if h.state.Zoom() != 1 {
		t.Fatalf("empty double-click scaled the view to %v under OnZoom", h.state.Zoom())
	}
	// Node clicks keep selecting and opening under the hook.
	if err := h.tt.Click("Build"); err != nil {
		t.Fatal(err)
	}
	h.settle()
	if last := h.selected[len(h.selected)-1]; last != "build" {
		t.Fatalf("OnZoom mode click selected %v", h.selected)
	}
}

// TestPinchScalesUnderOnZoom settles the split: with OnZoom set, plain
// wheel and scroll reach the hook, but a pinch or modified scroll (which
// the gesture layer marks as a pinch) scales the camera at the pointer.
func TestPinchScalesUnderOnZoom(t *testing.T) {
	h := newHarness(t)
	var got []bool
	h.opts.Navigation = NavigationMap
	h.opts.OnZoom = func(zoomIn bool) { got = append(got, zoomIn) }
	h.settle()
	v := &h.state.view
	before := v.shown
	gx, gy := (300-before.x)/before.zoom, (300-before.y)/before.zoom
	v.navigate(gesture.Event{Kind: gesture.Zoom, Scale: 1.5, X: 300, Y: 300, Pinch: true})
	if len(got) != 0 {
		t.Fatalf("a pinch reached OnZoom: %v", got)
	}
	if math.Abs(v.shown.zoom-1.5*before.zoom) > 1e-6 {
		t.Fatalf("pinch scaled to %v, want %v", v.shown.zoom, 1.5*before.zoom)
	}
	if math.Abs(v.shown.x+gx*v.shown.zoom-300) > 1e-6 || math.Abs(v.shown.y+gy*v.shown.zoom-300) > 1e-6 {
		t.Fatalf("pinch did not keep the point under the pointer: %+v", v.shown)
	}
	// A plain scroll still reaches the hook, untouched camera.
	v.navigate(gesture.Event{Kind: gesture.Zoom, Scale: 2, X: 10, Y: 10})
	if len(got) != 1 || !got[0] {
		t.Fatalf("plain scroll reached OnZoom as %v, want [true]", got)
	}
	if math.Abs(v.shown.zoom-1.5*before.zoom) > 1e-6 {
		t.Fatalf("the hook scaled the camera to %v", v.shown.zoom)
	}
}

// TestOnZoomAtReceivesThePointer settles the hook precedence: OnZoomAt
// wins over OnZoom, receives the point of the gesture, and pinches still
// scale the camera.
func TestOnZoomAtReceivesThePointer(t *testing.T) {
	h := newHarness(t)
	calls, x, y, in := 0, float32(0), float32(0), false
	h.opts.Navigation = NavigationMap
	h.opts.OnZoom = func(bool) { t.Fatal("OnZoomAt should take the gestures") }
	h.opts.OnZoomAt = func(px, py float32, zoomIn bool) { calls++; x, y, in = px, py, zoomIn }
	h.settle()
	v := &h.state.view
	if v.onZoomAtFn == nil {
		t.Fatal("OnZoomAt did not reach the view")
	}
	v.navigate(gesture.Event{Kind: gesture.Zoom, Scale: 2, X: 30, Y: 40})
	if calls != 1 || x != 30 || y != 40 || !in {
		t.Fatalf("OnZoomAt saw %v calls at (%v,%v) zoomIn=%v", calls, x, y, in)
	}
	if h.state.Zoom() != 1 {
		t.Fatalf("hook mode scaled the view to %v", h.state.Zoom())
	}
	// The map-mode double-click on empty space passes its point too.
	h.tt.ClickAt(4, 690)
	h.tt.ClickAt(4, 690)
	h.settle()
	if calls != 2 || x != 4 || y != 690 {
		t.Fatalf("double-click reached OnZoomAt %v times at (%v,%v)", calls, x, y)
	}
	if h.state.Zoom() != 1 {
		t.Fatalf("double-click under OnZoomAt scaled the view to %v", h.state.Zoom())
	}
	// A pinch still scales.
	v.navigate(gesture.Event{Kind: gesture.Zoom, Scale: 1.5, X: 10, Y: 10, Pinch: true})
	if h.state.Zoom() <= 1 {
		t.Fatalf("pinch left the camera at %v under OnZoomAt", h.state.Zoom())
	}
}

func TestCameraRoundTrip(t *testing.T) {
	h := newHarness(t)
	h.settle()
	if _, _, z := h.state.Camera(); z != h.state.Zoom() {
		t.Fatal("Camera disagrees with Zoom")
	}
	h.state.SetCamera(120, 60, 1.5)
	h.settle()
	if x, y, z := h.state.Camera(); x != 120 || y != 60 || z != 1.5 {
		t.Fatalf("SetCamera round-trip = (%v,%v,%v)", x, y, z)
	}
	// A relayout does not steal the camera back.
	h.data.Edges = append(h.data.Edges, Edge{Source: "test", Target: "fixed"})
	h.settle()
	if x, y, z := h.state.Camera(); x != 120 || y != 60 || z != 1.5 {
		t.Fatalf("a relayout moved the camera to (%v,%v,%v)", x, y, z)
	}
	// Fit takes it, and an out-of-range scale is clamped.
	h.state.Fit()
	if h.state.view.userMoved {
		t.Fatal("Fit left the camera user-placed")
	}
	h.state.SetCamera(0, 0, 10)
	if _, _, z := h.state.Camera(); z != 2.5 {
		t.Fatalf("SetCamera clamped zoom to %v, want 2.5", z)
	}
}

func TestHitAndHovered(t *testing.T) {
	h := newHarness(t)
	h.settle()
	r, ok := h.tt.Find("Build")
	if !ok {
		t.Fatal("Build node missing")
	}
	if id := h.state.Hit(r.X+r.W/2, r.Y+r.H/2); id != "build" {
		t.Fatalf("Hit at the Build center = %q", id)
	}
	if id := h.state.Hit(4, 690); id != "" {
		t.Fatalf("Hit on empty space = %q", id)
	}
	if id := h.state.Hovered(); id != "" {
		t.Fatalf("Hovered before any pointer move = %q", id)
	}
	h.tt.Move(r.X+r.W/2, r.Y+r.H/2)
	h.settle()
	if id := h.state.Hovered(); id != "build" {
		t.Fatalf("Hovered after moving over Build = %q", id)
	}
}

func TestDefaultModeDoesNotPanFromNodes(t *testing.T) {
	h := newHarness(t)
	v := &h.state.view
	r, ok := h.tt.Find("Build")
	if !ok {
		t.Fatal("Build node missing")
	}
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	before := v.shown
	h.tt.Press(cx, cy)
	h.tt.Move(cx+30, cy+20)
	h.tt.Release(cx+30, cy+20)
	h.settle()
	if v.shown != before {
		t.Fatalf("default mode panned from a node: %v to %v", before, v.shown)
	}
	// MyGo counts a press and release on the same element as a click no
	// matter how far the pointer travelled between them, so a drag that
	// starts and ends inside a card selects it. NavigationMap exists to
	// change exactly that: there the same gesture is a pan, never a click.
	if len(h.selected) != 1 || h.selected[0] != "build" {
		t.Fatalf("drag on a node selected %v", h.selected)
	}
}
