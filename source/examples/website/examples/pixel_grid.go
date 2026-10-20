package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = registerWindow(&canvas{})

type canvas struct{}

func (canvas) View(c *ui.Context) {
	t := c.Theme()
	components.PixelGrid(c, components.PixelGridOptions{Spacing: 20, MajorEvery: 5}).Fill().Children(func() {
		ui.Column(c).Fill().Center().Gap(t.Space(2)).Children(func() {
			ui.Column(c).Padding(t.Space(4)).Gap(t.Space(1)).Radius(t.Radius*1.5).
				Background(t.Surface).Border(1, t.Border).Children(func() {
				components.Text(c, "Artboard").Bold()
				components.Text(c, "Children draw above the grid.").TextColor(t.TextMuted)
			})
		})
	})
}
