//go:build !js

package main

import (
	"strings"
	"testing"
)

func TestShownSourceLeavesOutPackageAndRegistration(t *testing.T) {
	src := `package examples

import "github.com/egoist/mygo/ui"

// The registration's comment goes with it.
var _ = register(&greeting{})

type greeting struct{}

func (g *greeting) View(c *ui.Context) {}
`
	got, err := shownSource("greeting.go", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	want := "import \"github.com/egoist/mygo/ui\"\n\ntype greeting struct{}\n\nfunc (g *greeting) View(c *ui.Context) {}\n"
	if got != want {
		t.Errorf("shown source =\n%s\nwant\n%s", got, want)
	}
	if _, err := shownSource("none.go", []byte("package examples\n")); err == nil {
		t.Error("a file without a registration was accepted")
	}
}

func TestManifestMatchesRegistry(t *testing.T) {
	entries, err := readExamples(sourceDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if strings.Contains(e.Code, "register") || strings.HasPrefix(e.Code, "package") {
			t.Errorf("%s shows its package clause or registration", e.Name)
		}
	}
	if err := writeManifestTo(discard{}); err != nil {
		t.Fatal(err)
	}
}

func TestAPIReadsDocumentedPackages(t *testing.T) {
	api, err := readAPI(moduleRoot())
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"components.Switch", "components.SwitchOptions", "chart.Bar", "plain.NewShell", "theme.Onevra"} {
		if _, ok := api[key]; !ok {
			t.Errorf("API has no %s", key)
		}
	}
	if f := api["components.SwitchOptions"].Fields; len(f) == 0 || f[0].Doc == "" {
		t.Errorf("SwitchOptions fields = %+v, want documented fields", f)
	}
}

type discard struct{}

func (discard) Write(p []byte) (int, error) { return len(p), nil }
