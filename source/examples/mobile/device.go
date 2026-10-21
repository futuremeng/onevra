package main

import (
	"math"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// platform is the operating system whose conventions the simulated screen
// follows.
type platform int

const (
	platformIOS platform = iota
	platformAndroid
)

// orientation is how the device is held. Landscape left turns it a quarter
// counterclockwise, so its top, with the cutout, is on the left, as iOS's
// UIDeviceOrientation.landscapeLeft; landscape right turns it the other way.
type orientation int

const (
	portrait orientation = iota
	landscapeLeft
	landscapeRight
)

// cutout is the screen's camera housing.
type cutout int

const (
	cutoutNone cutout = iota
	cutoutIsland
	cutoutNotch
	cutoutHole
)

// insets are the four edges of a box, in points (dp on Android).
type insets struct{ top, right, bottom, left float32 }

// rect is a box in screen points, from the screen's top-left corner.
type rect struct{ x, y, w, h float32 }

// device is a simple model of a phone or tablet: its portrait screen, the
// body around it and the system's safe areas. Sizes are logical points, close
// to common devices rather than to any one product.
type device struct {
	// name is the message key of the model's display name.
	name     string
	platform platform
	// width and height are the portrait screen.
	width, height float32
	// radius rounds the screen's corners; the body follows it.
	radius float32
	// bezel surrounds the screen; chin adds to it above and below, where
	// home-button models keep their sensor and button.
	bezel, chin float32
	cutout      cutout
	homeButton  bool
	// portraitSafe and landscapeSafe are the system's safe-area insets,
	// landscapeSafe with the device turned landscape left.
	portraitSafe, landscapeSafe insets
	// keyboard is the software keyboard's height in portrait and landscape,
	// including the edge below it.
	keyboard [2]float32
	// handle draws a home indicator or gesture handle along the bottom.
	handle bool
}

// devices are the models the preview offers, phones before tablets.
var devices = []device{
	{name: "deviceIsland", platform: platformIOS, width: 393, height: 852, radius: 55, bezel: 14,
		cutout: cutoutIsland, handle: true,
		portraitSafe: insets{59, 0, 34, 0}, landscapeSafe: insets{0, 59, 21, 59}, keyboard: [2]float32{336, 209}},
	{name: "deviceNotch", platform: platformIOS, width: 390, height: 844, radius: 47, bezel: 14,
		cutout: cutoutNotch, handle: true,
		portraitSafe: insets{47, 0, 34, 0}, landscapeSafe: insets{0, 47, 21, 47}, keyboard: [2]float32{336, 209}},
	{name: "deviceHomeButton", platform: platformIOS, width: 375, height: 667, bezel: 18, chin: 78,
		homeButton:   true,
		portraitSafe: insets{20, 0, 0, 0}, keyboard: [2]float32{260, 200}},
	{name: "deviceAndroidPhone", platform: platformAndroid, width: 412, height: 915, radius: 34, bezel: 12,
		cutout: cutoutHole, handle: true,
		portraitSafe: insets{52, 0, 24, 0}, landscapeSafe: insets{28, 0, 24, 52}, keyboard: [2]float32{300, 220}},
	{name: "deviceIPad", platform: platformIOS, width: 820, height: 1180, radius: 18, bezel: 24, handle: true,
		portraitSafe: insets{24, 0, 20, 0}, landscapeSafe: insets{24, 0, 20, 0}, keyboard: [2]float32{313, 398}},
	{name: "deviceAndroidTablet", platform: platformAndroid, width: 800, height: 1280, radius: 20, bezel: 28,
		portraitSafe: insets{28, 0, 48, 0}, landscapeSafe: insets{28, 0, 48, 0}, keyboard: [2]float32{360, 330}},
}

func (o orientation) landscape() bool { return o != portrait }

// angle is how far the device is turned, in degrees clockwise.
func (o orientation) angle() float32 {
	switch o {
	case landscapeLeft:
		return -90
	case landscapeRight:
		return 90
	}
	return 0
}

// edges turns portrait edges with the device: turned left, the top edge
// becomes the left one.
func (o orientation) edges(e insets) insets {
	switch o {
	case landscapeLeft:
		return insets{top: e.right, right: e.bottom, bottom: e.left, left: e.top}
	case landscapeRight:
		return insets{top: e.left, right: e.top, bottom: e.right, left: e.bottom}
	}
	return e
}

// place turns a box given in portrait screen points, which may lie outside
// the w × h screen, with the device. Corner radii follow as edges do, in
// top-left, top-right, bottom-right, bottom-left order.
func (o orientation) place(r rect, w, h float32) rect {
	switch o {
	case landscapeLeft:
		return rect{r.y, w - r.x - r.w, r.h, r.w}
	case landscapeRight:
		return rect{h - r.y - r.h, r.x, r.h, r.w}
	}
	return r
}

func (o orientation) corners(c [4]float32) [4]float32 {
	switch o {
	case landscapeLeft:
		return [4]float32{c[1], c[2], c[3], c[0]}
	case landscapeRight:
		return [4]float32{c[3], c[0], c[1], c[2]}
	}
	return c
}

// screen returns the screen's size as held.
func (d device) screen(o orientation) (w, h float32) {
	if o.landscape() {
		return d.height, d.width
	}
	return d.width, d.height
}

// body returns the bezels around the screen as held.
func (d device) body(o orientation) insets {
	return o.edges(insets{d.bezel + d.chin, d.bezel, d.bezel + d.chin, d.bezel})
}

// safe returns the safe-area insets as held. A shown keyboard raises the
// bottom edge to its top, as both systems report the keyboard's inset.
func (d device) safe(o orientation, keyboard bool) insets {
	s := d.portraitSafe
	if o.landscape() {
		s = d.landscapeSafe
		if o == landscapeRight {
			s.left, s.right = s.right, s.left
		}
	}
	if keyboard {
		s.bottom = max(s.bottom, d.keyboardHeight(o))
	}
	return s
}

func (d device) keyboardHeight(o orientation) float32 {
	if o.landscape() {
		return d.keyboard[1]
	}
	return d.keyboard[0]
}

// compact reports a screen narrow enough for one pane. iOS keeps phones in
// the compact width class when turned; Android's window size classes follow
// the current width.
func (d device) compact(o orientation) bool {
	if d.platform == platformIOS {
		return min(d.width, d.height) < 600
	}
	w, _ := d.screen(o)
	return w < 600
}

// cutoutShape returns the camera housing in portrait screen points with its
// corner radii, or false without one.
func (d device) cutoutShape() (rect, [4]float32, bool) {
	switch d.cutout {
	case cutoutIsland:
		return rect{(d.width - 126) / 2, 11, 126, 37}, [4]float32{18.5, 18.5, 18.5, 18.5}, true
	case cutoutNotch:
		return rect{(d.width - 162) / 2, 0, 162, 32}, [4]float32{0, 0, 20, 20}, true
	case cutoutHole:
		return rect{(d.width - 24) / 2, 14, 24, 24}, [4]float32{12, 12, 12, 12}, true
	}
	return rect{}, [4]float32{}, false
}

// fit returns the scale that shows a w × h box turned by degrees within
// availW × availH, never above actual size.
func fit(w, h, degrees, availW, availH float32) float32 {
	if w <= 0 || h <= 0 || availW <= 0 || availH <= 0 {
		return 1
	}
	rad := float64(degrees) * math.Pi / 180
	cos, sin := float32(math.Abs(math.Cos(rad))), float32(math.Abs(math.Sin(rad)))
	bw, bh := w*cos+h*sin, w*sin+h*cos
	return min(1, availW/bw, availH/bh)
}

// frameView builds the device turned as held and scaled into availW ×
// availH; screen builds what the screen shows within w × h points.
func (a *app) frameView(c *ui.Context, availW, availH float32, screen func(c *ui.Context, w, h float32, safe insets)) {
	d, o := devices[a.device], a.orientation
	sw, sh := d.screen(o)
	b := d.body(o)
	fw, fh := sw+b.left+b.right, sh+b.top+b.bottom
	safe := d.safe(o, a.keyboardShown())

	frame := ui.Box(c).Key("device").Absolute().Attach(ui.AnchorCenter, ui.AnchorCenter).Size(fw, fh)
	// The frame is laid out as held and turns from the previous orientation,
	// as a simulator's window does.
	turn := frame.AnimateTo("turn", o.angle(), ui.Smooth) - o.angle()
	k := fit(fw, fh, turn, availW, availH)
	a.scale = k
	frame.Scale(k, k).Rotate(turn).
		Radius(d.radius+min(b.left, b.top)).
		Background(ui.Hex("#1d1d1f")).InsetShadow(0, 0, 0, 2, ui.Hex("#3a3a3c")).
		Shadow(0, 24, 48, 0, ui.RGBA(0, 0, 0, 0.28)).
		Children(func() {
			if d.homeButton {
				a.homeButtonView(c, d, o, b)
			}
			ui.Box(c).Absolute().Left(b.left).Top(b.top).Size(sw, sh).
				Radius(d.radius).Clip().Background(c.Theme().Background).
				Children(func() {
					screen(c, sw, sh, safe)
					a.hardwareView(c, d, o, sw, sh, safe)
				})
		})
}

// homeButtonView draws the round button and the earpiece of home-button
// models in the body, beside the screen.
func (a *app) homeButtonView(c *ui.Context, d device, o orientation, b insets) {
	at := func(r rect) ui.Element {
		p := o.place(r, d.width, d.height)
		return ui.Box(c).Absolute().Left(b.left+p.x).Top(b.top+p.y).Size(p.w, p.h).PassThrough()
	}
	at(rect{(d.width - 56) / 2, d.height + (d.chin+d.bezel-56)/2, 56, 56}).
		Radius(28).Border(2, ui.Hex("#48484a"))
	at(rect{(d.width - 60) / 2, -(d.chin + d.bezel + 6) / 2, 60, 6}).
		Radius(3).Background(ui.Hex("#2c2c2e"))
}

// hardwareView draws over the screen what the system and the hardware put
// there: the status bar, the cutout, the keyboard, the home indicator and the
// optional safe-area overlay.
func (a *app) hardwareView(c *ui.Context, d device, o orientation, w, h float32, safe insets) {
	t := c.Theme()
	if a.keyboardShown() {
		keyboardView(c, w, d.keyboardHeight(o))
	}
	statusBarView(c, d, o, w, safe)
	if r, corners, ok := d.cutoutShape(); ok {
		p, rc := o.place(r, d.width, d.height), o.corners(corners)
		ui.Box(c).Absolute().Left(p.x).Top(p.y).Size(p.w, p.h).PassThrough().
			Radius(rc[0], rc[1], rc[2], rc[3]).Background(ui.Hex("#000000"))
	}
	if d.handle && !a.keyboardShown() {
		handle := t.Text
		width := min(134, w*0.36)
		if d.platform == platformAndroid {
			width = 108
		}
		ui.Box(c).Absolute().Bottom(8).Left((w-width)/2).Size(width, 5).Radius(2.5).
			Background(handle.Alpha(0.85)).PassThrough()
	}
	if a.overlay {
		safeOverlay(c, w, h, safe)
	}
}

// statusBarView draws the system's clock and indicators in the top safe area,
// beside the cutout. iPhones hide the status bar in landscape.
func statusBarView(c *ui.Context, d device, o orientation, w float32, safe insets) {
	if safe.top == 0 {
		return
	}
	bar := min(safe.top, 44)
	if d.platform == platformAndroid {
		bar = min(safe.top, 32)
	}
	top := (safe.top - bar) / 2
	if d.cutout == cutoutIsland && !o.landscape() {
		top = 11
		bar = 37
	}
	side := max(safe.left, safe.right, d.radius*0.45, 8)
	if d.platform == platformAndroid {
		side = max(safe.left, safe.right, 20)
	}
	ui.Row(c).Absolute().Top(top).Left(side).Right(side).Height(bar).PassThrough().
		FontSize(min(15, bar*0.6)).Bold().Children(func() {
		clock := ui.Text(c, "9:41")
		if d.cutout == cutoutIsland || d.cutout == cutoutNotch {
			// Face ID iPhones center the clock in the space left of the cutout.
			clock.Width(w/2 - side - 81).TextAlign(ui.Center)
		}
		ui.Spacer(c)
		ui.Row(c).Gap(6).Children(func() {
			for _, glyph := range []func() *ui.SVG{lucide.SignalHigh, lucide.Wifi, lucide.BatteryFull} {
				components.Icon(c, glyph).Size(17, 17)
			}
		})
	})
}

// keyboardView draws a placeholder software keyboard along the bottom: three
// rows of letter keys and a space bar.
func keyboardView(c *ui.Context, w, h float32) {
	t := c.Theme()
	keyHeight := min(42, (h-56)/4)
	ui.Column(c).Absolute().Left(0).Right(0).Bottom(0).Height(h).
		Background(t.Surface.Mix(t.Text, 0.08)).Padding(8, 4).Gap(10).PassThrough().
		Children(func() {
			for _, keys := range []int{10, 9, 7, 0} {
				ui.Row(c).Gap(6).Justify(ui.Center).Children(func() {
					if keys == 0 {
						ui.Box(c).Width(w * 0.5).Height(keyHeight).Radius(6).Background(t.Background)
						return
					}
					for range keys {
						ui.Box(c).Grow(1).MaxWidth(w / 11).Height(keyHeight).Radius(6).Background(t.Background)
					}
				})
			}
		})
}

// safeOverlay tints the insets outside the safe area and outlines the area
// itself.
func safeOverlay(c *ui.Context, w, h float32, s insets) {
	tint := ui.RGBA(255, 59, 48, 0.22)
	band := func() ui.Element { return ui.Box(c).Absolute().Background(tint).PassThrough() }
	if s.top > 0 {
		band().Top(0).Left(0).Size(w, s.top)
	}
	if s.bottom > 0 {
		band().Bottom(0).Left(0).Size(w, s.bottom)
	}
	if s.left > 0 {
		band().Top(s.top).Left(0).Size(s.left, h-s.top-s.bottom)
	}
	if s.right > 0 {
		band().Top(s.top).Right(0).Size(s.right, h-s.top-s.bottom)
	}
	ui.Box(c).Absolute().Top(s.top).Left(s.left).Size(w-s.left-s.right, h-s.top-s.bottom).
		Border(1, ui.RGBA(255, 59, 48, 0.9)).BorderStyle(ui.BorderDashed).PassThrough()
}
