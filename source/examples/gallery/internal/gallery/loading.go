package gallery

import (
	"github.com/egoist/mygo/plugins/loading"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

// exampleLoading shows every indicator of MyGo's loading plugin under its
// API name. The titlebar pauses them (s.on) or draws them in the accent
// color (s.other); each moves only while it shows.
func exampleLoading(c *ui.Context, s *demoState) {
	t := c.Theme()
	o := loading.Options{Size: 28, Paused: s.on}
	if s.other {
		o.Color = t.Accent
	}
	ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).Children(func() {
		for _, spinner := range loading.All() {
			ui.Column(c).Key(spinner.Name).Size(t.Space(28), t.Space(24)).Radius(t.Radius*1.5).Border(1, t.Border).
				AlignItems(ui.Center).Justify(ui.Center).Gap(t.Space(3)).Children(func() {
				spinner.New(c, o).Label(spinner.Name)
				components.Text(c, spinner.Name).FontSize(t.Rem(.75)).TextColor(t.TextMuted)
			})
		}
	})
}

// loadingToolbar pauses the indicators and tints them from the titlebar.
func loadingToolbar(c *ui.Context, s *demoState) {
	components.ToggleButton(c, &s.on, tr(c, "loadingPause")).KeepFocus()
	components.ToggleButton(c, &s.other, tr(c, "accent")).KeepFocus()
}
