package main

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// Tabs of the simulated application. Search is iOS's prominent tab; Android
// opens it from the top app bar.
const (
	tabComponents = iota
	tabOverview
	tabDevice
	tabSearch
)

var tabs = []struct {
	label string
	icon  func() *ui.SVG
}{
	{"tabComponents", lucide.LayoutGrid},
	{"tabOverview", lucide.Layers},
	{"tabDevice", lucide.Smartphone},
	{"tabSearch", lucide.Search},
}

// pageFrame places a page's scrolling content among the bars around it.
type pageFrame struct {
	// pad keeps the content clear of bars and safe areas, beyond its own
	// margins.
	pad insets
	// title, when set, is the large title at the top of the content.
	title  string
	scroll *ui.ScrollState
	// split shows the open component beside the list.
	split bool
}

// screenView builds the simulated application in the conventions of the
// device's platform, then the sheet over it.
func (a *app) screenView(c *ui.Context, w, h float32, safe insets) {
	d := devices[a.device]
	size := float32(17)
	if d.platform == platformAndroid {
		size = 16
	}
	ui.Box(c).Fill().FontSize(size).Children(func() {
		if d.platform == platformIOS {
			a.iosScreen(c, d, w, h, safe)
		} else {
			a.androidScreen(c, d, w, safe)
		}
	})
	if d.platform == platformIOS {
		a.iosSheet(c, d, w, safe)
	} else {
		a.androidSheet(c, w, safe)
	}
}

// pushed reports a component open over the list on a single-pane screen.
func (a *app) pushed(split bool) bool {
	return a.tab == tabComponents && a.selected != "" && !split
}

// title returns the current page's title.
func (a *app) title(split bool) string {
	if a.pushed(split) {
		return a.selected
	}
	return a.text(tabs[a.tab].label)
}

// scrollState returns the retained scroll position of the page key.
func (a *app) scrollState(key string) *ui.ScrollState {
	if a.scrolls == nil {
		a.scrolls = map[string]*ui.ScrollState{}
	}
	s, ok := a.scrolls[key]
	if !ok {
		s = &ui.ScrollState{}
		a.scrolls[key] = s
	}
	return s
}

// page builds the current tab's content, or the open component.
func (a *app) page(c *ui.Context, f pageFrame) {
	switch {
	case a.pushed(f.split):
		a.detail(c, f)
	case a.tab == tabOverview:
		a.overview(c, f)
	case a.tab == tabDevice:
		a.deviceInfo(c, f)
	case a.tab == tabSearch:
		a.searchPage(c, f)
	case f.split:
		// The list and the open component side by side, as a split view.
		t := c.Theme()
		list, detail := f, f
		list.pad.right, detail.pad.left, detail.title = 0, 0, a.selected
		detail.scroll = a.scrollState("detail")
		ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
			ui.Box(c).Width(320+f.pad.left).FillHeight().BorderWidth(0, 1, 0, 0).BorderColor(t.Border).
				Children(func() { a.componentList(c, list) })
			ui.Box(c).Grow(1).FillHeight().Children(func() { a.detail(c, detail) })
		})
	default:
		a.componentList(c, f)
	}
}

// pageKey names the scroll position of the current page.
func (a *app) pageKey(split bool) string {
	if a.pushed(split) {
		return "detail"
	}
	return tabs[a.tab].label
}

// scrollPage is a vertical scroll view laid out within f.
func scrollPage(c *ui.Context, f pageFrame, content func()) {
	s := components.ScrollView(c, components.ScrollViewOptions{}).Fill().
		Padding(f.pad.top+12, f.pad.right+16, f.pad.bottom+24, f.pad.left+16).Gap(20).
		ScrollbarInsets(f.pad.top, 0, f.pad.bottom)
	if f.scroll != nil {
		s.TrackScroll(f.scroll)
	}
	s.Children(func() {
		if f.title != "" {
			ui.Text(c, f.title).FontSize(34).Bold().SingleLine().Margin(0, 4, -8)
		}
		content()
	})
}

