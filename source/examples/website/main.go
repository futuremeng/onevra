// The website's live examples, the files of package examples, built for
// browsers by website/scripts/examples.ts. In a browser one instance serves
// every example of a page, each in a window of its own; natively it shows one
// example, for checking it without the website:
//
//	go run ./examples/website -example switch
//	go run ./examples/website -manifest
package main

import (
	"flag"
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/components"
	"github.com/oneo-me/onevra/source/examples/website/examples"
	"github.com/oneo-me/onevra/source/theme"
)

func main() {
	manifest := flag.Bool("manifest", false, "print the examples and their source as JSON")
	example := flag.String("example", "", "example to show natively, such as switch")
	language := flag.String("lang", "en", "interface language: en or zh-Hans")
	flag.Parse()

	if *manifest {
		if err := writeManifest(); err != nil {
			log.Fatal(err)
		}
		return
	}

	mygo.App.SetName("Onevra examples")
	mygo.App.WhenReady(func() {
		mygo.App.SetLocale(*language)
		start(*example)
	})
	// A page closes the windows of the examples it no longer shows, which
	// leaves the app running for the next page's.
	mygo.App.OnWindowAllClosed(func() {})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

// content builds an example in Onevra's theme, centered in its window with a
// margin unless it lays out the window itself.
func content(e examples.Example, window bool) func(*ui.Context) {
	return func(c *ui.Context) {
		theme.Apply(c)
		t := c.Theme()
		ui.Box(c).Fill().Background(t.Background).Children(func() {
			if window {
				e.View(c)
				return
			}
			ui.Column(c).Fill().Center().Padding(t.Space(6)).Children(func() {
				e.View(c)
			})
		})
		// A Shell would host these; an example's window has none.
		components.Toasts(c)
	}
}
