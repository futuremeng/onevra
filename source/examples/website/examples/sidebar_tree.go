package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = registerWindow(&explorer{open: map[string]*bool{"source": new(true)}})

type projectFile struct {
	name     string
	children []projectFile
}

var projectFiles = []projectFile{
	{"source", []projectFile{
		{"components", []projectFile{{"button.go", nil}, {"tree.go", nil}}},
		{"main.go", nil},
		{"go.mod", nil},
	}},
	{"docs", []projectFile{{"build.md", nil}}},
	{"README.md", nil},
}

type explorer struct {
	open     map[string]*bool // by folder name
	selected string
}

// folder returns whether a folder is expanded, as a pointer the tree updates.
func (e *explorer) folder(name string) *bool {
	if e.open[name] == nil {
		e.open[name] = new(bool)
	}
	return e.open[name]
}

func (e *explorer) View(c *ui.Context) {
	t := c.Theme()
	var files func([]projectFile)
	files = func(nodes []projectFile) {
		for _, node := range nodes {
			var open *bool
			var children func()
			if node.children != nil {
				open = e.folder(node.name)
				children = func() { files(node.children) }
			}
			if components.TreeItem(c, node.name, open, e.selected == node.name, children).Clicked() {
				e.selected = node.name
			}
		}
	}
	ui.Row(c).Fill().Children(func() {
		ui.Column(c).Width(240).FillHeight().Padding(t.Space(4), 0).Background(t.Surface).Children(func() {
			components.Sidebar(c).Children(func() {
				components.SidebarSection(c, "Explorer", func() {
					components.Tree(c, func() { files(projectFiles) }).Label("Explorer")
				})
			})
		})
		ui.Column(c).Grow(1).FillHeight().Center().Children(func() {
			components.Text(c, e.selected).TextColor(t.TextMuted)
		})
	})
}
