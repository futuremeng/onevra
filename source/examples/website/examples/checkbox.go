package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&signup{newsletter: true})

type signup struct {
	terms, newsletter bool
}

func (s *signup) View(c *ui.Context) {
	ui.Column(c).Gap(12).Children(func() {
		components.Checkbox(c, &s.terms, "I accept the terms")
		components.Checkbox(c, &s.newsletter, "Send me product news")
		components.Button(c, "Create account", components.ButtonPrimary).Disabled(!s.terms)
	})
}
