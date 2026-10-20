package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&busy{})

type busy struct{}

func (busy) View(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).Gap(24).AlignItems(ui.Center).Children(func() {
		components.Spinner(c, components.SpinnerRing).Label("Loading")
		components.Spinner(c, components.SpinnerSpokes).Label("Loading")
		components.Spinner(c, components.SpinnerSpokes).Label("Syncing").TextColor(t.Accent)
	})
}
