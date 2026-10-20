package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&shipping{method: "standard"})

type shipping struct {
	method string
}

func (s *shipping) View(c *ui.Context) {
	ui.Column(c).Gap(12).Children(func() {
		components.Radio(c, &s.method, "standard", "Standard, 3–5 days")
		components.Radio(c, &s.method, "express", "Express, next day")
		components.Radio(c, &s.method, "pickup", "Pick up in store")
	})
}
