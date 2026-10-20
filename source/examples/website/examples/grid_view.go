package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = registerWindow(&fileGrid{})

type fileGrid struct {
	grid     components.GridState
	selected int
}

func (g *fileGrid) View(c *ui.Context) {
	t := c.Theme()
	g.grid.Selected = &g.selected
	g.grid.Label = func(i int) string { return fmt.Sprintf("Photo %d", i+1) }
	// Cells are at least 96 points wide; the grid fits as many per row as it can.
	components.GridView(c, &g.grid, 500, 96, 88, func(i int) {
		ui.Column(c).Fill().Center().Gap(t.Space(1)).Children(func() {
			components.Icon(c, lucide.Image).Size(28, 28).TextColor(t.Accent)
			components.Text(c, fmt.Sprintf("Photo %d", i+1)).FontSize(t.Rem(0.8))
		})
	}).Fill().Label("Photos").Padding(t.Space(3))
}
