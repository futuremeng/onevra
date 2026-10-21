package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

func TestDeviceGeometry(t *testing.T) {
	for _, d := range devices {
		w, h := d.screen(portrait)
		lw, lh := d.screen(landscapeLeft)
		if lw != h || lh != w {
			t.Errorf("%s: landscape screen %v × %v, want %v × %v", d.name, lw, lh, h, w)
		}
		pb, lb := d.body(portrait), d.body(landscapeLeft)
		if pb.top != lb.left || pb.left != lb.bottom || pb.top+pb.bottom != lb.left+lb.right {
			t.Errorf("%s: body %+v does not turn into %+v", d.name, pb, lb)
		}
		left, right := d.safe(landscapeLeft, false), d.safe(landscapeRight, false)
		if left.left != right.right || left.right != right.left || left.top != right.top || left.bottom != right.bottom {
			t.Errorf("%s: landscape right %+v does not mirror left %+v", d.name, right, left)
		}
		for _, o := range []orientation{portrait, landscapeLeft, landscapeRight} {
			if s := d.safe(o, true); s.bottom != max(d.safe(o, false).bottom, d.keyboardHeight(o)) {
				t.Errorf("%s %v: keyboard leaves bottom inset %v", d.name, o, s.bottom)
			}
			r, _, ok := d.cutoutShape()
			if !ok {
				continue
			}
			sw, sh := d.screen(o)
			p := o.place(r, d.width, d.height)
			if p.x < 0 || p.y < 0 || p.x+p.w > sw || p.y+p.h > sh {
				t.Errorf("%s %v: cutout %+v outside %v × %v", d.name, o, p, sw, sh)
			}
			// The cutout sits within the safe-area inset of its edge.
			s := d.safe(o, false)
			switch o {
			case portrait:
				if p.y+p.h > s.top {
					t.Errorf("%s: cutout below top inset %v", d.name, s.top)
				}
			case landscapeLeft:
				if p.x+p.w > s.left {
					t.Errorf("%s: cutout beyond left inset %v", d.name, s.left)
				}
			case landscapeRight:
				if p.x < sw-s.right {
					t.Errorf("%s: cutout beyond right inset %v", d.name, s.right)
				}
			}
		}
	}
}

func TestDeviceCompact(t *testing.T) {
	phone, tablet, android := devices[0], devices[4], devices[3]
	if !phone.compact(landscapeLeft) || tablet.compact(portrait) {
		t.Error("iOS width classes")
	}
	if !android.compact(portrait) || android.compact(landscapeLeft) {
		t.Error("Android width classes follow the current width")
	}
}

func TestFit(t *testing.T) {
	if k := fit(100, 200, 0, 400, 100); k != 0.5 {
		t.Errorf("upright fit %v", k)
	}
	if k := fit(100, 200, 90, 100, 400); k != 0.5 {
		t.Errorf("turned fit %v", k)
	}
	if k := fit(100, 200, 0, 1000, 1000); k != 1 {
		t.Errorf("fit %v enlarges", k)
	}
}

func TestRotate(t *testing.T) {
	a := &app{}
	a.rotate(-1)
	if a.orientation != landscapeLeft {
		t.Fatalf("counterclockwise turn gives %v", a.orientation)
	}
	a.rotate(-1)
	if a.orientation != landscapeLeft {
		t.Fatalf("turning past landscape gives %v", a.orientation)
	}
	a.rotate(1)
	a.rotate(1)
	if a.orientation != landscapeRight {
		t.Fatalf("clockwise turns give %v", a.orientation)
	}
}

// tester shows the preview with the device and orientation chosen. Reduced
// motion settles rotations and sheets at once.
func tester(t *testing.T, a *app) *ui.Tester {
	t.Helper()
	tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); a.view(c) }, 1120, 830)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Frame()
	tt.Advance(100 * time.Millisecond)
	return tt
}

