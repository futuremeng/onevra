package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&paint{color: ui.RGB(110, 90, 230)})

type paint struct {
	color ui.Color
}

func (p *paint) View(c *ui.Context) {
	components.ColorView(c, &p.color).Label("Color")
}
