package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components/chart"
)

var _ = register(&trafficSources{state: chart.State{Index: -1}})

type trafficSources struct {
	state chart.State
}

func (t *trafficSources) View(c *ui.Context) {
	slices := []chart.Slice{
		{ID: "organic", Label: "Organic", Value: 41260},
		{ID: "direct", Label: "Direct", Value: 24610},
		{ID: "paid", Label: "Paid", Value: 15830},
		{ID: "social", Label: "Social", Value: 8420},
		{ID: "referral", Label: "Referral", Value: 5310},
	}
	ui.Column(c).Width(480).MaxWidthPercent(100).Children(func() {
		chart.SplitBar(c, &t.state, slices, chart.Options{
			Label:       "Traffic sources",
			Description: "Sessions this month",
		})
	})
}
