package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = register(&advancedOptions{})

type advancedOptions struct {
	open, verbose, telemetry bool
}

func (a *advancedOptions) View(c *ui.Context) {
	ui.Column(c).Width(320).MaxWidthPercent(100).Children(func() {
		components.Collapsible(c, "Advanced", &a.open, components.CollapsibleOptions{
			ArrowPosition: components.ArrowLeft,
			Leading:       func() { components.Icon(c, lucide.SlidersHorizontal).Size(14, 14).Shrink(0) },
			Trailing:      func() { components.Badge(c, "2") },
		}, func() {
			ui.Column(c).Gap(10).Padding(8, 0, 0, 0).Children(func() {
				components.Checkbox(c, &a.verbose, "Verbose logging")
				components.Checkbox(c, &a.telemetry, "Share usage data")
			})
		}).Root.FillWidth()
	})
}
