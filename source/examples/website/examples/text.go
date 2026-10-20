package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&textStyles{title: "Untitled project"})

type textStyles struct {
	title string
}

func (s *textStyles) View(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Width(380).MaxWidthPercent(100).Gap(10).Children(func() {
		// Editable text edits in place, without a field around it.
		components.Text(c, s.title, components.TextOptions{Editable: &s.title}).
			Label("Project name").FontSize(t.Rem(1.4)).Bold()
		components.Text(c, "Click the title above to rename the project.").TextColor(t.TextMuted)
		components.Text(c, "Long text wraps within the available width, including in narrow windows.").
			FillWidth()
		components.Text(c, "A single line truncates when it does not fit, rather than wrapping.").
			SingleLine().Ellipsis("…").FillWidth()
	})
}
