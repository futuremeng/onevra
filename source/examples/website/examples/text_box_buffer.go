package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&notes{buffer: ui.NewTextBuffer("A buffer edits in place,\nso long documents stay fast.")})

type notes struct {
	// buffer is kept across frames; editing does not copy the whole text.
	buffer *ui.TextBuffer
}

func (n *notes) View(c *ui.Context) {
	ui.Column(c).Width(360).MaxWidthPercent(100).Gap(8).Children(func() {
		components.TextBoxBuffer(c, n.buffer, components.TextBoxOptions{Multiline: true}).
			Label("Notes").Lines(4, 8).FillWidth()
		components.Text(c, fmt.Sprintf("%d bytes", n.buffer.Len())).TextColor(c.Theme().TextMuted)
	})
}