// section builds a titled group of rows on a rounded surface, as iOS's inset
// grouped lists and Android's cards.
func section(c *ui.Context, title string, rows func()) {
	t := c.Theme()
	ui.Column(c).Gap(6).Children(func() {
		if title != "" {
			ui.Text(c, title).FontSize(13).TextColor(t.TextMuted).Padding(0, 4)
		}
		ui.Column(c).Radius(touchRadius).Background(t.Surface).Clip().Dividers(1, t.Border).Children(rows)
	})
}

// touchRadius rounds grouped rows.
const touchRadius = 12

// componentRow builds a catalog entry's row with its prediction on the
// device's platform; choosing it opens the entry.
func (a *app) componentRow(c *ui.Context, e entry) {
	t := c.Theme()
	row := ui.ButtonBase(c).Key(e.id).Label(e.id).
		Row().AlignItems(ui.Center).FillWidth().Height(touchTarget(a)+4).Padding(0, 12, 0, 16).Gap(8)
	if row.Pressed() || a.selected == e.id {
		row.Background(t.SurfacePressed)
	}
	row.Children(func() {
		ui.Text(c, e.id).SingleLine().Grow(1).Shrink(1)
		s := e.on(devices[a.device].platform)
		components.Badge(c, a.text(s.name()), s.badge()).AlignSelf(ui.Center)
		components.Icon(c, lucide.ChevronRight).Size(16, 16).TextColor(t.TextMuted)
	})
	if row.Clicked() {
		a.selected, a.tab = e.id, tabComponents
	}
}

// componentList builds the catalog grouped by purpose.
func (a *app) componentList(c *ui.Context, f pageFrame) {
	scrollPage(c, f, func() {
		for _, g := range groups {
			section(c, a.text(g), func() {
				for _, e := range catalog {
					if e.group == g {
						a.componentRow(c, e)
					}
				}
			})
		}
	})
}

// matches returns the entries whose API name contains the query, ignoring
// case.
func matches(query string) []entry {
	query = strings.ToLower(strings.TrimSpace(query))
	var found []entry
	for _, e := range catalog {
		if strings.Contains(strings.ToLower(e.id), query) {
			found = append(found, e)
		}
	}
	return found
}

// searchPage lists the entries matching the query. fieldOnTop puts the
// search field above them, where iPhone keeps it in the tab bar instead.
func (a *app) searchPage(c *ui.Context, f pageFrame) {
	t := c.Theme()
	scrollPage(c, f, func() {
		if devices[a.device].platform == platformAndroid || !a.iosTabBar() {
			a.searchField(c)
		}
		found := matches(a.query)
		if len(found) == 0 {
			ui.Text(c, a.text("searchEmpty")).TextColor(t.TextMuted).Margin(24, ui.Auto, 0)
			return
		}
		section(c, "", func() {
			for _, e := range found {
				a.componentRow(c, e)
			}
		})
	})
}

// detail builds the open component's predictions, note and live preview.
func (a *app) detail(c *ui.Context, f pageFrame) {
	t := c.Theme()
	e, ok := find(a.selected)
	if !ok {
		ui.Box(c).Fill().Center().Children(func() {
			ui.Text(c, a.text("emptyDetail")).TextColor(t.TextMuted)
		})
		return
	}
	scrollPage(c, f, func() {
		section(c, a.text("prediction"), func() {
			for _, p := range []struct {
				name  string
				level support
			}{{"iOS", e.ios}, {"Android", e.android}} {
				ui.Row(c).AlignItems(ui.Center).Height(touchTarget(a)).Padding(0, 16).Children(func() {
					ui.Text(c, p.name).Grow(1)
					components.Badge(c, a.text(p.level.name()), p.level.badge()).AlignSelf(ui.Center)
				})
			}
			ui.Text(c, a.text(e.note)).Padding(12, 16).TextColor(t.TextMuted)
		})
		if e.preview != nil {
			section(c, a.text("preview"), func() {
				ui.Column(c).Padding(16).Gap(12).Children(func() { e.preview(c, a) })
			})
		}
	})
}

