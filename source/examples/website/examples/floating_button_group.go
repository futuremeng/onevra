package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = register(&boardToolbar{})

type boardToolbar struct {
	edits, undone int
	presenting    bool
}

func (b *boardToolbar) View(c *ui.Context) {
	present := "Present"
	if b.presenting {
		present = "Stop"
	}
	items := []components.FloatingButtonGroupItem{
		{ID: "undo", Label: "Undo", Icon: lucide.Undo2, IconOnly: true, Shortcut: "⌘Z", Disabled: b.edits == 0},
		{ID: "redo", Label: "Redo", Icon: lucide.Redo2, IconOnly: true, Shortcut: "⇧⌘Z", Disabled: b.undone == 0},
		{Separator: true},
		{ID: "comment", Label: "Comment", Icon: lucide.MessageCircle},
		{Separator: true},
		{ID: "present", Label: present, Icon: lucide.Play, Toggle: true, On: b.presenting,
			ReserveLabels: []string{"Present", "Stop"}},
	}
	group := components.FloatingButtonGroup(c, "Board actions",
		components.FloatingButtonGroupRaised, ui.Horizontal, items)
	switch group.Clicked() {
	case "comment":
		b.edits, b.undone = b.edits+1, 0
	case "undo":
		b.edits, b.undone = b.edits-1, b.undone+1
	case "redo":
		b.edits, b.undone = b.edits+1, b.undone-1
	case "present":
		b.presenting = !b.presenting
	}
}
