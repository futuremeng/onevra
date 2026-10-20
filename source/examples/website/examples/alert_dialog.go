package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = registerWindow(&confirmDelete{status: "3 files selected"})

type confirmDelete struct {
	open   bool
	status string
}

func (d *confirmDelete) View(c *ui.Context) {
	ui.Column(c).Fill().Center().Gap(12).Children(func() {
		components.Text(c, d.status).TextColor(c.Theme().TextMuted)
		if components.Button(c, "Delete files", components.ButtonSecondary).Clicked() {
			d.open = true
		}
	})
	// The result arrives in the frame of the choice; the call never blocks.
	switch components.AlertDialog(c, &d.open, components.AlertDialogOptions{
		Title:        "Delete 3 files?",
		Message:      "This cannot be undone.",
		ConfirmLabel: "Delete",
		CancelLabel:  "Cancel",
	}) {
	case components.DialogConfirmed:
		d.status = "Deleted"
	case components.DialogCancelled:
		d.status = "Kept the files"
	}
}
