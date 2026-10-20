package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components/chart"
)

var _ = register(&traffic{state: chart.State{Index: -1}})

type traffic struct {
	state chart.State
}

func (t *traffic) View(c *ui.Context) {
	data := chart.Data{
		Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"},
		Series: []chart.Series{
			{ID: "visits", Label: "Visits", Values: []float64{120, 180, 150, 240, 210, 300}},
			{ID: "signups", Label: "Sign-ups", Values: []float64{30, 45, 38, 70, 64, 92}},
		},
	}
	ui.Column(c).Width(460).MaxWidthPercent(100).Children(func() {
		chart.Area(c, &t.state, data, chart.Options{Label: "Traffic", PlotHeight: 200})
	})
}
