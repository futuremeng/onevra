// Package examples holds the live examples of the website's documentation.
// Each file declares one example, registered under the file's name, and the
// documentation shows it running beside the file's source. Write an example
// as an application would: state in a value, built by its View each frame.
package examples

import (
	_ "embed"
	"path"
	"runtime"
	"sort"
	"strings"

	"github.com/egoist/mygo/ui"
)

// Example is an application value whose View builds its content each frame.
type Example interface {
	View(c *ui.Context)
}

type entry struct {
	example Example
	window  bool
}

var registry = map[string]entry{}

// register adds the example under the name of the file declaring it: tabs.go
// is "tabs" and button_variants.go is "button-variants". The host centers it
// in its window, with a margin. Each example file calls register or
// registerWindow once, as var _ = register(&value{}); the website leaves that
// line out of the source it shows.
func register(e Example) Example { return add(e, false) }

// registerWindow adds an example that lays out its whole window itself, as an
// application's content does, from edge to edge.
func registerWindow(e Example) Example { return add(e, true) }

func add(e Example, window bool) Example {
	_, file, _, ok := runtime.Caller(2)
	if !ok {
		panic("examples: register cannot find its caller")
	}
	name := strings.ReplaceAll(strings.TrimSuffix(path.Base(file), ".go"), "_", "-")
	if _, dup := registry[name]; dup {
		panic("examples: " + name + " is registered twice")
	}
	registry[name] = entry{e, window}
	return e
}

// cat.jpeg is the photo the image examples show, decoded once as they
// register, since an application decodes images outside its builders.
//
//go:embed cat.jpeg
var catJPEG []byte

func decode(data []byte) *ui.Bitmap {
	b, err := ui.DecodeBitmap(data)
	if err != nil {
		panic(err)
	}
	return b
}

// Lookup returns the example of a name, and whether it fills its window.
func Lookup(name string) (e Example, window, ok bool) {
	r, ok := registry[name]
	return r.example, r.window, ok
}

// Names lists the examples in order.
func Names() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
