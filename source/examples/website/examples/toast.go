package examples

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = registerWindow(&notifications{items: 3})

type notifications struct {
	items int
}

func (n *notifications) View(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Fill().Center().Gap(t.Space(3)).Children(func() {
		ui.Row(c).Gap(8).Wrap().Justify(ui.Center).Children(func() {
			if components.Button(c, "Saved", components.ButtonSecondary).Clicked() {
				c.AddToast(ui.Toast{Title: "Changes saved", Type: components.ToastSuccess})
			}
			if components.Button(c, "Offline", components.ButtonSecondary).Clicked() {
				c.AddToast(ui.Toast{Title: "You are offline", Description: "Changes sync when you reconnect.",
					Type: components.ToastWarning})
			}
			if components.Button(c, "Delete item", components.ButtonSecondary).Disabled(n.items == 0).Clicked() {
				n.items--
				c.AddToast(ui.Toast{Title: "Item deleted", Action: "Undo", Timeout: 8 * time.Second,
					OnAction: func() { n.items++ }})
			}
		})
		components.Text(c, fmt.Sprintf("%d items", n.items)).TextColor(t.TextMuted)
	})
}
