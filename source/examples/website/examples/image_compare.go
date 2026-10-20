package examples

import (
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
)

var _ = register(&beforeAfter{photo: decode(catJPEG), position: 50})

type beforeAfter struct {
	photo    *ui.Bitmap // decoded once, when the application starts
	position float32    // the share, from 0 to 100, that "before" covers
}

func (b *beforeAfter) View(c *ui.Context) {
	side := func(gray bool) func() {
		return func() {
			image := components.Image(c, b.photo).Fill().Fit(ui.Cover)
			if gray {
				image.Grayscale()
			}
		}
	}
	components.ImageCompare(c, &b.position, side(true), side(false), components.ImageCompareOptions{
		BeforeLabel: "Before",
		AfterLabel:  "After",
	}).Label("Color restoration").Width(420).MaxWidthPercent(100)
}
