package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components/chart"
)

var _ = register(&weeklyHours{state: chart.State{Index: -1}})

type weeklyHours struct {
	state chart.State // Index -1 starts with nothing selected
}

func (w *weeklyHours) View(c *ui.Context) {
	data := chart.Data{
		Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
		Series: []chart.Series{
			{ID: "this", Label: "This week", Values: []float64{38, 52, 46, 68, 61, 82}},
			{ID: "last", Label: "Last week", Values: []float64{26, 34, 28, 48, 43, 57}},
		},
	}
	ui.Column(c).Width(460).MaxWidthPercent(100).Children(func() {
		chart.Bar(c, &w.state, data, chart.Options{
			Label:       "Focus hours",
			Description: "Use the arrow keys to explore",
			PlotHeight:  200,
			FormatValue: func(v float64) string { return fmt.Sprintf("%.0f h", v) },
		})
	})
}
