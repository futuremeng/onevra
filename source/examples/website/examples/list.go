package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&inbox{})

type inbox struct {
	list     components.ListState
	selected int
}

func (b *inbox) View(c *ui.Context) {
	t := c.Theme()
	b.list.Selected = &b.selected
	ui.Column(c).Width(360).MaxWidthPercent(100).Gap(8).Children(func() {
		// Only the visible rows of the 10,000 are built.
		components.List(c, &b.list, 10000, func(i int) {
			ui.Column(c).FillWidth().Padding(t.Space(2), t.Space(3)).Children(func() {
				components.Text(c, fmt.Sprintf("Message %d", i+1)).Bold()
				components.Text(c, "A preview of the message's first line").TextColor(t.TextMuted).SingleLine()
			})
		}).Height(260).FillWidth().Label("Messages").Radius(t.Radius).Border(1, t.Border)
		components.Text(c, fmt.Sprintf("Selected message %d", b.selected+1)).TextColor(t.TextMuted)
	})
}
