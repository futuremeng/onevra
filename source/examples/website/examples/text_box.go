package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = register(&accountForm{name: "Ada Lovelace", id: "ADA-1815"})

type accountForm struct {
	name, query, price, id string
}

func (f *accountForm) View(c *ui.Context) {
	ui.Column(c).Width(320).MaxWidthPercent(100).Gap(12).Children(func() {
		components.TextBox(c, &f.name, components.TextBoxOptions{}).
			Label("Name").Placeholder("Your name").FillWidth()
		components.TextBox(c, &f.query, components.TextBoxOptions{LeadingIcon: lucide.Search()}).
			Label("Search").Placeholder("Search").FillWidth()
		components.TextBox(c, &f.price, components.TextBoxOptions{LeadingText: "$"}).
			Label("Price").Placeholder("0.00").FillWidth()
		components.TextBox(c, &f.id, components.TextBoxOptions{ReadOnly: true}).
			Label("Account ID").FillWidth()
	})
}
