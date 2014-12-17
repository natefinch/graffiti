package tags

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
	"text/template"
)

// Options contains the data needed to generate tags for a target (file or
// package).
type Options struct {
	Target    string
	IsPackage bool
	Tags      []string
	Template  *template.Template
	Mapping   map[string]string
	Types     []string
}

// Generate generates tags according to the given options.
func Generate(o Options) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, o.Target, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error reading file: %s", err)
	}
	if err := walk(visitor(o), f); err != nil {
		return err
	}

	c := printer.Config{Mode: printer.RawFormat}
	if err := c.Fprint(os.Stdout, fset, f); err != nil {
		return fmt.Errorf("error printing output: %s", err)
	}

	return nil
}

// walk is a function that wraps ast.Walk in a defer recover. This is because
// ast.Walk doesn't give us a way to bail out with an error, so we have to panic
// instead.
func walk(v ast.Visitor, f *ast.File) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("unexpected panic while parsing file: %v", r)
			}
		}
	}()
	ast.Walk(v, f)
	return nil
}

type visitor Options

func (v visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	if t, ok := n.(*ast.TypeSpec); ok {
		if s, ok := t.Type.(*ast.StructType); ok {
			if !v.shouldGen(t.Name.Name) {
				return v
			}
			for _, f := range s.Fields.List {
				if len(f.Names) > 1 {
					// skip fields declared as a, b, c int
					continue
				}
				name := f.Names[0].Name
				if !ast.IsExported(name) {
					// skip non-exported names
					continue
				}
				if f.Tag == nil {
					f.Tag = &ast.BasicLit{}
				}
				f.Tag.Value = v.gen(name)
			}
		}
	}
	return v
}

func (v visitor) shouldGen(name string) bool {
	if len(v.Types) == 0 {
		return true
	}
	for _, typ := range v.Types {
		if typ == name {
			return true
		}
	}
	return false
}

func (v visitor) gen(name string) string {
	if m, ok := v.Mapping[name]; ok {
		name = m
	} else {
		name = strings.ToLower(name)
	}
	if len(v.Tags) > 0 {
		vals := make([]string, len(v.Tags))
		for i, t := range v.Tags {
			vals[i] = fmt.Sprintf("%s:%q", t, name)
		}
		return strings.Join(vals, " ")
	}
	// no tages means we have a template
	buf := &bytes.Buffer{}
	err := v.Template.Execute(buf, struct{ F string }{name})
	if err != nil {
		panic(err)
	}
	return buf.String()
}
