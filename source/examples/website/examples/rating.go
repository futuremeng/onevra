package examples

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&review{stars: 4})

type review struct {
	stars int
}

func (r *review) View(c *ui.Context) {
	ui.Column(c).Gap(8).AlignItems(ui.Center).Children(func() {
		components.Rating(c, &r.stars, 5).Label("Rating")
		text := "Not rated"
		if r.stars > 0 {
			text = fmt.Sprintf("%d of 5 stars", r.stars)
		}
		components.Text(c, text).TextColor(c.Theme().TextMuted)
	})
}
