package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/gesture"
)

var _ = registerWindow(&zoomCanvas{x: 40, y: 40, zoom: 1})

type zoomCanvas struct {
	nav        gesture.Navigation
	x, y, zoom float32 // the artboard's position and scale
}

func (z *zoomCanvas) View(c *ui.Context) {
	t := c.Theme()
	canvas := components.PixelGrid(c, components.PixelGridOptions{Spacing: 20}).Fill().Clip()
	// Scroll or drag two fingers to pan; use the wheel or pinch to zoom.
	z.nav.Handle(c, canvas, func(e gesture.Event) {
		switch e.Kind {
		case gesture.Pan:
			z.x += e.DX
			z.y += e.DY
		case gesture.Zoom:
			next := min(max(z.zoom*float32(e.Scale), 0.25), 4)
			// Zoom around the pointer: keep the point under it in place.
			f := next / z.zoom
			z.x = e.X - (e.X-z.x)*f
			z.y = e.Y - (e.Y-z.y)*f
			z.zoom = next
		}
	}, nil)
	canvas.Children(func() {
		ui.Column(c).Absolute().Left(z.x).Top(z.y).Width(240*z.zoom).Height(160*z.zoom).
			Padding(t.Space(4)*z.zoom).Radius(t.Radius*1.5*z.zoom).
			Background(t.Surface).Border(1, t.Border).Children(func() {
			components.Text(c, "Artboard").Bold().FontSize(t.FontSize * z.zoom)
		})
	})
}
