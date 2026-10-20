package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = register(&iconSizes{})

type iconSizes struct{}

func (iconSizes) View(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Gap(16).AlignItems(ui.Center).Children(func() {
		ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
			for _, size := range []float32{16, 20, 24, 32, 48} {
				components.Icon(c, lucide.Rocket).Size(size, size)
			}
		})
		ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
			components.Icon(c, lucide.Heart).Size(24, 24).TextColor(t.Danger)
			components.Icon(c, lucide.CircleCheck).Size(24, 24).TextColor(t.Success)
			components.Icon(c, lucide.TriangleAlert).Size(24, 24).TextColor(t.Warning)
			components.Icon(c, lucide.Sparkles).Size(24, 24).TextColor(t.Accent)
		})
	})
}
