package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = register(&statuses{unread: 3})

type statuses struct {
	unread int
}

func (s *statuses) View(c *ui.Context) {
	ui.Column(c).Gap(16).AlignItems(ui.Center).Children(func() {
		ui.Row(c).Gap(8).Wrap().Justify(ui.Center).Children(func() {
			components.Badge(c, "Draft")
			components.Badge(c, "New", components.BadgeOptions{Tone: components.BadgeAccent})
			components.Badge(c, "Passed", components.BadgeOptions{Tone: components.BadgeSuccess, Icon: lucide.Check()})
			components.Badge(c, "Pending", components.BadgeOptions{Tone: components.BadgeWarning, Icon: lucide.Clock()})
			components.Badge(c, "Failed", components.BadgeOptions{Tone: components.BadgeDanger, Icon: lucide.X()})
		})
		// A changed text morphs in place.
		ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
			components.Badge(c, fmt.Sprintf("%d unread", s.unread), components.BadgeOptions{Tone: components.BadgeAccent})
			if components.Button(c, "New message", components.ButtonGhost).Clicked() {
				s.unread++
			}
		})
	})
}
