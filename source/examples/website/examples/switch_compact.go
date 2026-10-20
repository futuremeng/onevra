package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&channelList{on: []bool{true, false, true}})

type channelList struct {
	on []bool
}

func (l *channelList) View(c *ui.Context) {
	t := c.Theme()
	compact := components.SwitchOptions{Compact: true, HideLabel: true}
	ui.Column(c).Width(280).Gap(t.Space(1)).Children(func() {
		for i, name := range []string{"Email", "Push", "SMS"} {
			ui.Row(c).Key(name).FillWidth().AlignItems(ui.Center).Padding(t.Space(1.5), 0).Children(func() {
				components.Text(c, name).Grow(1)
				// HideLabel keeps the label as the switch's accessible name.
				components.Switch(c, &l.on[i], name, compact)
			})
		}
	})
}
