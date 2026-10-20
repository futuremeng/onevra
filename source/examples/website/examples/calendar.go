package examples

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&dueDate{day: time.Date(2026, time.October, 20, 9, 30, 0, 0, time.Local)})

type dueDate struct {
	day time.Time
}

func (d *dueDate) View(c *ui.Context) {
	ui.Column(c).Gap(12).AlignItems(ui.Center).Children(func() {
		components.Calendar(c, &d.day)
		components.Text(c, d.day.Format("Monday, January 2")).TextColor(c.Theme().TextMuted)
	})
}
