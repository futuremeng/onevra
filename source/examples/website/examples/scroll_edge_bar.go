package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = registerWindow(&photoFeed{})

type photoFeed struct {
	style int
}

var feedStyles = []struct {
	name  string
	style components.ScrollEdgeStyle
}{
	{"Glass", components.ScrollEdgeGlass},
	{"Soft", components.ScrollEdgeSoft},
	{"Hard", components.ScrollEdgeHard},
	{"Blur", components.ScrollEdgeBlur},
}

var feedColors = []ui.Color{
	ui.RGB(236, 120, 82), ui.RGB(79, 124, 255), ui.RGB(72, 168, 128),
	ui.RGB(226, 178, 60), ui.RGB(206, 96, 176), ui.RGB(66, 176, 176),
}

func (f *photoFeed) View(c *ui.Context) {
	options := components.ScrollEdgeBarOptions{Style: feedStyles[f.style].style}
	// Over glass, buttons stay plain; the other styles leave the bar
	// transparent, so the buttons bring their own face.
	face := components.IconButtonToolbar
	if options.Style == components.ScrollEdgeGlass {
		face = components.IconButtonPlain
	}
	inset := components.ScrollEdgeBarExtent(c, options)
	ui.Box(c).Fill().Clip().Children(func() {
		components.ScrollView(c, components.ScrollViewOptions{}).Fill().
			Padding(inset+8, 16, 16).ScrollbarInsets(inset, 0, 0).Gap(10).Children(func() {
			for i := range 16 {
				hue := feedColors[i%len(feedColors)]
				ui.Row(c).Key(i).Gap(12).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(64, 44).Radius(8).Shrink(0).
						LinearGradient(ui.LinearGradient{From: hue, To: hue.Alpha(.4), Angle: 135})
					components.Text(c, fmt.Sprintf("Photo %d", i+1))
				})
			}
		})
		components.ScrollEdgeBar(c, options, func() {
			components.IconButton(c, lucide.ChevronLeft(), "Back", face)
			ui.Row(c).Grow(1).Justify(ui.Center).Children(func() {
				names := make([]string, len(feedStyles))
				for i, s := range feedStyles {
					names[i] = s.name
				}
				components.Tabs(c, &f.style, components.TabsSegmented, names...)
			})
			components.IconButton(c, lucide.Share(), "Share", face)
		}).Label("Toolbar")
	})
}
