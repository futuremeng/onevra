package examples

import (
	"fmt"
	"slices"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = registerWindow(&logTable{})

type logTable struct {
	table    components.TableVirtualizedState
	sort     ui.SortOrder
	rows     []int // the order rows show in, re-sorted only when the user sorts
	selected int
}

// logLatency is a stand-in for a row's data.
func logLatency(i int) int { return i*7919%977 + 1 }

func (l *logTable) View(c *ui.Context) {
	if l.rows == nil {
		l.rows = make([]int, 100000)
		for i := range l.rows {
			l.rows[i] = i
		}
		l.table.Sort = &l.sort
		l.table.List.Selected = &l.selected
		l.table.List.Key = func(row int) any { return l.rows[row] }
		l.table.List.KeepOffset = true
		l.table.OnSort = func() {
			slices.SortStableFunc(l.rows, func(a, b int) int {
				order := a - b
				if l.sort.Column == "latency" {
					order = logLatency(a) - logLatency(b)
				}
				if l.sort.Descending {
					return -order
				}
				return order
			})
		}
	}
	columns := []ui.TableColumn{
		{Title: "Request", ID: "request", Sortable: true, MinWidth: 160},
		{Title: "Status", ID: "status", Width: 90},
		{Title: "Latency", ID: "latency", Width: 100, Align: ui.End, Sortable: true},
	}
	components.TableVirtualized(c, &l.table, columns, len(l.rows), func(row, col int) {
		i := l.rows[row]
		switch col {
		case 0:
			components.Text(c, fmt.Sprintf("GET /api/items/%d", i+1)).SingleLine()
		case 1:
			components.Text(c, "200 OK").SingleLine().TextColor(c.Theme().Success)
		case 2:
			components.Text(c, fmt.Sprintf("%d ms", logLatency(i))).SingleLine()
		}
	}).Fill().Label("Requests")
}
