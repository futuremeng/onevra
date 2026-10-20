package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/chart"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = registerWindow(&dashboard{
	chart:  chart.State{Index: -1},
	notify: true,
	volume: 64,
})

// dashboard is a small application window: a sidebar of pages beside an
// overview with a chart, a few settings and a release's progress.
type dashboard struct {
	page, period int
	chart        chart.State
	notify, sync bool
	volume       float64
	deploys      int
}

var dashboardPages = []struct {
	label string
	icon  func() *ui.SVG
}{
	{"Overview", lucide.LayoutDashboard},
	{"Projects", lucide.FolderKanban},
	{"Messages", lucide.MessageSquare},
	{"Settings", lucide.Settings},
}

func (d *dashboard) View(c *ui.Context) {
	t := c.Theme()
	ui.Row(c).Fill().MinWidth(0).Children(func() {
		ui.Column(c).Width(200).FillHeight().Shrink(0).Padding(t.Space(4), 0).
			Background(t.Surface).Children(func() {
			components.Sidebar(c).Children(func() {
				components.SidebarSection(c, "Workspace", func() {
					for i, p := range dashboardPages {
						if components.SidebarItem(c, p.label, p.icon(), d.page == i).Clicked() {
							d.page = i
						}
					}
				})
			})
		})
		ui.Box(c).Width(1).FillHeight().Background(t.Border.Alpha(.6))
		ui.Column(c).Grow(1).FillHeight().MinWidth(0).Padding(t.Space(5)).Gap(t.Space(4)).Children(func() {
			ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(3)).Children(func() {
				components.Text(c, dashboardPages[d.page].label).FontSize(t.Rem(1.25)).Bold().Grow(1)
				components.Tabs(c, &d.period, components.TabsSegmented, "Week", "Month")
				if components.Button(c, "Deploy", components.ButtonPrimary).Clicked() {
					d.deploys++
					c.Toast(fmt.Sprintf("Deploy #%d started", d.deploys))
				}
			})
			ui.Row(c).FillWidth().Gap(t.Space(4)).Children(func() {
				d.stat(c, "Visitors", []string{"18.2k", "74.9k"}[d.period], "+12%", components.BadgeSuccess)
				d.stat(c, "Build time", []string{"1m 42s", "1m 51s"}[d.period], "−8s", components.BadgeAccent)
				d.stat(c, "Errors", []string{"3", "17"}[d.period], "+2", components.BadgeWarning)
			})
			ui.Row(c).FillWidth().Grow(1).MinHeight(0).Gap(t.Space(4)).Children(func() {
				dashboardCard(c).Grow(1).FillHeight().MinWidth(0).Children(func() {
					components.Text(c, "Traffic").Bold()
					chart.Area(c, &d.chart, d.traffic(), chart.Options{
						Label:      "Traffic",
						PlotHeight: 190,
					})
				})
				dashboardCard(c).Width(240).FillHeight().Shrink(0).Gap(t.Space(3)).Children(func() {
					components.Text(c, "Preferences").Bold()
					components.Switch(c, &d.notify, "Notifications")
					components.Switch(c, &d.sync, "Sync on launch")
					components.Text(c, "Volume").TextColor(t.TextMuted)
					components.Slider(c, &d.volume, 0, 100).Label("Volume").FillWidth()
					ui.Spacer(c)
					components.Text(c, "Release 2.4").TextColor(t.TextMuted)
					components.Progress(c, .72)
				})
			})
		})
	})
}

func (d *dashboard) stat(c *ui.Context, label, value, change string, tone components.BadgeTone) {
	t := c.Theme()
	dashboardCard(c).Grow(1).MinWidth(0).Gap(t.Space(1.5)).Children(func() {
		components.Text(c, label).TextColor(t.TextMuted)
		ui.Row(c).AlignItems(ui.Center).Gap(t.Space(2)).Children(func() {
			components.Text(c, value).FontSize(t.Rem(1.5)).Bold()
			components.Badge(c, change, components.BadgeOptions{Tone: tone})
		})
	})
}

func (d *dashboard) traffic() chart.Data {
	if d.period == 1 {
		return chart.Data{
			Labels: []string{"W1", "W2", "W3", "W4"},
			Series: []chart.Series{
				{ID: "visits", Label: "Visits", Values: []float64{16.4, 18.9, 17.2, 22.4}},
				{ID: "signups", Label: "Sign-ups", Values: []float64{4.1, 5.3, 4.8, 6.9}},
			},
		}
	}
	return chart.Data{
		Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
		Series: []chart.Series{
			{ID: "visits", Label: "Visits", Values: []float64{2.1, 2.8, 2.4, 3.3, 3.1, 2.2, 2.3}},
			{ID: "signups", Label: "Sign-ups", Values: []float64{0.6, 0.9, 0.7, 1.2, 1.0, 0.5, 0.6}},
		},
	}
}

// dashboardCard is a raised surface for a group of content.
func dashboardCard(c *ui.Context) ui.Element {
	t := c.Theme()
	return ui.Column(c).Padding(t.Space(4)).Gap(t.Space(2)).Radius(t.Radius*1.5).
		Background(t.Surface).Border(1, t.Border.Alpha(.6))
}
