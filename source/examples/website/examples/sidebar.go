package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = registerWindow(&navigation{page: "Home"})

type navigation struct {
	page string
}

func (n *navigation) View(c *ui.Context) {
	t := c.Theme()
	item := func(label string, icon *ui.SVG) {
		if components.SidebarItem(c, label, icon, n.page == label).Clicked() {
			n.page = label
		}
	}
	ui.Row(c).Fill().Children(func() {
		ui.Column(c).Width(220).FillHeight().Padding(t.Space(4), 0).Background(t.Surface).Children(func() {
			components.Sidebar(c).Children(func() {
				components.SidebarSection(c, "Library", func() {
					item("Home", lucide.House())
					item("Recent", lucide.Clock())
					item("Shared", lucide.Users())
				})
				components.SidebarSectionWithActions(c, "Projects", func(c *ui.Context) {
					components.IconButton(c, lucide.Plus(), "New project", components.IconButtonCompact)
				}, func() {
					item("Website", lucide.Globe())
					item("Desktop app", lucide.AppWindow())
				})
			})
		})
		ui.Column(c).Grow(1).FillHeight().Center().Children(func() {
			components.Text(c, n.page).FontSize(t.Rem(1.4)).Bold()
		})
	})
}
