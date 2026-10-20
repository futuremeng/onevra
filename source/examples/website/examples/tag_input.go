package examples

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&labels{tags: []string{"design", "go", "desktop"}})

type labels struct {
	tags []string
}

func (l *labels) View(c *ui.Context) {
	ui.Column(c).Width(380).MaxWidthPercent(100).Gap(8).Children(func() {
		tags := components.TagInput(c, &l.tags).Label("Labels").Placeholder("Add a label")
		tags.FillWidth()
		if tags.Changed() {
			// Applications normalize tags after a change.
			for i, tag := range l.tags {
				l.tags[i] = strings.ToLower(tag)
			}
		}
		components.Text(c, "Type and press Enter or a comma. Backspace removes the last label.").
			TextColor(c.Theme().TextMuted)
	})
}
