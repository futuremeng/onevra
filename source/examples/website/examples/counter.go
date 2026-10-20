package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&counter{})

type counter struct {
	count int
}

func (s *counter) View(c *ui.Context) {
	ui.Row(c).Gap(c.Theme().Space(3)).Children(func() {
		if components.Button(c, "Add one", components.ButtonPrimary).Clicked() {
			s.count++
		}
		if components.Button(c, "Reset", components.ButtonSecondary).Clicked() {
			s.count = 0
		}
		components.Text(c, fmt.Sprintf("Clicked %d times", s.count))
	})
}
