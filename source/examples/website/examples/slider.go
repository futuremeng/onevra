package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&mixer{volume: 60, quality: 40, low: 20, high: 80})

type mixer struct {
	volume, quality, low, high float64
}

func (m *mixer) View(c *ui.Context) {
	t := c.Theme()
	row := func(label, value string, slider func()) {
		ui.Column(c).Key(label).FillWidth().Gap(t.Space(1)).Children(func() {
			ui.Row(c).FillWidth().Children(func() {
				components.Text(c, label).Grow(1)
				components.Text(c, value).TextColor(t.TextMuted)
			})
			slider()
		})
	}
	ui.Column(c).Width(320).MaxWidthPercent(100).Gap(16).Children(func() {
		row("Volume", fmt.Sprintf("%.0f", m.volume), func() {
			components.Slider(c, &m.volume, 0, 100).Label("Volume").FillWidth()
		})
		row("Quality", fmt.Sprintf("%.0f", m.quality), func() {
			components.Slider(c, &m.quality, 0, 100, components.SliderOptions{Step: 20}).
				Label("Quality").FillWidth()
		})
		row("Price", fmt.Sprintf("$%.0f – $%.0f", m.low, m.high), func() {
			components.Slider(c, &m.low, 0, 100, components.SliderOptions{Range: &m.high}).
				Label("Price").FillWidth()
		})
	})
}
