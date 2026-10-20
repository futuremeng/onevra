package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&release{})

type release struct {
	reads int
}

func (r *release) View(c *ui.Context) {
	t := c.Theme()
	components.RichText(c,
		ui.Span{Text: "Onevra 2.4 ", Weight: 600},
		ui.Span{Text: "adds "},
		ui.Span{Text: "TagInput", Background: t.SurfaceHover, Radius: 4},
		ui.Span{Text: ", faster tables and "},
		ui.Span{Text: "fewer surprises", Italic: true, Underline: true},
		ui.Span{Text: ". "},
	).Children(func() {
		if components.Button(c, "Read the notes", components.ButtonLink).Clicked() {
			r.reads++
		}
		if r.reads > 0 {
			components.Text(c, fmt.Sprintf(" (opened %d times)", r.reads)).TextColor(t.TextMuted)
		}
	}).MaxWidth(360)
}
