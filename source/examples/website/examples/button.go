package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&buttonVariants{})

type buttonVariants struct {
	last string
}

func (b *buttonVariants) View(c *ui.Context) {
	ui.Column(c).Gap(16).AlignItems(ui.Center).Children(func() {
		ui.Row(c).Wrap().Gap(8).Justify(ui.Center).Children(func() {
			for _, v := range []struct {
				label   string
				variant components.ButtonVariant
			}{
				{"Primary", components.ButtonPrimary},
				{"Secondary", components.ButtonSecondary},
				{"Outline", components.ButtonOutline},
				{"Ghost", components.ButtonGhost},
				{"Link", components.ButtonLink},
			} {
				ui.Box(c).Key(v.label).Children(func() {
					if components.Button(c, v.label, v.variant).Clicked() {
						b.last = v.label
					}
				})
			}
			components.Button(c, "Disabled", components.ButtonPrimary).Disabled(true)
		})
		if b.last != "" {
			components.Text(c, "Clicked "+b.last).TextColor(c.Theme().TextMuted)
		}
	})
}
