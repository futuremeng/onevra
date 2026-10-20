package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&dimensions{quantity: 2, width: 1280, opacity: 0.8})

type dimensions struct {
	quantity, width, opacity float64
}

func (d *dimensions) View(c *ui.Context) {
	ui.Column(c).Width(240).Gap(12).Children(func() {
		components.NumberInput(c, &d.quantity, 1, 99, 1).Label("Quantity").FillWidth()
		// A leading text drags the value: press on "W" and move sideways.
		components.NumberInput(c, &d.width, 0, 4096, 1, components.NumberInputOptions{
			LeadingText: "W",
		}).Label("Width").FillWidth()
		components.NumberInput(c, &d.opacity, 0, 1, 0.05, components.NumberInputOptions{
			HideStepper: true,
		}).Label("Opacity").FillWidth()
	})
}
