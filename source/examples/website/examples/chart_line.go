package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/chart"
)

var _ = register(&temperatures{state: chart.State{Index: -1}})

type temperatures struct {
	state chart.State
	round int
}

func (t *temperatures) View(c *ui.Context) {
	// Updating the data animates the lines to their new values.
	values := [][]float64{{12, 14, 13, 17, 19, 18, 21}, {15, 13, 16, 14, 18, 22, 20}}
	data := chart.Data{
		Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
		Series: []chart.Series{
			{ID: "lisbon", Label: "Lisbon", Values: values[t.round%2]},
			{ID: "oslo", Label: "Oslo", Values: []float64{4, 6, 5, 7, 3, 5, 8}},
		},
	}
	ui.Column(c).Width(460).MaxWidthPercent(100).Gap(8).Children(func() {
		chart.Line(c, &t.state, data, chart.Options{
			Label:       "Temperature",
			PlotHeight:  200,
			FormatValue: func(v float64) string { return fmt.Sprintf("%.0f °C", v) },
		})
		if components.Button(c, "Next week", components.ButtonSecondary).Clicked() {
			t.round++
		}
	})
}
