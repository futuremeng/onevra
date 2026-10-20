//go:build !js

package main

import (
	"log"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/examples/website/examples"
)

// start shows one example in a window, closing the app with it.
func start(name string) {
	e, window, ok := examples.Lookup(name)
	if !ok {
		log.Fatalf("no example named %q; choose one of %v", name, examples.Names())
	}
	w := mygo.NewWindow(mygo.WindowOptions{
		Title:   name,
		Width:   640,
		Height:  400,
		Content: ui.View(content(e, window)),
	})
	w.OnClosed(func() { mygo.App.Quit() })
}
