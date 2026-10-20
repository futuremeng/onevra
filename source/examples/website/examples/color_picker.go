package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&brandColors{
	primary: ui.RGB(189, 72, 27),
	accent:  ui.RGB(72, 150, 214),
})

type brandColors struct {
	primary, accent ui.Color
}

func (b *brandColors) View(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Gap(16).AlignItems(ui.Center).Children(func() {
		ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
			components.ColorPicker(c, &b.primary).Label("Primary color")
			components.ColorPicker(c, &b.accent).Label("Accent color").
				Size(40, 40).Radius(20).Padding(0)
		})
		ui.Box(c).Width(220).Height(40).Radius(t.Radius).
			LinearGradient(ui.LinearGradient{From: b.primary, To: b.accent, Angle: 90})
	})
}
