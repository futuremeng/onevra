package examples

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = register(&shortlist{})

type destination struct {
	id, name, country string
	icon              func() *ui.SVG
	tint              ui.Color
}

var destinations = []destination{
	{"lisbon", "Lisbon", "Portugal", lucide.TramFront, ui.RGB(226, 132, 72)},
	{"kyoto", "Kyoto", "Japan", lucide.Landmark, ui.RGB(206, 86, 110)},
	{"reykjavik", "Reykjavík", "Iceland", lucide.Snowflake, ui.RGB(72, 150, 214)},
	{"cape-town", "Cape Town", "South Africa", lucide.Mountain, ui.RGB(72, 168, 128)},
}

type shortlist struct {
	cards components.CardStackState
}

func (s *shortlist) View(c *ui.Context) {
	t := c.Theme()
	items := make([]components.CardStackItem, len(destinations))
	for i, d := range destinations {
		items[i] = components.CardStackItem{ID: d.id, Label: d.name, Content: func(c *ui.Context, radius float32) {
			ui.Column(c).Fill().Children(func() {
				ui.Box(c).FillWidth().Height(t.Space(28)).Center().Radius(radius, radius, 0, 0).
					LinearGradient(ui.LinearGradient{From: d.tint.Alpha(.32), To: d.tint.Alpha(.08), Angle: 135}).
					Children(func() {
						components.Icon(c, d.icon).Size(t.Space(12), t.Space(12)).TextColor(d.tint)
					})
				ui.Column(c).FillWidth().Padding(t.Space(4)).Gap(t.Space(1)).Children(func() {
					components.Text(c, d.name).FontSize(t.Rem(1.2)).FontWeight(600)
					components.Text(c, d.country).TextColor(t.TextMuted)
				})
			})
		}}
	}
	ui.Column(c).Width(360).MaxWidthPercent(100).Gap(t.Space(4)).Children(func() {
		components.CardStack(c, &s.cards, components.CardStackOptions{
			Items:     items,
			Label:     "Destinations",
			LeftLabel: "Skip", RightLabel: "Shortlist",
		})
		var kept []string
		for _, d := range s.cards.Decisions {
			if d.Decision == components.CardRight {
				for _, dest := range destinations {
					if dest.id == d.ID {
						kept = append(kept, dest.name)
					}
				}
			}
		}
		summary := "Drag a card left or right, or use the arrow keys."
		if len(kept) > 0 {
			summary = "Shortlisted: " + strings.Join(kept, ", ")
		}
		components.Text(c, summary).TextColor(t.TextMuted).TextAlign(ui.Center).FillWidth()
	})
}
