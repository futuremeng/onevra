package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = registerWindow(&shareMenu{})

type shareMenu struct {
	open, comments, download bool
}

func (s *shareMenu) View(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Fill().Padding(t.Space(6)).AlignItems(ui.Center).Children(func() {
		anchor := components.Button(c, "Share", components.ButtonOutline)
		if anchor.Clicked() {
			s.open = !s.open
		}
		components.Popover(c, anchor, &s.open, func(panel ui.Element) {
			panel.Label("Sharing").Width(260)
			components.Text(c, "Anyone with the link").Bold()
			components.Text(c, "People you share with can view the board.").TextColor(t.TextMuted)
			components.Checkbox(c, &s.comments, "Allow comments")
			components.Checkbox(c, &s.download, "Allow downloads")
			if components.Button(c, "Done", components.ButtonPrimary).Clicked() {
				s.open = false
			}
		})
	})
}
