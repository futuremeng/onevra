package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&checkout{current: 1})

type checkout struct {
	current int
}

func (o *checkout) View(c *ui.Context) {
	steps := []components.StepperStep{
		{Label: "Cart", Description: "3 items"},
		{Label: "Shipping", Description: "Address and method"},
		{Label: "Payment", Description: "Card or transfer"},
		{Label: "Review"},
	}
	ui.Column(c).Width(520).MaxWidthPercent(100).Gap(20).Children(func() {
		stepper := components.Stepper(c, steps, o.current, components.StepperOptions{Selectable: true})
		stepper.FillWidth()
		if i, ok := stepper.Chosen(); ok {
			o.current = i
		}
		ui.Row(c).Gap(8).Justify(ui.Center).FillWidth().Children(func() {
			if components.Button(c, "Back", components.ButtonGhost).Disabled(o.current == 0).Clicked() {
				o.current--
			}
			if components.Button(c, "Continue", components.ButtonPrimary).Disabled(o.current == len(steps)).Clicked() {
				o.current++
			}
		})
	})
}
