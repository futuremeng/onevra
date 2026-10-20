package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = register(&formatting{bold: true})

type formatting struct {
	bold, italic, underline bool
}

func (f *formatting) View(c *ui.Context) {
	size := c.Theme().Space(4)
	ui.Row(c).Gap(8).Children(func() {
		components.ToggleButton(c, &f.bold, "Bold")
		// An empty label takes custom content; Label names it.
		components.ToggleButton(c, &f.italic, "").Label("Italic").Children(func() {
			components.Icon(c, lucide.Italic).Size(size, size)
		})
		components.ToggleButton(c, &f.underline, "").Label("Underline").Children(func() {
			components.Icon(c, lucide.Underline).Size(size, size)
		})
	})
}
