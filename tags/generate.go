package tags

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// Options contains the data needed to generate tags for a target (file or
// package).
type Options struct {
	// Target is the file or directory to parse, or current dir if empty.
	Target string
	// Tags contains schema names like json or yaml.
	Tags []string
	// Template is used to generate the contents of the struct tag field if Tags is empty.
	Template *template.Template
	// Mapping contains field name conversions. If a field isn't in the map, lowercase of field name is used.
	Mapping map[string]string
	// Types to generate tags for, if empty, all structs will have tags generated.
	Types []string
	// DryRun indicates whether we should simply write to StdOut rather than writing to the files.
	DryRun bool
}

// Generate generates tags according to the given options.
func Generate(o Options) error {
	i, err := os.Stat(o.Target)
	if err != nil {
		return err
	}
	if !i.IsDir() {
		return genfile(o, o.Target)
	}

	p, err := build.Default.ImportDir(o.Target, 0)
	if err != nil {
		return err
	}
	for _, f := range p.GoFiles {
		if err := genfile(o, f); err != nil {
			return err
		}
	}
	return nil
}

// genfile generates struct tags for the given file.
func genfile(o Options, file string) error {
	fset := token.NewFileSet()
	n, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	b, err := gen(o, fset, n)
	if err != nil {
		return err
	}
	if b == nil {
		// no changes
		return nil
	}
	if o.DryRun {
		_, err := fmt.Fprintf(os.Stdout, "%s\n", b)
		return err
	}
	return ioutil.WriteFile(file, b, 0644)
}

func gen(o Options, fset *token.FileSet, n ast.Node) ([]byte, error) {
	v := &visitor{Options: o}
	ast.Walk(v, n)
	if v.err != nil {
		return nil, v.err
	}
	if !v.changed {
		return nil, nil
	}
	c := printer.Config{Mode: printer.RawFormat}
	buf := &bytes.Buffer{}
	if err := c.Fprint(buf, fset, n); err != nil {
		return nil, fmt.Errorf("error printing output: %s", err)
	}
	b, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}
	return b, nil
}

// visitor is a wrapper around Options that implement the ast.Visitor interface
// and some helper methods.  Since ast.Walk doesn't let you return values,
// we instead set the return values in this struct.
type visitor struct {
	Options
	// changed is true if the AST was changed by our code.
	changed bool
	// err is non-nil if there was an error processing the file.
	err error
}

// Visit implements ast.Visitor and does the meat of the tag generation.
func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil || v.err != nil {
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
				val, err := v.gen(name)
				if err != nil {
					v.err = err
					return nil
				}
				f.Tag.Value = val
				v.changed = true
			}
		}
	}
	return v
}

// shouldGen reports whether graffiti should generate tags for the struct with
// the given name.
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

// gen creates the struct tag for the given field name, according to the options
// set.
func (v visitor) gen(name string) (string, error) {
	if m, ok := v.Mapping[name]; ok {
		name = m
	} else {
		name = TitleToCamel(name)
	}
	if len(v.Tags) > 0 {
		vals := make([]string, len(v.Tags))
		for i, t := range v.Tags {
			vals[i] = fmt.Sprintf("%s:%q", t, name)
		}
		return "`" + strings.Join(vals, " ") + "`", nil
	}

	// no tags means we have a template
	buf := &bytes.Buffer{}
	err := v.Template.Execute(buf, struct{ F string }{name})
	if err != nil {
		return "", err
	}
	return "`" + buf.String() + "`", nil
}
