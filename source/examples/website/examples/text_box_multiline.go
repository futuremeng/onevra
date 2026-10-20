package examples

import (
	"fmt"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&feedback{})

type feedback struct {
	message string
}

func (f *feedback) View(c *ui.Context) {
	ui.Column(c).Width(360).MaxWidthPercent(100).Gap(8).Children(func() {
		components.TextBox(c, &f.message, components.TextBoxOptions{Multiline: true}).
			Label("Feedback").Placeholder("Tell us what you think…").Lines(4, 8).FillWidth()
		components.Text(c, fmt.Sprintf("%d characters", utf8.RuneCountInString(f.message))).
			TextColor(c.Theme().TextMuted)
	})
}
