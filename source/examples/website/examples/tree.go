package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = register(&projectTree{source: true, selected: "main.go"})

type projectTree struct {
	source, docs bool
	selected     string
}

func (p *projectTree) View(c *ui.Context) {
	folder := func(c *ui.Context, size float32) {
		components.Icon(c, lucide.Folder).Size(size, size).TextColor(c.Theme().Accent)
	}
	file := func(name string) {
		if components.TreeItem(c, name, nil, p.selected == name, nil).Clicked() {
			p.selected = name
		}
	}
	ui.Column(c).Width(260).Children(func() {
		components.Tree(c, func() {
			if components.TreeItemWithLeading(c, "source", &p.source, p.selected == "source", folder, func() {
				file("main.go")
				file("app.go")
				file("go.mod")
			}).Clicked() {
				p.selected = "source"
			}
			if components.TreeItemWithLeading(c, "docs", &p.docs, p.selected == "docs", folder, func() {
				file("build.md")
			}).Clicked() {
				p.selected = "docs"
			}
			file("README.md")
		}).Label("Project").FillWidth()
	})
}
