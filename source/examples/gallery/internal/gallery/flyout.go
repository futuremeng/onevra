package gallery

import (
	"runtime"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/theme"
)

// flyoutDemo retains application state separately from the transient window.
type flyoutDemo struct {
	gallery *Gallery
	window  *mygo.Window
	note    string
}

func (d *flyoutDemo) close() {
	if d.window != nil {
		d.window.Destroy()
	}
	d.window = nil
}

// flyoutStyle is how the page's flyout looks: the editor draws its own
// surface, the system popover is the platform's, and the material flyout
// shows a blurred material cut to rounded corners where the platform has one.
type flyoutStyle int

const (
	flyoutEditor flyoutStyle = iota
	flyoutPopover
	flyoutMaterial
)

func (d *flyoutDemo) show(bounds ui.Rect, style flyoutStyle) {
	if d.gallery == nil || d.gallery.Shell == nil {
		return
	}
	d.close()
	parent := d.gallery.Shell.Window()
	o := mygo.FlyoutOptions{
		Parent:    parent,
		Anchor:    mygo.Rectangle{X: int(bounds.X), Y: int(bounds.Y), Width: int(bounds.W), Height: int(bounds.H)},
		Placement: mygo.PlacementBottomStart, Gap: 6,
		Width: 300, Height: 180, Focusable: true, Shadow: true, Popover: style == flyoutPopover,
		Content: ui.View(d.view),
	}
	if style == flyoutMaterial {
		// Windows 11 rounds flyouts itself; CornerRadius rounds them on macOS.
		o.Vibrancy, o.CornerRadius = mygo.VibrancyPopover, 10
	}
	w := mygo.NewFlyout(o)
	d.window = w
	w.OnClosed(func() {
		if d.window == w {
			d.window = nil
		}
		parent.Invalidate()
	})
}

func (d *flyoutDemo) view(c *ui.Context) {
	theme.Apply(c)
	c.Root().Background(ui.Transparent)
	d.gallery.withLanguage(c, func() {
		t := c.Theme()
		panel := ui.Column(c).Fill().Padding(16).Gap(12)
		if !c.Vibrancy() {
			panel.Background(t.Surface)
			switch runtime.GOOS {
			case "darwin":
				panel.Radius(10)
			case "windows":
			default:
				panel.Radius(10).Border(1, t.Border)
			}
		}
		panel.Children(func() {
			components.TextBox(c, &d.note, components.TextBoxOptions{}).Label(tr(c, "notes")).FillWidth().AutoFocus()
			components.Button(c, tr(c, "close"), components.ButtonOutline).OnClick(func() {
				if d.window != nil {
					d.window.Close()
				}
			})
		})
	})
}

func exampleFlyout(c *ui.Context, s *demoState) {
	d := &s.flyout
	components.Text(c, tr(c, "flyoutDescription")).Wrap()
	for i, label := range []string{"flyoutEditor", "flyoutPopover", "flyoutMaterial"} {
		b := components.Button(c, tr(c, label), components.ButtonOutline).Disabled(d.gallery == nil || d.gallery.Shell == nil)
		if b.Clicked() {
			d.show(b.Bounds(), flyoutStyle(i))
		}
	}
	components.Text(c, d.note).Wrap()
}