// overview counts the predictions on each platform.
func (a *app) overview(c *ui.Context, f pageFrame) {
	t := c.Theme()
	scrollPage(c, f, func() {
		section(c, a.text("overviewCounts"), func() {
			ui.Row(c).Height(36).Padding(0, 16).TextColor(t.TextMuted).FontSize(13).Children(func() {
				ui.Spacer(c)
				ui.Text(c, "iOS").Width(64).TextAlign(ui.End)
				ui.Text(c, "Android").Width(72).TextAlign(ui.End)
			})
			for _, s := range supports {
				ui.Row(c).Height(touchTarget(a)).Padding(0, 16).Children(func() {
					ui.Box(c).Grow(1).AlignItems(ui.Start).Children(func() {
						components.Badge(c, a.text(s.name()), s.badge())
					})
					ui.Textf(c, "%d", count(platformIOS, s)).Width(64).TextAlign(ui.End)
					ui.Textf(c, "%d", count(platformAndroid, s)).Width(72).TextAlign(ui.End)
				})
			}
		})
		section(c, a.text("overviewLegend"), func() {
			for _, s := range supports {
				ui.Column(c).Padding(12, 16).Gap(4).Children(func() {
					ui.Text(c, a.text(s.name())).Bold()
					ui.Text(c, a.text(s.summary())).TextColor(t.TextMuted)
				})
			}
		})
	})
}

// deviceInfo reports the simulated screen's metrics as an application would
// read them.
func (a *app) deviceInfo(c *ui.Context, f pageFrame) {
	t := c.Theme()
	d, o := devices[a.device], a.orientation
	w, h := d.screen(o)
	s := d.safe(o, a.keyboardShown())
	unit := "pt"
	if d.platform == platformAndroid {
		unit = "dp"
	}
	class := "classCompact"
	if !d.compact(o) {
		class = "classRegular"
	}
	row := func(label, value string) {
		ui.Row(c).Height(touchTarget(a)).Padding(0, 16).Gap(12).Children(func() {
			ui.Text(c, label).Grow(1).SingleLine()
			ui.Text(c, value).TextColor(t.TextMuted).SingleLine()
		})
	}
	scrollPage(c, f, func() {
		section(c, a.text("deviceSection"), func() {
			row(a.text("deviceModel"), a.text(d.name))
			row(a.text("deviceScreen"), fmt.Sprintf("%g × %g %s", w, h, unit))
			row(a.text("deviceOrientation"), a.text(orientationNames[o]))
			row(a.text("deviceClass"), a.text(class))
		})
		section(c, a.text("safeSection"), func() {
			row(a.text("safeTop"), fmt.Sprintf("%g", s.top))
			row(a.text("safeRight"), fmt.Sprintf("%g", s.right))
			row(a.text("safeBottom"), fmt.Sprintf("%g", s.bottom))
			row(a.text("safeLeft"), fmt.Sprintf("%g", s.left))
		})
		section(c, "", func() {
			ui.Column(c).Padding(12, 16).Gap(12).Children(func() {
				components.Switch(c, &a.overlay, a.text("overlay"))
				components.Switch(c, &a.keyboard, a.text("keyboard"))
			})
		})
	})
}

// backdrop dims the screen under a sheet and dismisses it when pressed. It
// stays built so the dimming fades by color rather than by an opacity layer
// over the whole screen, and reports whether the sheet is open.
func (a *app) backdrop(c *ui.Context, dim float32) bool {
	b := ui.Box(c).Key("backdrop").Absolute().Fill()
	if !a.sheet {
		dim = 0
	}
	b.Background(ui.RGBA(0, 0, 0, b.AnimateTo("dim", dim, ui.Smooth)))
	if !a.sheet {
		b.PassThrough()
		return false
	}
	if b.Clicked() {
		a.sheet = false
	}
	return true
}

// sheetMotion slides the sheet up from below the screen and back.
var sheetMotion = ui.ElementTransition{
	Animation: &ui.Smooth,
	Enter:     &ui.Motion{Y: 480, Opacity: 1},
	Exit:      &ui.Motion{Y: 480, Opacity: 1},
}
