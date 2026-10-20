package examples

import (
	"fmt"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&fileTable{})

type tableFile struct {
	name, kind string
	size       int
}

var tableFiles = []tableFile{
	{"main.go", "Go source", 12}, {"go.mod", "Module", 1},
	{"README.md", "Markdown", 4}, {"logo.png", "Image", 86}, {"LICENSE", "Text", 2},
}

type fileTable struct {
	table    components.TableState
	sort     ui.SortOrder
	selected string
}

func (f *fileTable) View(c *ui.Context) {
	// Column IDs stay stable; the application sorts its own rows.
	columns := []ui.TableColumn{
		{Title: "Name", ID: "name", Sortable: true, MinWidth: 140},
		{Title: "Kind", ID: "kind", Width: 110},
		{Title: "Size", ID: "size", Width: 80, Align: ui.End, Sortable: true},
	}
	f.table.Sort = &f.sort
	files := slices.Clone(tableFiles)
	slices.SortStableFunc(files, func(a, b tableFile) int {
		order := 0
		switch f.sort.Column {
		case "name":
			order = strings.Compare(a.name, b.name)
		case "size":
			order = a.size - b.size
		}
		if f.sort.Descending {
			return -order
		}
		return order
	})
	components.Table(c, &f.table, columns, func() {
		for _, file := range files {
			if components.TableRow(c, file.name, file.name, f.selected == file.name, func(col int) {
				switch col {
				case 0:
					components.Text(c, file.name).SingleLine()
				case 1:
					components.Text(c, file.kind).SingleLine().TextColor(c.Theme().TextMuted)
				case 2:
					components.Text(c, fmt.Sprintf("%d KB", file.size)).SingleLine()
				}
			}).Clicked() {
				f.selected = file.name
			}
		}
	}).Label("Files").Width(420).MaxWidthPercent(100)
}
