package examples

import (
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&timezone{zone: "Europe/Lisbon"})

var timezones = []string{
	"America/New_York", "America/Los_Angeles", "Asia/Shanghai", "Asia/Tokyo",
	"Europe/Berlin", "Europe/Lisbon", "Europe/London", "Pacific/Auckland",
}

type timezone struct {
	zone, query string
}

func (z *timezone) View(c *ui.Context) {
	t := c.Theme()
	sel := components.Select(c, &z.zone)
	sel.Trigger.Label("Time zone").Children(func() {
		components.Text(c, z.zone).SingleLine()
	})
	sel.Popup(func(panel ui.Element) {
		panel.MinWidth(t.Rem(16))
		sel.SearchField(&z.query).Label("Search time zones").Placeholder("Search")
		matches := 0
		for _, zone := range timezones {
			if !strings.Contains(strings.ToLower(zone), strings.ToLower(z.query)) {
				continue
			}
			matches++
			sel.Item(zone).Label(zone).Children(func() {
				components.Text(c, zone).SingleLine()
			})
		}
		if matches == 0 {
			components.Text(c, "No time zones match").Padding(t.Space(2)).TextColor(t.TextMuted)
		}
	})
}
