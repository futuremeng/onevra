package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&photoFits{photo: decode(catJPEG)})

type photoFits struct {
	photo *ui.Bitmap // decoded once, when the application starts
}

func (p *photoFits) View(c *ui.Context) {
	t := c.Theme()
	sample := func(name string, image ui.Element) {
		ui.Column(c).Key(name).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Background(t.Surface).Radius(t.Radius).Children(func() { image.Label(name) })
			components.Text(c, name).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
		})
	}
	ui.Row(c).Gap(16).Wrap().Justify(ui.Center).Children(func() {
		sample("Contain", components.Image(c, p.photo).Size(140, 100))
		sample("Cover", components.Image(c, p.photo).Size(140, 100).Fit(ui.Cover))
		sample("Rounded", components.Image(c, p.photo).Size(100, 100).Fit(ui.Cover).Radius(50))
		sample("Grayscale", components.Image(c, p.photo).Size(140, 100).Fit(ui.Cover).Grayscale())
	})
}
