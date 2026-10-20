package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/chart"
)

var _ = register(&languages{state: chart.State{Index: -1}, count: 5})

var languageShares = []chart.Slice{
	{ID: "go", Label: "Go", Value: 48}, {ID: "ts", Label: "TypeScript", Value: 31},
	{ID: "rust", Label: "Rust", Value: 18}, {ID: "swift", Label: "Swift", Value: 12},
	{ID: "py", Label: "Python", Value: 26}, {ID: "kt", Label: "Kotlin", Value: 9},
	{ID: "c", Label: "C", Value: 7},
}

type languages struct {
	state chart.State
	count int
}

func (l *languages) View(c *ui.Context) {
	ui.Column(c).Width(440).MaxWidthPercent(100).Gap(8).Children(func() {
		// A bubble's area, not its radius, follows its value.
		chart.Bubble(c, &l.state, languageShares[:l.count], chart.Options{Label: "Languages", PlotHeight: 220})
		ui.Row(c).Gap(8).Children(func() {
			if components.Button(c, "Add", components.ButtonSecondary).Disabled(l.count == len(languageShares)).Clicked() {
				l.count++
			}
			if components.Button(c, "Remove", components.ButtonGhost).Disabled(l.count == 1).Clicked() {
				l.count--
			}
		})
	})
}
