package examples

import (
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/theme"
)

// Every example builds in both themes without the website.
func TestExamplesBuild(t *testing.T) {
	if len(Names()) == 0 {
		t.Fatal("no examples are registered")
	}
	for _, name := range Names() {
		e, _, _ := Lookup(name)
		for _, dark := range []bool{false, true} {
			tt := ui.NewTester(func(c *ui.Context) {
				c.SetTheme(theme.Palette(dark))
				e.View(c)
			}, 720, 480)
			tt.Frame()
		}
	}
}
