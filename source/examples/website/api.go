//go:build !js

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

const modulePath = "github.com/oneo-me/onevra/source"

// apiPackages are the packages, from the module's root, whose API the
// documentation can show with <GoAPI name="components.Switch" />.
var apiPackages = []string{
	"components",
	"components/chart",
	"components/chat",
	"components/command",
	"components/flow",
	"components/gesture",
	"components/graph",
	"components/shaders",
	"components/treemap",
	"shells/plain",
	"shells/sidebar",
	"theme",
}

// apiEntry is one exported function or type, keyed by its package's name and
// its own, as components.Switch. Doc is the doc comment as Markdown.
type apiEntry struct {
	Kind    string     `json:"kind"`
	Package string     `json:"package"`
	Decl    string     `json:"decl"`
	Doc     string     `json:"doc,omitempty"`
	Fields  []apiField `json:"fields,omitempty"`
	Values  []apiValue `json:"values,omitempty"`
	Methods []apiFunc  `json:"methods,omitempty"`
}

type apiField struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Doc  string `json:"doc,omitempty"`
}

// apiValue is a group of a type's constants or variables, as declared.
type apiValue struct {
	Decl string `json:"decl"`
	Doc  string `json:"doc,omitempty"`
}

type apiFunc struct {
	Name string `json:"name"`
	Decl string `json:"decl"`
	Doc  string `json:"doc,omitempty"`
}

// readAPI reads the exported API of apiPackages from the module at root.
func readAPI(root string) (map[string]apiEntry, error) {
	api := map[string]apiEntry{}
	for _, rel := range apiPackages {
		if err := readPackage(api, root, rel); err != nil {
			return nil, err
		}
	}
	return api, nil
}

func readPackage(api map[string]apiEntry, root, rel string) error {
	fset := token.NewFileSet()
	dir := filepath.Join(root, filepath.FromSlash(rel))
	pkgs, err := parser.ParseDir(fset, dir, func(info fs.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, parser.ParseComments)
	if err != nil {
		return err
	}
	var files []*ast.File
	for _, p := range pkgs {
		for _, f := range p.Files {
			files = append(files, f)
		}
	}
	importPath := modulePath + "/" + rel
	p, err := doc.NewFromFiles(fset, files, importPath)
	if err != nil {
		return err
	}
	r := docReader{fset: fset, files: files, pkg: p}
	add := func(name string, e apiEntry) error {
		key := p.Name + "." + name
		if _, dup := api[key]; dup {
			return fmt.Errorf("%s is declared by two packages named %s", key, p.Name)
		}
		e.Package = importPath
		api[key] = e
		return nil
	}
	for _, f := range p.Funcs {
		if err := add(f.Name, apiEntry{Kind: "func", Decl: r.signature(f.Decl), Doc: r.markdown(f.Doc)}); err != nil {
			return err
		}
	}
	for _, t := range p.Types {
		e := apiEntry{Kind: "type", Decl: r.node(t.Decl, false), Doc: r.markdown(t.Doc), Fields: r.fields(t)}
		for _, group := range append(t.Consts, t.Vars...) {
			e.Values = append(e.Values, apiValue{Decl: r.node(group.Decl, true), Doc: r.markdown(group.Doc)})
		}
		for _, m := range t.Methods {
			e.Methods = append(e.Methods, apiFunc{Name: m.Name, Decl: r.signature(m.Decl), Doc: r.markdown(m.Doc)})
		}
		if err := add(t.Name, e); err != nil {
			return err
		}
		// Constructors are listed under their type by go/doc, and on their
		// own here, as the documentation names them.
		for _, f := range t.Funcs {
			if err := add(f.Name, apiEntry{Kind: "func", Decl: r.signature(f.Decl), Doc: r.markdown(f.Doc)}); err != nil {
				return err
			}
		}
	}
	return nil
}

type docReader struct {
	fset  *token.FileSet
	files []*ast.File
	pkg   *doc.Package
}

// markdown converts a doc comment, linking names to pkg.go.dev.
func (r docReader) markdown(text string) string {
	p := r.pkg.Printer()
	p.DocLinkBaseURL = "https://pkg.go.dev"
	p.HeadingLevel = 4
	return strings.TrimSpace(string(p.Markdown(r.pkg.Parser().Parse(text))))
}

// plain converts a doc comment to unwrapped text, for a table cell.
func (r docReader) plain(text string) string {
	p := r.pkg.Printer()
	p.TextWidth = -1
	return strings.TrimSpace(string(p.Text(r.pkg.Parser().Parse(text))))
}

// signature prints a function's declaration without its body or comment.
func (r docReader) signature(f *ast.FuncDecl) string {
	return r.print(&ast.FuncDecl{Recv: f.Recv, Name: f.Name, Type: f.Type})
}

// node prints a declaration, with the comments inside it when asked.
func (r docReader) node(decl *ast.GenDecl, comments bool) string {
	if !comments {
		d := *decl
		d.Doc = nil
		return r.print(&d)
	}
	for _, f := range r.files {
		if f.FileStart <= decl.Pos() && decl.End() <= f.FileEnd {
			d := *decl
			d.Doc = nil
			return r.print(&printer.CommentedNode{Node: &d, Comments: f.Comments})
		}
	}
	return r.print(decl)
}

func (r docReader) print(node any) string {
	var b bytes.Buffer
	cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}
	if err := cfg.Fprint(&b, r.fset, node); err != nil {
		return ""
	}
	return b.String()
}

// fields lists a struct type's exported fields; go/doc has already removed
// the unexported ones.
func (r docReader) fields(t *doc.Type) []apiField {
	if len(t.Decl.Specs) != 1 {
		return nil
	}
	spec, ok := t.Decl.Specs[0].(*ast.TypeSpec)
	if !ok {
		return nil
	}
	st, ok := spec.Type.(*ast.StructType)
	if !ok {
		return nil
	}
	var out []apiField
	for _, f := range st.Fields.List {
		typ := r.print(f.Type)
		text := f.Doc.Text()
		if text == "" {
			text = f.Comment.Text()
		}
		names := f.Names
		if len(names) == 0 {
			// An embedded field is named by its type.
			name := strings.TrimPrefix(typ, "*")
			if i := strings.LastIndex(name, "."); i >= 0 {
				name = name[i+1:]
			}
			names = []*ast.Ident{{Name: name}}
		}
		for _, n := range names {
			if !n.IsExported() {
				continue
			}
			out = append(out, apiField{Name: n.Name, Type: typ, Doc: r.plain(text)})
		}
	}
	return out
}
