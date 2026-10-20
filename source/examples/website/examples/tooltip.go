package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&publish{})

type publish struct{}

func (publish) View(c *ui.Context) {
	ui.Row(c).Gap(8).Children(func() {
		save := components.Button(c, "Save draft", components.ButtonSecondary)
		components.Tooltip(c, save, "Saves without publishing")
		send := components.Button(c, "Publish", components.ButtonPrimary)
		components.Tooltip(c, send, "Makes the post visible to everyone")
	})
}
