package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components/chart"
)

var _ = register(&revenue{state: chart.State{Index: -1}})

type revenue struct {
	state chart.State
}

func (r *revenue) View(c *ui.Context) {
	data := chart.Data{
		Labels: []string{"Q1", "Q2", "Q3", "Q4"},
		Series: []chart.Series{
			{ID: "apps", Label: "Apps", Values: []float64{42, 48, 51, 63}},
			{ID: "services", Label: "Services", Values: []float64{28, 31, 37, 40}},
			{ID: "hardware", Label: "Hardware", Values: []float64{18, 15, 22, 30}},
		},
	}
	ui.Column(c).Width(460).MaxWidthPercent(100).Children(func() {
		// The legend hides and shows series; totals stack.
		chart.StackedBar(c, &r.state, data, chart.Options{Label: "Revenue", PlotHeight: 200})
	})
}
