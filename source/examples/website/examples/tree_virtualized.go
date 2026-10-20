package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&largeProject{})

// assets holds 10,000 files, of which only the visible rows are built.
var largeProjectFiles = func() []string {
	files := make([]string, 10000)
	for i := range files {
		files[i] = fmt.Sprintf("image-%05d.png", i+1)
	}
	return files
}()

type largeProject struct {
	tree components.TreeVirtualizedState[string]
}

func (p *largeProject) View(c *ui.Context) {
	components.TreeVirtualized(c, &p.tree, []string{"assets", "README.md"},
		func(name string) []string {
			if name == "assets" {
				return largeProjectFiles
			}
			return nil // a leaf
		},
		func(name string) {
			components.Text(c, name).SingleLine()
		},
	).Width(300).Height(260).Label("Large project").Border(1, c.Theme().Border).Radius(c.Theme().Radius)
}
