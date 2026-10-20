package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&exportFormat{format: "png"})

type exportFormat struct {
	format string
}

func (e *exportFormat) View(c *ui.Context) {
	components.NativeSelect(c, &e.format, []components.NativeSelectOption[string]{
		{Value: "png", Label: "PNG image"},
		{Value: "jpg", Label: "JPEG image"},
		{Value: "pdf", Label: "PDF document"},
		{Value: "svg", Label: "SVG (coming soon)", Disabled: true},
	}).Trigger.Label("Export format")
}
