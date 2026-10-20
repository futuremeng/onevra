package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&deleteProject{})

type deleteProject struct {
	deleted bool
}

func (d *deleteProject) View(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).Gap(12).AlignItems(ui.Center).Children(func() {
		hold := components.HoldToConfirm(c.Key("delete"), &d.deleted, "Hold to delete", components.HoldToConfirmOptions{
			ConfirmedLabel: "Deleted",
			Tone:           components.HoldDanger,
		})
		switch {
		case hold.Holding():
			components.Text(c, "Keep holding…").TextColor(t.TextMuted)
		case d.deleted:
			if components.Button(c, "Undo", components.ButtonGhost).Clicked() {
				d.deleted = false
			}
		}
	})
}
