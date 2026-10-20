package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = registerWindow(&editorLayout{files: 180, console: 30})

type editorLayout struct {
	files   float32 // the file list's width in DIPs
	console float32 // the console's share of the height, in percent
}

func (e *editorLayout) View(c *ui.Context) {
	t := c.Theme()
	pane := func(title, body string, background ui.Color) func() {
		return func() {
			ui.Column(c).Fill().Padding(t.Space(4)).Gap(t.Space(1)).Background(background).Children(func() {
				components.Text(c, title).Bold()
				components.Text(c, body).TextColor(t.TextMuted)
			})
		}
	}
	editor := func() {
		components.SplitPercent(c, &e.console,
			pane("Console", fmt.Sprintf("%.0f%% of the height", e.console), t.Surface),
			pane("Editor", "Drag the dividers, or focus one and use the arrow keys.", t.Background),
			components.SplitOptions{Orientation: components.Vertical, Default: 30}).
			Label("Console divider").Fill()
	}
	components.Split(c, &e.files,
		pane("Files", fmt.Sprintf("%.0f points wide", e.files), t.Surface),
		editor,
		components.SplitOptions{Default: 180, FirstMin: 120}).
		Label("Sidebar divider").Fill()
}
