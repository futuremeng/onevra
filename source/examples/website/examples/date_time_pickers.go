package examples

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&meeting{at: time.Date(2026, time.October, 20, 14, 30, 0, 0, time.Local)})

type meeting struct {
	at time.Time
}

func (m *meeting) View(c *ui.Context) {
	ui.Column(c).Gap(12).AlignItems(ui.Center).Children(func() {
		ui.Row(c).Gap(8).Children(func() {
			// Both edit the same time.Time: the date keeps the time of day,
			// and the time keeps the date.
			components.CalendarPicker(c, &m.at).Label("Date")
			components.TimePicker(c, &m.at).Label("Time")
		})
		components.Text(c, m.at.Format("Mon, Jan 2 at 15:04")).TextColor(c.Theme().TextMuted)
	})
}
