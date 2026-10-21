package main

import (
	"github.com/egoist/mygo/plugins/glass"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// The iOS screens follow iOS 27's Liquid Glass. Its tab bar measures, from
// UIKit, a capsule 62 points tall and 21 points from the screen's edges,
// with tabs 86 points wide; bar buttons are 44-point circles.
const (
	iosBarHeight = 62
	iosBarInset  = 21
	iosTabWidth  = 86
	iosNavHeight = 44
	iosButton    = 44
	// iosPaneInset separates floating panes, the sidebar and sheets, from
	// the screen's edges.
	iosPaneInset = 8
)

// iosTabBar reports whether the screen shows the floating tab bar: iPhone in
// portrait. iOS 27 turns it into a sidebar on iPhone in landscape, as iPad
// shows one.
func (a *app) iosTabBar() bool {
	return devices[a.device].compact(a.orientation) && !a.orientation.landscape()
}

// glass returns the material of bars, buttons and sheets. glassOpacity
// stands in for iOS 27's setting between clear and tinted glass.
func (a *app) glass(c *ui.Context) glass.Glass {
	return glass.Glass{Tint: c.Theme().Background.Alpha(float32(a.glassOpacity) * 0.9)}
}

// glassOutline is the hairline iOS 27 draws around glass, which separates
// it from plain surfaces instead of a shadow.
func glassOutline(c *ui.Context) ui.Color {
	if c.Theme().Dark {
		return ui.RGBA(255, 255, 255, 0.12)
	}
	return ui.RGBA(0, 0, 0, 0.08)
}

// glassPlatter marks the chosen tab. iOS 27's lighter dark glass marks it
// darker rather than lighter.
func glassPlatter(c *ui.Context) ui.Color {
	if c.Theme().Dark {
		return ui.RGBA(0, 0, 0, 0.3)
	}
	return ui.RGBA(0, 0, 0, 0.07)
}

// glassButton builds a circular bar button of interactive glass; prominent
// tints it with the accent, as a confirming action.
func (a *app) glassButton(c *ui.Context, icon func() *ui.SVG, label string, prominent bool) ui.Element {
	t := c.Theme()
	g, color := a.glass(c), t.Text
	if prominent {
		g.Tint, color = t.Accent, t.AccentText
	}
	g.Interactive = true
	b := ui.ButtonBase(c).Label(label).Size(iosButton, iosButton).Radius(iosButton/2).Center().
		Material(g).Border(1, glassOutline(c)).TextColor(color)
	b.Children(func() { components.Icon(c, icon).Size(20, 20) })
	return b
}

// iosScreen builds the content edge to edge under floating glass: the
// navigation bar's buttons over a hard scroll edge, and the tab bar or the
// sidebar.
func (a *app) iosScreen(c *ui.Context, d device, w, h float32, safe insets) {
	tabBar := a.iosTabBar()
	paneLeft, paneWidth := max(iosPaneInset, safe.left), float32(240)
	if !d.compact(a.orientation) {
		paneWidth = 300
	}
	left, s := float32(0), safe
	if !tabBar {
		left, s.left = paneLeft+paneWidth, iosPaneInset
	}
	split := a.tab == tabComponents && w-left >= 700
	key := a.pageKey(split)
	scroll := a.scrollState(key)
	a.trackMinimize(key, scroll)

	bottom := s.bottom
	if tabBar && (!a.keyboardShown() || a.tab == tabSearch) {
		bottom = a.barBottom(safe) + iosBarHeight
	}
	f := pageFrame{pad: insets{s.top + iosNavHeight, s.right, bottom, s.left}, scroll: scroll, split: split}
	// Large titles head roots on tall screens and shrink into the bar once
	// scrolled away.
	large := h >= 500 && !a.pushed(split)
	if large {
		f.title = a.title(split)
	}
	ui.Box(c).Absolute().Top(0).Bottom(0).Left(left).Right(0).Clip().Children(func() {
		a.page(c, f)
		// iOS 27's default scroll edge is hard: content frosts under the
		// bar, which a hairline ends.
		if scroll.Y > 1 {
			components.ScrollEdge(c, components.ScrollEdgeOptions{Style: components.ScrollEdgeHard, Extent: s.top + iosNavHeight + 4})
		}
		title := ""
		if !large || scroll.Y > 40 {
			title = a.title(split)
		}
		a.iosNavBar(c, s, a.pushed(split), title)
	})
	if tabBar {
		a.iosTabs(c, safe)
	} else {
		a.iosSidebar(c, d, paneLeft, paneWidth, safe)
	}
}

// barBottom is the tab bar's distance from the screen's bottom edge, beside
// the home indicator.
func (a *app) barBottom(safe insets) float32 {
	return max(devices[a.device].safe(a.orientation, false).bottom-13, iosPaneInset)
}

// trackMinimize shrinks the tab bar while the page scrolls down and restores
// it when it scrolls back up, as tabBarMinimizeBehavior's onScrollDown.
func (a *app) trackMinimize(key string, scroll *ui.ScrollState) {
	if key != a.minimizeKey {
		a.minimizeKey, a.lastY, a.minimized = key, scroll.Y, false
		return
	}
	switch {
	case scroll.Y <= 0:
		a.minimized = false
	case scroll.Y > a.lastY+2 && scroll.Y > 60:
		a.minimized = true
	case scroll.Y < a.lastY-2:
		a.minimized = false
	}
	a.lastY = scroll.Y
}

// iosNavBar floats the navigation bar's glass buttons and inline title over
// the content.
func (a *app) iosNavBar(c *ui.Context, s insets, back bool, title string) {
	ui.Row(c).Absolute().Top(s.top).Left(s.left + 16).Right(s.right + 16).Height(iosNavHeight).
		AlignItems(ui.Center).Gap(8).PassThrough().Children(func() {
		if back && a.glassButton(c, lucide.ChevronLeft, a.text("back"), false).Clicked() {
			a.selected = ""
		}
		if title != "" {
			ui.Text(c, title).Bold().SingleLine().Absolute().Attach(ui.AnchorCenter, ui.AnchorCenter).PassThrough()
		}
		ui.Spacer(c)
		if a.glassButton(c, lucide.Plus, a.text("sampleOpenSheet"), false).Clicked() {
			a.sheet = true
		}
	})
}

// iosTabs floats the tab bar over the content's bottom: a glass capsule of
// tabs and the prominent Search tab beside it. Choosing Search turns the bar
// into the search field, which rides above the keyboard.
func (a *app) iosTabs(c *ui.Context, safe insets) {
	bottom := a.barBottom(safe)
	left, right := safe.left+iosBarInset, safe.right+iosBarInset
	if a.tab == tabSearch {
		if a.keyboardShown() {
			bottom = safe.bottom + iosPaneInset
		}
		ui.Row(c).Key("searchbar").Absolute().Left(left).Right(right).Bottom(bottom).Height(52).Gap(10).
			Children(func() {
				last := tabs[a.lastTab]
				back := ui.ButtonBase(c).Label(a.text(last.label)).Size(52, 52).Radius(26).Center().
					Material(a.glass(c)).Border(1, glassOutline(c))
				back.Children(func() { components.Icon(c, last.icon).Size(22, 22) })
				if back.Clicked() {
					a.tab = a.lastTab
				}
				a.searchField(c)
			})
		return
	}
	if a.keyboardShown() {
		// The keyboard covers the tab bar.
		return
	}
	width := float32(3*iosTabWidth + 8)
	if a.minimized {
		width = iosBarHeight
	}
	ui.Row(c).Key("tabbar").Absolute().Left(left).Bottom(bottom).Size(width, iosBarHeight).
		Radius(iosBarHeight/2).Padding(4).Material(a.glass(c)).Border(1, glassOutline(c)).ClipX().
		Transition(ui.ElementTransition{Animation: &ui.Snappy, Size: true}).
		Children(func() {
			for i := range tabSearch {
				if !a.minimized || i == a.tab {
					a.iosTab(c, i)
				}
			}
		})
	// iOS 27 can make any tab prominent; Search is, apart from the others.
	search := ui.ButtonBase(c).Label(a.text("tabSearch")).Absolute().Right(right).Bottom(bottom).
		Size(iosBarHeight, iosBarHeight).Radius(iosBarHeight/2).Center().
		Material(a.glass(c)).Border(1, glassOutline(c))
	search.Children(func() { components.Icon(c, lucide.Search).Size(24, 24) })
	if search.Clicked() {
		a.tab = tabSearch
	}
}

// iosTab builds one tab of the bar; the minimized bar shows the current tab
// alone, without its label, and expands when it is pressed.
func (a *app) iosTab(c *ui.Context, i int) {
	t := c.Theme()
	width := float32(iosTabWidth)
	if a.minimized {
		width = iosBarHeight - 8
	}
	item := ui.ButtonBase(c).Key(tabs[i].label).Label(a.text(tabs[i].label)).Column().
		Size(width, iosBarHeight-8).Radius((iosBarHeight - 8) / 2).Center().Gap(1).Shrink(0)
	if a.tab == i {
		item.Background(glassPlatter(c)).TextColor(t.Accent)
	}
	item.Children(func() {
		components.Icon(c, tabs[i].icon).Size(24, 24)
		if !a.minimized {
			ui.Text(c, a.text(tabs[i].label)).FontSize(10).FontWeight(600).SingleLine()
		}
	})
	if item.Clicked() {
		a.choose(i)
	}
}

// choose shows tab i; choosing the current tab again returns to its root,
// and a minimized bar expands first.
func (a *app) choose(i int) {
	switch {
	case a.minimized:
		a.minimized = false
		return
	case a.tab == i:
		a.selected = ""
	}
	a.tab = i
	if i != tabSearch {
		a.lastTab = i
	}
}

// iosSidebar floats the sidebar, a glass pane listing the tabs, inside the
// safe area at the screen's leading edge.
func (a *app) iosSidebar(c *ui.Context, d device, left, width float32, safe insets) {
	t := c.Theme()
	base := d.safe(a.orientation, false)
	ui.Column(c).Key("sidebar").Absolute().Left(left).Top(max(safe.top, iosPaneInset)).
		Bottom(max(base.bottom, iosPaneInset)).Width(width).
		Radius(max(d.radius-iosPaneInset, 16)).Material(a.glass(c)).Border(1, glassOutline(c)).
		Padding(12, 10).Gap(2).Children(func() {
		ui.Text(c, a.text("title")).FontSize(22).Bold().SingleLine().Padding(4, 10, 10)
		for i, tab := range tabs {
			item := ui.ButtonBase(c).Key(tab.label).Label(a.text(tab.label)).Row().AlignItems(ui.Center).Justify(ui.Start).
				Height(iosButton).Radius(touchRadius).Padding(0, 10).Gap(10)
			if a.tab == i {
				item.Background(glassPlatter(c)).TextColor(t.Accent)
			}
			item.Children(func() {
				components.Icon(c, tab.icon).Size(20, 20)
				ui.Text(c, a.text(tab.label)).SingleLine()
			})
			if item.Clicked() {
				a.minimized = false
				a.choose(i)
			}
		}
	})
}

// searchField builds the search capsule: glass on iOS, a filled bar on
// Android. Focusing it shows the keyboard.
func (a *app) searchField(c *ui.Context) {
	t := c.Theme()
	field := ui.Row(c).Key("search").Grow(1).Height(52).Radius(26).Padding(0, 18).Gap(10).
		AlignItems(ui.Center).Bind(&a.fields[1])
	if devices[a.device].platform == platformIOS {
		field.Material(a.glass(c)).Border(1, glassOutline(c))
	} else {
		field.Height(56).Radius(28).Background(t.Surface.Mix(t.Text, 0.06))
	}
	field.Children(func() {
		components.Icon(c, lucide.Search).Size(18, 18).TextColor(t.TextMuted)
		ui.Box(c).Grow(1).Children(func() {
			if a.query == "" {
				ui.Text(c, a.text("searchPlaceholder")).TextColor(t.TextMuted).SingleLine().
					Absolute().Left(0).PassThrough()
			}
			components.Text(c, a.query, components.TextOptions{Editable: &a.query}).
				Label(a.text("tabSearch")).FillWidth()
		})
	})
}

// iosSheet floats a sheet of glass at its partial height, inset from the
// screen with corners concentric to it, with close and confirm buttons in
// glass circles. Screens with square corners keep it on the bottom edge.
func (a *app) iosSheet(c *ui.Context, d device, w float32, safe insets) {
	if !a.backdrop(c, 0.12) {
		return
	}
	t := c.Theme()
	inset, radius := float32(iosPaneInset), max(d.radius-iosPaneInset, 16)
	pad := float32(24)
	if d.radius == 0 {
		inset, pad = 0, pad+safe.bottom
	}
	width := min(w-2*inset, 520)
	sheet := ui.Column(c).Key("sheet").Absolute().Bottom(inset).Left((w-width)/2).Width(width).
		Material(a.glass(c)).Border(1, glassOutline(c)).Padding(14, 14, pad, 14).Gap(14).Transition(sheetMotion)
	if inset == 0 {
		sheet.Radius(radius, radius, 0, 0)
	} else {
		sheet.Radius(radius)
	}
	sheet.Children(func() {
		ui.Row(c).Height(iosButton).AlignItems(ui.Center).Children(func() {
			if a.glassButton(c, lucide.X, a.text("sampleCancel"), false).Clicked() {
				a.sheet = false
			}
			ui.Text(c, "Sheet").Bold().Absolute().Attach(ui.AnchorCenter, ui.AnchorCenter).PassThrough()
			ui.Spacer(c)
			if a.glassButton(c, lucide.Check, a.text("sampleDone"), true).Clicked() {
				a.sheet = false
			}
		})
		ui.Text(c, a.text("noteSheet")).TextColor(t.TextMuted).Padding(0, 6)
	})
}
