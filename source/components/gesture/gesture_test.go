package gesture

import (
	"math"
	"testing"
	"weak"

	"github.com/egoist/mygo/ui"
)

// Tests run off the main thread, so devices come from Precise as they do
// outside macOS.
func TestClassify(t *testing.T) {
	scroll := func(dx, dy float32, mods ui.Modifiers, precise bool) Event {
		e, ok := Classify(ui.InputEvent{Kind: ui.InputScroll, X: 10, Y: 20, DX: dx, DY: dy, Mods: mods, Precise: precise})
		if !ok {
			t.Fatal("scroll was not classified")
		}
		return e
	}
	if e := scroll(0, -120, 0, false); e.Kind != Zoom || e.Device != Wheel || e.Pinch || math.Abs(e.Scale-1.5) > 1e-9 || e.X != 10 || e.Y != 20 {
		t.Fatalf("wheel = %+v, want ×1.5 zoom", e)
	}
	if e := scroll(0, 48, ui.Ctrl, true); e.Kind != Zoom || !e.Pinch || math.Abs(e.Scale-1/1.5) > 1e-9 {
		t.Fatalf("touchpad Control-scroll = %+v, want ÷1.5 pinch zoom", e)
	}
	if e := scroll(0, -120, ui.Super, false); e.Kind != Zoom || e.Pinch {
		t.Fatalf("Command-wheel = %+v, want a wheel zoom", e)
	}
	if e := scroll(12, -30, 0, true); e.Kind != Pan || e.Device != Touchpad || e.DX != -12 || e.DY != 30 {
		t.Fatalf("touchpad scroll = %+v, want a pan", e)
	}
	if e := scroll(40, 0, ui.Shift, false); e.Kind != Pan || e.DX != -40 || e.DY != 0 {
		t.Fatalf("sideways wheel = %+v, want a pan across", e)
	}
	if _, ok := Classify(ui.InputEvent{Kind: ui.InputPointerMove}); ok {
		t.Fatal("pointer move was classified")
	}
}

func TestScrollZoomMode(t *testing.T) {
	scroll := func(dx, dy float32, mods ui.Modifiers, precise bool) Event {
		e, ok := classify(ui.InputEvent{Kind: ui.InputScroll, X: 10, Y: 20, DX: dx, DY: dy, Mods: mods, Precise: precise}, ScrollZoom)
		if !ok {
			t.Fatal("scroll was not classified")
		}
		return e
	}
	// A plain touchpad scroll zooms at the pointer, and it is not a pinch.
	if e := scroll(12, -30, 0, true); e.Kind != Zoom || e.Device != Touchpad || e.Pinch || e.X != 10 || e.Y != 20 {
		t.Fatalf("touchpad scroll under ScrollZoom = %+v, want a zoom at the pointer", e)
	}
	// Touchpad scrolls scale by the pinch distance.
	if e := scroll(0, -48, 0, true); math.Abs(e.Scale-1.5) > 1e-9 {
		t.Fatalf("touchpad scroll = %+v, want ×1.5", e)
	}
	// A wheel zooms as it does under ScrollPan, and Control still pinches.
	if e := scroll(0, -120, 0, false); e.Kind != Zoom || e.Pinch || math.Abs(e.Scale-1.5) > 1e-9 {
		t.Fatalf("wheel = %+v, want ×1.5 zoom", e)
	}
	if e := scroll(0, 48, ui.Ctrl, true); e.Kind != Zoom || !e.Pinch || math.Abs(e.Scale-1/1.5) > 1e-9 {
		t.Fatalf("Control-scroll = %+v, want ÷1.5 pinch zoom", e)
	}
	// A sideways scroll neither pans nor scales.
	if e := scroll(40, 0, 0, true); e.Kind != Zoom || math.Abs(e.Scale-1) > 1e-9 {
		t.Fatalf("sideways scroll = %+v, want a no-op zoom", e)
	}
}

func TestSetScrollModeReachesHandle(t *testing.T) {
	var nav Navigation
	nav.SetScrollMode(ScrollZoom)
	if nav.ScrollMode() != ScrollZoom {
		t.Fatalf("SetScrollMode stored %v", nav.scroll)
	}
	var got []Event
	tt := ui.NewTester(func(c *ui.Context) {
		nav.Handle(c, ui.Box(c).Size(200, 200).Label("Canvas"), func(e Event) { got = append(got, e) }, nil)
	}, 200, 200)
	tt.Frame()
	tt.Scroll(50, 60, 0, -120)
	if len(got) != 1 || got[0].Kind != Zoom {
		t.Fatalf("scroll delivered %+v, want one zoom", got)
	}
}

func TestHandleRoutesInput(t *testing.T) {
	var nav Navigation
	var got []Event
	others := 0
	tt := ui.NewTester(func(c *ui.Context) {
		nav.Handle(c, ui.Box(c).Size(200, 200).Label("Canvas"), func(e Event) { got = append(got, e) },
			func(ev ui.InputEvent) bool { others++; return ev.Kind == ui.InputPointerDown })
	}, 200, 200)
	tt.Frame()
	tt.Scroll(50, 60, 0, -120)
	if len(got) != 1 || got[0].Kind != Zoom || got[0].X != 50 || got[0].Y != 60 {
		t.Fatalf("scroll delivered %+v", got)
	}
	tt.Click("Canvas")
	if others == 0 || len(got) != 1 {
		t.Fatal("pointer input did not reach the other handler")
	}
}

func TestPinchGoesWhereThePointerIs(t *testing.T) {
	var got []Event
	n := &Navigation{fn: func(e Event) { got = append(got, e) }, x: 30, y: 40}
	n.pointer = nativePointer{window: 1, x: 100, y: 200}
	pinchTarget = weak.Make(n)
	defer func() { pinchTarget = weak.Pointer[Navigation]{} }()

	if !pinch(nativePointer{window: 1, x: 100.5, y: 200}, .25) {
		t.Fatal("pinch under the pointer was not delivered")
	}
	if len(got) != 1 || got[0].Kind != Zoom || got[0].Scale != 1.25 || got[0].X != 30 || got[0].Y != 40 || !got[0].Pinch {
		t.Fatalf("pinch delivered %+v", got)
	}
	// The pointer has since left for another element or window.
	if pinch(nativePointer{window: 1, x: 300, y: 200}, .25) || pinch(nativePointer{window: 2, x: 100, y: 200}, .25) {
		t.Fatal("pinch elsewhere reached the element")
	}
}