func TestNavigation(t *testing.T) {
	for _, device := range []int{0, 3} {
		a := &app{device: device}
		tt := tester(t, a)
		if err := tt.Click("Switch"); err != nil {
			t.Fatal(err)
		}
		tt.Advance(100 * time.Millisecond)
		if a.selected != "Switch" || !tt.HasText("Prediction") || !tt.HasText("Wi-Fi") {
			t.Fatalf("device %d: detail of %q not shown", device, a.selected)
		}
		if err := tt.Click("Back"); err != nil {
			t.Fatal(err)
		}
		tt.Advance(100 * time.Millisecond)
		if a.selected != "" {
			t.Fatalf("device %d: back left %q open", device, a.selected)
		}
	}
	// Wide tablets keep the list beside the detail.
	a := &app{device: 4, orientation: landscapeLeft}
	tt := tester(t, a)
	if !tt.HasText("Choose a component") {
		t.Fatal("tablet has no empty detail")
	}
	if err := tt.Click("Slider"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(100 * time.Millisecond)
	if !tt.HasText("Prediction") || !tt.HasText("Button") {
		t.Fatal("tablet does not show list and detail together")
	}
}

func TestTabsAndSheet(t *testing.T) {
	a := &app{}
	tt := tester(t, a)
	if err := tt.Click("Device"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(100 * time.Millisecond)
	if a.tab != tabDevice || !tt.HasText("393 × 852 pt") || !tt.HasText("59") {
		t.Fatalf("device tab %d does not report the screen", a.tab)
	}
	if err := tt.Click("Open sheet"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(100 * time.Millisecond)
	if !a.sheet {
		t.Fatal("sheet did not open")
	}
	if err := tt.Click("Done"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(100 * time.Millisecond)
	if a.sheet {
		t.Fatal("sheet did not close")
	}
}

func TestIOSSearchTab(t *testing.T) {
	a := &app{}
	tt := tester(t, a)
	if err := tt.Click("Search"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(100 * time.Millisecond)
	if a.tab != tabSearch || !tt.HasText("Search components") {
		t.Fatal("the tab bar did not turn into the search field")
	}
	a.query = "slid"
	tt.Frame()
	if !tt.HasText("Slider") || tt.HasText("Switch") {
		t.Fatal("search does not filter the catalog")
	}
	// The button beside the field returns to the last tab.
	if err := tt.Click("Components"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(100 * time.Millisecond)
	if a.tab != tabComponents {
		t.Fatalf("returned to tab %d", a.tab)
	}
}

func TestIOSTabBarMinimizes(t *testing.T) {
	a := &app{}
	tt := tester(t, a)
	for range 4 {
		tt.Wheel(400, 400, 0, 120)
		tt.Advance(50 * time.Millisecond)
	}
	if !a.minimized || tt.HasText("Overview") {
		t.Fatal("scrolling down does not minimize the tab bar")
	}
	tt.Wheel(400, 400, 0, -120)
	tt.Advance(50 * time.Millisecond)
	if a.minimized {
		t.Fatal("scrolling up does not restore the tab bar")
	}
}

func TestIOSLandscapeSidebar(t *testing.T) {
	a := &app{orientation: landscapeLeft}
	tt := tester(t, a)
	if !tt.HasText("Onevra Mobile") {
		t.Fatal("landscape iPhone shows no sidebar")
	}
	if err := tt.Click("Overview"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(100 * time.Millisecond)
	if a.tab != tabOverview {
		t.Fatalf("sidebar chose tab %d", a.tab)
	}
}

func TestKeyboardFollowsTextField(t *testing.T) {
	a := &app{selected: "TextBox"}
	tt := tester(t, a)
	if a.keyboardShown() {
		t.Fatal("keyboard shown before focus")
	}
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(100 * time.Millisecond)
	if !a.keyboardShown() {
		t.Fatal("focusing the field does not show the keyboard")
	}
}

// TestCatalogMessages checks that every message key the preview uses exists
// in each catalog.
func TestCatalogMessages(t *testing.T) {
	keys := append([]string{}, groups...)
	keys = append(keys, orientationNames...)
	for _, d := range devices {
		keys = append(keys, d.name)
	}
	for _, s := range supports {
		keys = append(keys, s.name(), s.summary())
	}
	for _, tab := range tabs {
		keys = append(keys, tab.label)
	}
	seen := map[string]bool{}
	for _, e := range catalog {
		if seen[e.id] {
			t.Errorf("duplicate entry %s", e.id)
		}
		seen[e.id] = true
		keys = append(keys, e.group, e.note)
	}
	files, _ := filepath.Glob("locales/*.json")
	if len(files) < 2 {
		t.Fatalf("catalogs %v", files)
	}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		messages := map[string]any{}
		if err := json.Unmarshal(data, &messages); err != nil {
			t.Fatal(err)
		}
		for _, key := range keys {
			if _, ok := messages[key]; !ok {
				t.Errorf("%s lacks %q", file, key)
			}
		}
	}
}
