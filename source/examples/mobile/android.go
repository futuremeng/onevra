package main

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

// androidScreen builds Material's opaque top app bar and navigation bar
// around the content, with the floating action button. Widths from 600 dp
// show the list beside the open component.
func (a *app) androidScreen(c *ui.Context, d device, w float32, safe insets) {
	split := a.tab == tabComponents && !d.compact(a.orientation)
	ui.Column(c).Fill().Children(func() {
		a.androidTopBar(c, safe, split)
		ui.Box(c).Grow(1).MinHeight(0).Children(func() {
			a.page(c, pageFrame{pad: insets{right: safe.right, left: safe.left},
				scroll: a.scrollState(a.pageKey(split)), split: split})
		})
		if a.keyboardShown() {
			// The keyboard covers the navigation bar; content ends at its top.
			ui.Box(c).Height(safe.bottom)
		} else {
			a.androidNavBar(c, safe)
		}
	})
	if a.tab == tabComponents && (a.selected == "" || split) && !a.keyboardShown() {
		a.fab(c, safe)
	}
}

// androidTopBar builds the top app bar, with a back arrow on pushed screens
// and search among its actions.
func (a *app) androidTopBar(c *ui.Context, safe insets, split bool) {
	t := c.Theme()
	pushed := a.pushed(split) || a.tab == tabSearch
	bar := ui.Row(c).FillWidth().Height(safe.top+64).Padding(safe.top, safe.right+8, 0, safe.left+8).
		Gap(4).Background(t.Background)
	if pushed || split {
		bar.BorderWidth(0, 0, 1, 0).BorderColor(t.Border)
	}
	bar.Children(func() {
		if pushed {
			if components.IconButton(c, lucide.ArrowLeft(), a.text("back"), components.IconButtonPlain).Clicked() {
				if a.tab == tabSearch {
					a.tab = a.lastTab
				} else {
					a.selected = ""
				}
			}
		} else {
			ui.Box(c).Width(8)
		}
		ui.Text(c, a.title(split)).FontSize(22).SingleLine().Grow(1)
		if a.tab != tabSearch &&
			components.IconButton(c, lucide.Search(), a.text("tabSearch"), components.IconButtonPlain).Clicked() {
			a.tab = tabSearch
		}
	})
}

// androidNavBar builds the navigation bar, padded by the bottom safe area,
// with a pill behind the current destination's icon.
func (a *app) androidNavBar(c *ui.Context, safe insets) {
	t := c.Theme()
	ui.Row(c).FillWidth().Padding(0, safe.right, safe.bottom, safe.left).Background(t.Surface).
		BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
		for i := range tabSearch {
			tab := tabs[i]
			item := ui.ButtonBase(c).Key(tab.label).Label(a.text(tab.label)).Column().Grow(1).Basis(0).
				Height(80).Center().Gap(4)
			color := t.TextMuted
			if a.tab == i {
				color = t.Accent
			}
			item.TextColor(color).Children(func() {
				icon := ui.Box(c).Size(64, 32).Radius(16).Center().
					Children(func() { components.Icon(c, tab.icon).Size(24, 24) })
				if a.tab == i {
					icon.Background(t.Accent.Alpha(0.18))
				}
				ui.Text(c, a.text(tab.label)).FontSize(12).SingleLine()
			})
			if item.Clicked() {
				a.choose(i)
			}
		}
	})
}

// fab builds the floating action button above the navigation bar.
func (a *app) fab(c *ui.Context, safe insets) {
	t := c.Theme()
	fab := ui.ButtonBase(c).Label(a.text("sampleOpenSheet")).Absolute().
		Right(safe.right+16).Bottom(safe.bottom+80+16).Size(56, 56).Radius(16).Center().
		Background(t.Accent).TextColor(t.AccentText).Shadow(0, 4, 12, 0, ui.RGBA(0, 0, 0, 0.25))
	fab.Children(func() { components.Icon(c, lucide.Plus).Size(24, 24) })
	if fab.Clicked() {
		a.sheet = true
	}
}

// androidSheet builds a modal bottom sheet on the bottom edge, with a drag
// handle and its content padded by the bottom safe area.
func (a *app) androidSheet(c *ui.Context, w float32, safe insets) {
	if !a.backdrop(c, 0.32) {
		return
	}
	t := c.Theme()
	width := min(w, 640)
	ui.Column(c).Key("sheet").Absolute().Bottom(0).Left((w-width)/2).Width(width).
		Radius(28, 28, 0, 0).Background(t.Surface).
		Padding(8, safe.right+20, safe.bottom+20, safe.left+20).Gap(16).
		Transition(sheetMotion).
		Children(func() {
			ui.Box(c).Size(32, 4).Radius(2).Background(t.TextMuted.Alpha(0.5)).Margin(14, ui.Auto, 8)
			ui.Text(c, "Sheet").FontSize(22)
			ui.Text(c, a.text("noteSheet")).TextColor(t.TextMuted)
			if components.Button(c, a.text("sampleDone"), components.ButtonPrimary).Height(touchTarget(a)).Clicked() {
				a.sheet = false
			}
		})
}
