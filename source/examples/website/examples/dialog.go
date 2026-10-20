package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = registerWindow(&renameDialog{name: "Quarterly report"})

type renameDialog struct {
	open, discard bool
	name, draft   string
}

func (d *renameDialog) View(c *ui.Context) {
	t := c.Theme()
	ui.Column(c).Fill().Center().Gap(t.Space(3)).Children(func() {
		components.Text(c, d.name).FontSize(t.Rem(1.2)).Bold()
		if components.Button(c, "Rename…", components.ButtonSecondary).Clicked() {
			d.draft, d.open = d.name, true
		}
	})
	components.Dialog(c, &d.open, components.DialogOptions{Label: "Rename"}, func(panel ui.Element) {
		panel.Width(340)
		components.Text(c, "Rename document").Bold()
		components.TextBox(c, &d.draft, components.TextBoxOptions{}).Label("Name").FillWidth()
		ui.Row(c).FillWidth().Gap(8).Justify(ui.End).Children(func() {
			if components.Button(c, "Cancel", components.ButtonGhost).Clicked() {
				// A nested AlertDialog opens over this one.
				d.discard = d.draft != d.name
				d.open = d.discard
			}
			if components.Button(c, "Save", components.ButtonPrimary).Clicked() {
				d.name, d.open = d.draft, false
			}
		})
		if components.AlertDialog(c, &d.discard, components.AlertDialogOptions{
			Title:        "Discard the new name?",
			Message:      "The document keeps its current name.",
			ConfirmLabel: "Discard",
			CancelLabel:  "Keep editing",
		}) == components.DialogConfirmed {
			d.open = false
		}
	})
}
