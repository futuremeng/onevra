// Mobile is a pre-research gallery that predicts how Onevra's components
// would carry over to iOS and Android. It simulates phones and tablets on the
// desktop: a simple device body with its cutout, the system's safe areas,
// portrait and both landscape orientations, and a software keyboard. The
// screen runs a small application in each platform's conventions whose
// catalog rates every component: iOS 27's Liquid Glass, with a floating tab
// bar, glass bar buttons over a hard scroll edge, a sidebar in landscape and
// floating sheets, on MyGo's glass plugin; Material on Android. A slider
// stands in for iOS 27's setting between clear and tinted glass. MyGo does not run on mobile systems today,
// so nothing here calls a mobile API; the metrics are approximations of
// common devices.
package main

import (
	"fmt"
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/examples/mobile/locales"
	"github.com/oneo-me/onevra/source/shells/plain"
	"github.com/oneo-me/onevra/source/theme"
)

// app keeps the simulation and the simulated application's state in memory.
type app struct {
	device      int
	orientation orientation
	// overlay tints the safe-area insets; keyboard shows the software
	// keyboard, which focusing a text field also does (typing).
	overlay, keyboard, typing bool
	// glassOpacity goes from clear to tinted Liquid Glass, as iOS 27's
	// setting does.
	glassOpacity float64

	tab, lastTab int
	selected     string // the catalog entry shown, or ""
	query        string
	sheet        bool
	demo         demo
	// fields are the text fields that bring up the keyboard: the TextBox
	// preview and the search field.
	fields  [2]ui.Handle
	scrolls map[string]*ui.ScrollState
	// minimized shrinks the iOS tab bar while minimizeKey's page scrolls
	// down from lastY.
	minimized   bool
	minimizeKey string
	lastY       float32

	// scale is the frame's last fitted scale, reported beside the controls.
	scale  float32
	locale string
}

// panelWidth is the width of the controls beside the device.
const panelWidth = 300

var orientationNames = []string{"portrait", "landscapeLeft", "landscapeRight"}

// text returns the static message key in the frame's language.
func (a *app) text(key string) string { return locales.Catalog.Format(a.locale, key, nil) }

func (a *app) keyboardShown() bool { return a.keyboard || a.typing }

func (a *app) view(c *ui.Context) {
	a.locale = c.Locale()
	a.typing = c.FocusWithin(a.fields[0]) || c.FocusWithin(a.fields[1])
	t := c.Theme()
	// Command+Left and Command+Right turn the device, as in Xcode's Simulator.
	if c.Shortcut(ui.Cmd, ui.KeyLeft) {
		a.rotate(-1)
	}
	if c.Shortcut(ui.Cmd, ui.KeyRight) {
		a.rotate(1)
	}
	w, h := c.Size()
	ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
		stage := ui.Box(c).Grow(1).FillHeight().ClipX().ClipY()
		if !t.Dark {
			stage.Background(t.Surface)
		}
		stage.Children(func() {
			a.frameView(c, w-panelWidth-48, h-plain.HeaderHeight(c)-48, a.screenView)
		})
		ui.Box(c).Width(panelWidth).FillHeight().BorderWidth(0, 0, 0, 1).BorderColor(t.Border).
			Children(func() { a.controls(c) })
	})
}

// rotate turns the device a quarter, counterclockwise for negative turns.
func (a *app) rotate(turn int) {
	order := []orientation{landscapeLeft, portrait, landscapeRight}
	i := 1
	for j, o := range order {
		if o == a.orientation {
			i = j
		}
	}
	a.orientation = order[min(max(i+turn, 0), len(order)-1)]
}

// controls builds the simulation's settings beside the device.
func (a *app) controls(c *ui.Context) {
	t := c.Theme()
	heading := func(key string) {
		ui.Text(c, a.text(key)).FontSize(t.FontSize*0.9).Bold().TextColor(t.TextMuted).Margin(8, 0, 0, 0)
	}
	components.ScrollView(c, components.ScrollViewOptions{}).Fill().Padding(16, 20).Gap(10).Children(func() {
		heading("deviceSection")
		for i, d := range devices {
			if components.Radio(c, &a.device, i, a.text(d.name)).Changed() {
				a.selected, a.sheet, a.minimized = "", false, false
			}
		}
		heading("deviceOrientation")
		choice := int(a.orientation)
		labels := make([]string, len(orientationNames))
		for i, key := range orientationNames {
			labels[i] = a.text(key)
		}
		if components.Tabs(c, &choice, components.TabsSegmented, labels...).List.Changed() {
			a.orientation = orientation(choice)
		}
		heading("simulate")
		components.Switch(c, &a.overlay, a.text("overlay"))
		components.Switch(c, &a.keyboard, a.text("keyboard"))
		heading("glassOpacity")
		components.Slider(c, &a.glassOpacity, 0, 1)
		heading("scale")
		ui.Text(c, fmt.Sprintf("%.0f%%", a.scale*100)).TextColor(t.TextMuted)
	})
}

func main() {
	mygo.App.SetName("Onevra Mobile")
	mygo.App.WhenReady(func() {
		a := &app{demo: demo{value: 0.4, on: true}}
		win := plain.NewShell(plain.ShellOptions{
			Title: locales.In(mygo.App.Locale()).Title(),
			Width: 1120, Height: 880, MinWidth: 820, MinHeight: 740,
			StateKey: "mobile",
			Theme:    theme.Onevra,
			Content:  a.view,
		}).Window()
		win.OnClosed(func() { mygo.App.Quit() })
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
