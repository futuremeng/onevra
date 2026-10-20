package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = register(&buttonIcons{})

type buttonIcons struct{}

func (buttonIcons) View(c *ui.Context) {
	t := c.Theme()
	// A button with an empty label takes any content as its children; Label
	// keeps its accessible name.
	iconButton := func(label string, icon *ui.SVG, variant components.ButtonVariant) {
		ui.Box(c).Key(label).Children(func() {
			components.Button(c, "", variant).Label(label).Gap(t.Space(1.5)).Children(func() {
				components.Icon(c, icon).Size(t.Rem(1), t.Rem(1))
				components.Text(c, label).SingleLine()
			})
		})
	}
	ui.Row(c).Wrap().Gap(8).Justify(ui.Center).Children(func() {
		iconButton("New", lucide.Plus(), components.ButtonPrimary)
		iconButton("Download", lucide.Download(), components.ButtonSecondary)
		iconButton("Share", lucide.Share(), components.ButtonOutline)
		iconButton("Edit", lucide.Pencil(), components.ButtonGhost)
	})
}
