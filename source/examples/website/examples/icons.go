package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = register(&iconSet{})

type iconSet struct{}

func (iconSet) View(c *ui.Context) {
	t := c.Theme()
	icons := []struct {
		name string
		svg  func() *ui.SVG
	}{
		{"House", lucide.House}, {"Search", lucide.Search}, {"Settings", lucide.Settings},
		{"Folder", lucide.Folder}, {"Download", lucide.Download}, {"Bell", lucide.Bell},
		{"Heart", lucide.Heart}, {"Sparkles", lucide.Sparkles},
	}
	ui.Row(c).Wrap().Gap(t.Space(2)).Justify(ui.Center).Children(func() {
		for _, icon := range icons {
			ui.Column(c).Key(icon.name).Width(84).Padding(t.Space(3)).Gap(t.Space(2)).
				AlignItems(ui.Center).Radius(t.Radius).Background(t.Surface).
				Border(1, t.Border.Alpha(.6)).Children(func() {
				components.Icon(c, icon.svg).Size(24, 24)
				components.Text(c, icon.name).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
			})
		}
	})
}
