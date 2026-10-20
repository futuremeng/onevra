//go:build js

package main

import (
	"errors"
	"sync"
	"syscall/js"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
	"github.com/oneo-me/onevra/source/examples/website/examples"
)

var mounted = struct {
	sync.Mutex
	next    int
	windows map[int]*mygo.Window
}{windows: map[int]*mygo.Window{}}

// start installs globalThis.onevraExamples for the page, which opens an
// example's window in the element it gives the ID "mygo" (see the website's
// example host), and closes it when the element leaves the page.
func start(string) {
	api := js.Global().Get("Object").New()
	api.Set("mount", js.FuncOf(func(_ js.Value, args []js.Value) any {
		name := args[0].String()
		return promise(func() (any, error) {
			e, window, ok := examples.Lookup(name)
			if !ok {
				return nil, errors.New("no example named " + name)
			}
			// Hidden, then shown without the keyboard: an example takes it
			// when clicked, not when it scrolls into view.
			w := mygo.NewWindow(mygo.WindowOptions{Hidden: true, Content: ui.View(content(e, window))})
			w.ShowInactive()
			mounted.Lock()
			defer mounted.Unlock()
			mounted.next++
			mounted.windows[mounted.next] = w
			return mounted.next, nil
		})
	}))
	api.Set("unmount", js.FuncOf(func(_ js.Value, args []js.Value) any {
		id := args[0].Int()
		mounted.Lock()
		w := mounted.windows[id]
		delete(mounted.windows, id)
		mounted.Unlock()
		if w != nil {
			go w.Close()
		}
		return nil
	}))
	api.Set("setLocale", js.FuncOf(func(_ js.Value, args []js.Value) any {
		tag := args[0].String()
		go mygo.App.SetLocale(tag)
		return nil
	}))
	js.Global().Set("onevraExamples", api)
	js.Global().Call("dispatchEvent", js.Global().Get("Event").New("onevra-examples-ready"))
}

// writeManifest reads source files, which browsers do not have; the website
// runs it natively.
func writeManifest() error { return errors.New("the manifest is written natively") }

// promise runs fn outside the browser's call, which cannot wait for the main
// loop it would block, and settles a promise with its result.
func promise(fn func() (any, error)) js.Value {
	var executor js.Func
	executor = js.FuncOf(func(_ js.Value, args []js.Value) any {
		resolve, reject := args[0], args[1]
		go func() {
			defer executor.Release()
			v, err := fn()
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New(err.Error()))
				return
			}
			resolve.Invoke(v)
		}()
		return nil
	})
	return js.Global().Get("Promise").New(executor)
}
