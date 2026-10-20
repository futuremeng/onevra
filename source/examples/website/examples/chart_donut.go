package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components/chart"
)

var _ = register(&timeSpent{state: chart.State{Index: -1}})

type timeSpent struct {
	state chart.State
}

func (t *timeSpent) View(c *ui.Context) {
	slices := []chart.Slice{
		{ID: "design", Label: "Design", Value: 38},
		{ID: "build", Label: "Development", Value: 52},
		{ID: "research", Label: "Research", Value: 21},
		{ID: "planning", Label: "Planning", Value: 14},
	}
	ui.Column(c).Width(420).MaxWidthPercent(100).Children(func() {
		chart.Donut(c, &t.state, slices, chart.Options{Label: "Time spent", PlotHeight: 220})
	})
}
