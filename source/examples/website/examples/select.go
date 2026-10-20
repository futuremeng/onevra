package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&sizeChoice{size: "Medium"})

type sizeChoice struct {
	size string
}

func (s *sizeChoice) View(c *ui.Context) {
	sel := components.Select(c, &s.size)
	sel.Trigger.Label("Size").Children(func() {
		components.Text(c, s.size).SingleLine()
	})
	sel.Popup(func(panel ui.Element) {
		for _, size := range []string{"Small", "Medium", "Large"} {
			sel.Item(size).Label(size).Children(func() {
				components.Text(c, size).SingleLine()
			})
		}
	})
}
