package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/icons/lucide"
)

var _ = register(&mailboxTabs{})

type mailboxTabs struct {
	selected int
}

func (s *mailboxTabs) View(c *ui.Context) {
	icons := []func() *ui.SVG{lucide.Inbox, lucide.Star, lucide.Archive}
	counts := []string{"12", "3", ""}
	components.TabsWith(c, &s.selected, components.TabsOptions{
		Variant: components.TabsSegmented,
		Leading: func(c *ui.Context, i int, size float32) {
			components.Icon(c, icons[i]).Size(size, size).Shrink(0)
		},
		Trailing: func(c *ui.Context, i int, _ float32) {
			if counts[i] != "" {
				components.Text(c, counts[i]).TextColor(c.Theme().TextMuted)
			}
		},
		Menu: func(i int) func(m *ui.Menu) {
			return func(m *ui.Menu) {
				if m.Item("Open").Chosen() {
					s.selected = i
				}
			}
		},
	}, "Inbox", "Starred", "Archive")
}
