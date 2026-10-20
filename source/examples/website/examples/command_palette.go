package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/command"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = registerWindow(&commandPalette{})

type commandPalette struct {
	open  bool
	state command.State
	last  string
}

var paletteCommands = []command.Item{
	{ID: "new", Title: "New document", Group: "File", Shortcut: "⌘N", Icon: lucide.FilePlus()},
	{ID: "open", Title: "Open…", Group: "File", Shortcut: "⌘O", Icon: lucide.FolderOpen()},
	{ID: "export", Title: "Export as PDF", Group: "File", Keywords: []string{"print"}, Icon: lucide.FileDown()},
	{ID: "theme", Title: "Toggle dark mode", Group: "View", Keywords: []string{"appearance"}, Icon: lucide.Moon()},
	{ID: "zoom", Title: "Zoom in", Group: "View", Shortcut: "⌘+", Icon: lucide.ZoomIn()},
	{ID: "settings", Title: "Settings", Group: "App", Shortcut: "⌘,", Keywords: []string{"preferences"}, Icon: lucide.Settings()},
}

func (p *commandPalette) View(c *ui.Context) {
	t := c.Theme()
	if c.Shortcut(ui.Cmd, ui.KeyK) {
		p.open = true
	}
	ui.Column(c).Fill().Center().Gap(t.Space(3)).Children(func() {
		if components.Button(c, "Open command palette", components.ButtonSecondary).Clicked() {
			p.open = true
		}
		text := "Or press ⌘K (Ctrl+K)"
		if p.last != "" {
			text = "Ran " + p.last
		}
		components.Text(c, text).TextColor(t.TextMuted)
	})
	command.CommandDialog(c, &p.open, &p.state, paletteCommands, command.Options{
		OnActivate: func(id string) {
			p.last, p.open = id, false
			p.state.Query = ""
		},
	}, components.DialogOptions{})
}
