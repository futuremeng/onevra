package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&todoList{items: []todo{
	{"Sketch the settings window", true},
	{"Wire up the sidebar", false},
	{"Package for macOS", false},
}})

type todo struct {
	title string
	done  bool
}

type todoList struct {
	draft string
	items []todo
}

func (l *todoList) View(c *ui.Context) {
	ui.Column(c).Width(340).Gap(12).Children(func() {
		ui.Row(c).FillWidth().Gap(8).Children(func() {
			input := components.TextBox(c, &l.draft, components.TextBoxOptions{}).
				Placeholder("Add a task").Label("New task").Grow(1)
			add := components.Button(c, "Add", components.ButtonPrimary)
			if (input.Submitted() || add.Clicked()) && l.draft != "" {
				l.items = append(l.items, todo{title: l.draft})
				l.draft = ""
			}
		})
		left := 0
		for i := range l.items {
			ui.Box(c).Key(i).Children(func() {
				components.Checkbox(c, &l.items[i].done, l.items[i].title)
			})
			if !l.items[i].done {
				left++
			}
		}
		components.Text(c, fmt.Sprintf("%d of %d left", left, len(l.items))).
			TextColor(c.Theme().TextMuted)
	})
}
