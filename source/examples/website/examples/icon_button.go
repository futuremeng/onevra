package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = register(&iconButtons{})

type iconButtons struct {
	starred bool
}

func (b *iconButtons) View(c *ui.Context) {
	ui.Column(c).Gap(16).AlignItems(ui.Center).Children(func() {
		// Plain buttons tint a transparent face; hover one for its tooltip.
		ui.Row(c).Gap(4).Children(func() {
			components.IconButton(c, lucide.Copy(), "Copy", components.IconButtonPlain)
			star := lucide.Star()
			if b.starred {
				star = lucide.StarOff()
			}
			if components.IconButton(c, star, "Favorite", components.IconButtonPlain).Clicked() {
				b.starred = !b.starred
			}
			components.IconButton(c, lucide.RefreshCw(), "Refresh", components.IconButtonPlain)
			components.IconButton(c, lucide.Trash(), "Delete", components.IconButtonPlain).Disabled(true)
		})
		// Toolbar buttons raise their face, as ButtonToolbar does.
		ui.Row(c).Gap(8).Children(func() {
			components.Button(c, "Filter", components.ButtonToolbar)
			components.IconButton(c, lucide.Undo(), "Undo", components.IconButtonToolbar)
			components.IconButton(c, lucide.Redo(), "Redo", components.IconButtonToolbar)
		})
	})
}
