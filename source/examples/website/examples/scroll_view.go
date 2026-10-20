package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&scrollers{})

type scrollers struct{}

func (scrollers) View(c *ui.Context) {
	t := c.Theme()
	frame := func(e ui.Element) ui.Element {
		return e.Radius(t.Radius).Background(t.Surface).Border(1, t.Border)
	}
	ui.Row(c).Gap(16).Wrap().Justify(ui.Center).Children(func() {
		frame(components.ScrollView(c, components.ScrollViewOptions{}).Width(220).Height(200)).
			Padding(t.Space(3)).Gap(t.Space(2)).Children(func() {
			for i := 1; i <= 20; i++ {
				components.Text(c, fmt.Sprintf("Row %d", i)).Shrink(0)
			}
		})
		frame(components.ScrollView(c, components.ScrollViewOptions{Axes: components.ScrollHorizontal}).
			Width(220).Height(200)).Padding(t.Space(3)).Gap(t.Space(2)).Children(func() {
			for i := 1; i <= 8; i++ {
				ui.Column(c).Width(90).FillHeight().Shrink(0).Center().Radius(t.Radius).
					Background(t.SurfaceHover).Children(func() {
					components.Text(c, fmt.Sprintf("Card %d", i))
				})
			}
		})
	})
}
