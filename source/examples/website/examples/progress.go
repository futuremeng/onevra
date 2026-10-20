package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&upload{done: 0.4})

type upload struct {
	done float64
}

func (u *upload) View(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Width(320).MaxWidthPercent(100).Gap(16).Children(func() {
		ui.Row(c).FillWidth().Children(func() {
			components.Text(c, "Uploading photos").Grow(1)
			components.Text(c, fmt.Sprintf("%.0f%%", u.done*100)).TextColor(t.TextMuted)
		})
		components.Progress(c, u.done)
		ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
			components.ProgressRing(c, u.done)
			if components.Button(c, "Advance", components.ButtonSecondary).Clicked() {
				u.done = min(1, u.done+0.2)
			}
			if components.Button(c, "Reset", components.ButtonGhost).Clicked() {
				u.done = 0
			}
		})
	})
}
