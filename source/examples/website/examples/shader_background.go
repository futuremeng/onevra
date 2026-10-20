package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/components/shaders"
)

var _ = registerWindow(&welcomePage{intensity: 1})

type welcomePage struct {
	intensity float64
	still     bool
}

func (w *welcomePage) View(c *ui.Context) {
	t := c.Theme()
	// Move the pointer over the field, and click to send a ripple.
	shaders.SignalField(c, shaders.Options{
		Intensity: float32(w.intensity),
		Still:     w.still,
	}).Column().Children(func() {
		ui.Column(c).Fill().Center().Gap(t.Space(3)).Children(func() {
			components.Text(c, "Welcome back").FontSize(t.Rem(1.6)).Bold()
			components.Text(c, "Nothing open yet. Start a project or open one.").TextColor(t.TextMuted)
			ui.Row(c).Width(260).Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), 0).Children(func() {
				components.Text(c, "Intensity").TextColor(t.TextMuted)
				components.Slider(c, &w.intensity, 0.2, 2).Label("Intensity").Grow(1)
			})
			components.Switch(c, &w.still, "Still")
		})
	})
}
