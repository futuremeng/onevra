package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&team{members: []components.AvatarGroupMember{
	{Name: "Ada Lovelace", Status: components.AvatarOnline},
	{Name: "Grace Hopper"},
	{Name: "Alan Turing"},
	{Name: "Katherine Johnson", Status: components.AvatarOnline},
	{Name: "Edsger Dijkstra"},
	{Name: "Barbara Liskov"},
}})

type team struct {
	members []components.AvatarGroupMember
}

var candidates = []string{"Margaret Hamilton", "Dennis Ritchie", "Frances Allen", "Ken Thompson"}

func (t *team) View(c *ui.Context) {
	ui.Column(c).Gap(16).AlignItems(ui.Center).Children(func() {
		// Point at the stack to loosen it; members past Max fold into "+N".
		components.AvatarGroup(c, t.members, components.AvatarGroupOptions{Max: 4})
		ui.Row(c).Gap(8).Children(func() {
			if components.Button(c, "Add", components.ButtonSecondary).Disabled(len(t.members) >= 10).Clicked() {
				t.members = append(t.members, components.AvatarGroupMember{Name: candidates[len(t.members)%len(candidates)]})
			}
			if components.Button(c, "Remove", components.ButtonGhost).Disabled(len(t.members) == 1).Clicked() {
				t.members = t.members[:len(t.members)-1]
			}
		})
	})
}
