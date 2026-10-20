package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&people{photo: decode(catJPEG)})

type people struct {
	photo *ui.Bitmap // decoded once, when the application starts
}

func (p *people) View(c *ui.Context) {
	ui.Row(c).Gap(16).AlignItems(ui.Center).Children(func() {
		components.Avatar(c, "Ada Lovelace", components.AvatarOptions{})
		components.Avatar(c, "Mochi", components.AvatarOptions{Image: p.photo, Status: components.AvatarOnline})
		components.Avatar(c, "Grace Hopper", components.AvatarOptions{Status: components.AvatarOffline})
		components.Avatar(c, "Alan Turing", components.AvatarOptions{Size: 56})
	})
}
