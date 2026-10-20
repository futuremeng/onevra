package examples

import (
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&profileCard{})

type profileCard struct {
	loading bool
	started time.Time
}

func (p *profileCard) View(c *ui.Context) {
	t := c.Theme()
	// Pretend to load for a moment; a real application would load elsewhere
	// and clear loading on the UI thread.
	if p.loading {
		if wait := 1500*time.Millisecond - c.Now().Sub(p.started); wait > 0 {
			c.After(wait)
		} else {
			p.loading = false
		}
	}
	ui.Column(c).Width(340).MaxWidthPercent(100).Gap(16).Children(func() {
		components.SkeletonContent(c, p.loading, components.SkeletonOptions{Lines: 2, Avatar: true}, func() {
			ui.Row(c).FillWidth().Gap(12).AlignItems(ui.Center).Children(func() {
				components.Avatar(c, "Grace Hopper", components.AvatarOptions{})
				ui.Column(c).Grow(1).Children(func() {
					components.Text(c, "Grace Hopper").Bold()
					components.Text(c, "Rear admiral, computer scientist").TextColor(t.TextMuted)
				})
			})
		})
		if components.Button(c, "Reload", components.ButtonSecondary).Disabled(p.loading).Clicked() {
			p.loading, p.started = true, c.Now()
		}
	})
}
