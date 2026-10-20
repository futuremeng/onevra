package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&palette{})

type palette struct{}

// View shows the frame theme's colors, which follow the light or dark
// appearance: switch the site's theme to see the other palette.
func (palette) View(c *ui.Context) {
	t := c.Theme()
	swatches := []struct {
		name  string
		color ui.Color
	}{
		{"Background", t.Background}, {"Surface", t.Surface}, {"Border", t.Border},
		{"Text", t.Text}, {"TextMuted", t.TextMuted}, {"Accent", t.Accent},
		{"Success", t.Success}, {"Warning", t.Warning}, {"Danger", t.Danger},
	}
	ui.Row(c).Wrap().Gap(t.Space(3)).Justify(ui.Center).Children(func() {
		for _, s := range swatches {
			ui.Column(c).Key(s.name).Width(96).Gap(t.Space(1.5)).Children(func() {
				ui.Box(c).FillWidth().Height(48).Radius(t.Radius).
					Background(s.color).Border(1, t.Border)
				components.Text(c, s.name).FontSize(t.Rem(0.8))
			})
		}
	})
}
