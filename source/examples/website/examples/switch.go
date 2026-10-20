package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&wirelessSettings{wifi: true})

type wirelessSettings struct {
	wifi, bluetooth bool
}

func (s *wirelessSettings) View(c *ui.Context) {
	ui.Column(c).Gap(c.Theme().Space(3)).Children(func() {
		components.Switch(c, &s.wifi, "Wi-Fi")
		components.Switch(c, &s.bluetooth, "Bluetooth")
	})
}
