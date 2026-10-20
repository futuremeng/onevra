package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&fileActions{})

type fileActions struct {
	last string
}

func (f *fileActions) View(c *ui.Context) {
	actions := func(m *ui.Menu) {
		for _, action := range []string{"Open", "Duplicate", "Rename"} {
			if m.Item(action).Chosen() {
				f.last = action
			}
		}
		m.Separator()
		if m.Item("Delete").Chosen() {
			f.last = "Delete"
		}
	}
	ui.Column(c).Gap(16).AlignItems(ui.Center).Children(func() {
		ui.Row(c).Gap(8).Children(func() {
			components.Button(c, "Actions", components.ButtonSecondary).Menu(actions)
			components.Button(c, "Right-click me", components.ButtonOutline).ContextMenu(actions)
		})
		text := "Choose an action"
		if f.last != "" {
			text = "Chose " + f.last
		}
		components.Text(c, text).TextColor(c.Theme().TextMuted)
	})
}
