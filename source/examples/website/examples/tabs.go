package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&tabVariants{})

type tabVariants struct {
	underline, segmented, highlight int
}

func (s *tabVariants) View(c *ui.Context) {
	labels := []string{"General", "Appearance", "Advanced"}
	ui.Column(c).Gap(20).AlignItems(ui.Center).Children(func() {
		components.Tabs(c, &s.underline, components.TabsUnderline, labels...)
		components.Tabs(c, &s.segmented, components.TabsSegmented, labels...)
		components.Tabs(c, &s.highlight, components.TabsHighlight, labels...)
	})
}
