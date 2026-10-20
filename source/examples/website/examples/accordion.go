package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&faq{open: []bool{true, false, false}})

type faq struct {
	open []bool
}

func (f *faq) View(c *ui.Context) {
	questions := []struct{ q, a string }{
		{"Does Onevra need a web view?", "No. MyGo draws every control itself, on the GPU."},
		{"Which platforms does it run on?", "macOS, Windows and Linux, from the same Go code."},
		{"Can I use MyGo widgets too?", "Yes. Onevra builds on MyGo and works beside it."},
	}
	ui.Column(c).Width(440).MaxWidthPercent(100).Children(func() {
		components.Accordion(c, func(group *components.AccordionGroup) {
			for i, item := range questions {
				if group.Item(item.q, &f.open[i], func() {
					components.Text(c, item.a).TextColor(c.Theme().TextMuted).FillWidth()
				}).Changed() && f.open[i] {
					// Keep one section open at a time.
					for j := range f.open {
						f.open[j] = j == i
					}
				}
			}
		})
	})
}
