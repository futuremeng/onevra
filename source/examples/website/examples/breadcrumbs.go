package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&folderPath{path: []string{"Home", "Projects", "Onevra", "Website", "Docs"}})

type folderPath struct {
	path   []string
	chosen int
}

func (f *folderPath) View(c *ui.Context) {
	ui.Column(c).Width(420).MaxWidthPercent(100).Gap(16).Children(func() {
		if components.Breadcrumbs(c, f.path, &f.chosen).Changed() {
			// Ancestors navigate: trim the path to the chosen folder.
			f.path = f.path[:f.chosen+1]
		}
		if len(f.path) < 5 && components.Button(c, "Reset path", components.ButtonOutline).Clicked() {
			f.path = []string{"Home", "Projects", "Onevra", "Website", "Docs"}
		}
	})
}
