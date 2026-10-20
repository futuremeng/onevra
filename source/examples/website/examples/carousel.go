package examples

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = register(&trips{})

var tripSlides = []struct {
	title, caption string
	icon           func() *ui.SVG
	color          ui.Color
}{
	{"Mountains", "Three days above the clouds", lucide.Mountain, ui.RGB(88, 128, 196)},
	{"Forest", "Quiet trails and cabins", lucide.Trees, ui.RGB(72, 168, 128)},
	{"Desert", "Dunes at sunrise", lucide.Sun, ui.RGB(226, 150, 60)},
	{"Camping", "Under the stars", lucide.TentTree, ui.RGB(196, 108, 148)},
	{"Harbour", "Boats and fresh fish", lucide.Sailboat, ui.RGB(66, 160, 190)},
}

type trips struct {
	state components.CarouselState
}

func (s *trips) View(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).FillWidth().Children(func() {
		components.Carousel(c, &s.state, components.CarouselOptions{
			Label:    "Trips",
			Count:    len(tripSlides),
			Interval: 5 * time.Second, // adds a play control
			Slide: func(c *ui.Context, i int) {
				slide := tripSlides[i]
				ui.Column(c).FillWidth().Height(200).Radius(t.Radius*2).Padding(t.Space(5)).Gap(t.Space(1)).
					LinearGradient(ui.LinearGradient{From: slide.color.Alpha(.32), To: slide.color.Alpha(.1), Angle: 135}).
					Border(1, slide.color.Alpha(.35)).Children(func() {
					components.Icon(c, slide.icon).Size(t.Space(8), t.Space(8)).TextColor(slide.color)
					ui.Spacer(c)
					components.Text(c, slide.title).FontWeight(600).FontSize(t.Rem(1.25))
					components.Text(c, slide.caption).TextColor(t.TextMuted)
				})
			},
		})
	})
}
