//go:build !js

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/oneo-me/onevra/source/examples/website/examples"
)

// sourceDir is package examples' directory: beside this file when the build
// keeps its path, as go run does, or else from the module's root.
func sourceDir() string {
	if _, file, _, ok := runtime.Caller(0); ok && filepath.IsAbs(file) {
		return filepath.Join(filepath.Dir(file), "examples")
	}
	return filepath.Join("examples", "website", "examples")
}

// manifest is what the website builds its documentation from.
type manifest struct {
	Examples []manifestEntry     `json:"examples"`
	API      map[string]apiEntry `json:"api"`
}

// moduleRoot is the Go module's directory, source/ of the repository.
func moduleRoot() string { return filepath.Join(sourceDir(), "..", "..", "..") }

type manifestEntry struct {
	Name string `json:"name"`
	// File is the example's path from the repository's root.
	File string `json:"file"`
	// Code is the file's source as the documentation shows it, without its
	// package clause and its registration.
	Code string `json:"code"`
}

// writeManifest prints every example with its source and the API the
// documentation describes, failing if a file does not register exactly one
// example or the registry has one no file declares.
func writeManifest() error { return writeManifestTo(os.Stdout) }

func writeManifestTo(w io.Writer) error {
	entries, err := readExamples(sourceDir())
	if err != nil {
		return err
	}
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name
	}
	if !slices.Equal(names, examples.Names()) {
		return fmt.Errorf("example files %v do not match the registered examples %v", names, examples.Names())
	}
	api, err := readAPI(moduleRoot())
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(manifest{Examples: entries, API: api}, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(append(data, '\n'))
	return err
}

func readExamples(dir string) ([]manifestEntry, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return nil, err
	}
	var entries []manifestEntry
	for _, path := range files {
		base := filepath.Base(path)
		if base == "examples.go" || strings.HasSuffix(base, "_test.go") {
			continue
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		code, err := shownSource(path, src)
		if err != nil {
			return nil, err
		}
		entries = append(entries, manifestEntry{
			Name: strings.ReplaceAll(strings.TrimSuffix(base, ".go"), "_", "-"),
			File: "source/examples/website/examples/" + base,
			Code: code,
		})
	}
	slices.SortFunc(entries, func(a, b manifestEntry) int { return strings.Compare(a.Name, b.Name) })
	return entries, nil
}

// shownSource cuts the package clause and the example's registration out of
// a file's source.
func shownSource(path string, src []byte) (string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return "", err
	}
	offset := func(p token.Pos) int { return fset.Position(p).Offset }
	type cut struct{ from, to int }
	cuts := []cut{{0, offset(f.Name.End())}}
	for _, decl := range f.Decls {
		g, ok := decl.(*ast.GenDecl)
		if !ok || g.Tok != token.VAR || !registers(g) {
			continue
		}
		from := g.Pos()
		if g.Doc != nil {
			from = g.Doc.Pos()
		}
		cuts = append(cuts, cut{offset(from), offset(g.End())})
	}
	if len(cuts) != 2 {
		return "", fmt.Errorf("%s: an example file registers exactly one example, with var _ = register(...) or registerWindow(...)", path)
	}
	var out bytes.Buffer
	last := 0
	for _, c := range cuts {
		out.Write(src[last:c.from])
		last = c.to
	}
	out.Write(src[last:])
	code := strings.TrimSpace(out.String())
	for strings.Contains(code, "\n\n\n") {
		code = strings.ReplaceAll(code, "\n\n\n", "\n\n")
	}
	return code + "\n", nil
}

// registers reports whether a declaration is var _ = register(...) or
// registerWindow(...).
func registers(g *ast.GenDecl) bool {
	for _, spec := range g.Specs {
		v, ok := spec.(*ast.ValueSpec)
		if !ok || len(v.Values) != 1 {
			continue
		}
		if call, ok := v.Values[0].(*ast.CallExpr); ok {
			if id, ok := call.Fun.(*ast.Ident); ok && (id.Name == "register" || id.Name == "registerWindow") {
				return true
			}
		}
	}
	return false
}
