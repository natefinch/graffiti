package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, `usage: "graffiti <file>"`)
		os.Exit(-1)
	}
	filename := os.Args[1]

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %s\n", err)
		os.Exit(-1)
	}

	ast.Walk(visitor(structTypes), f)
	c := printer.Config{Mode: printer.TabIndent}
	if err := c.Fprint(os.Stdout, fset, f); err != nil {
		fmt.Fprintf(os.Stderr, "error printing: %s", err)
		os.Exit(-1)
	}
}

func structTypes(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	if s, ok := n.(*ast.StructType); ok {
		for _, f := range s.Fields.List {
			if len(f.Names) > 1 {
				// skip fields declared as a, b, c int
				continue
			}
			if f.Tag == nil {
				f.Tag = &ast.BasicLit{}
			}
			if f.Tag.Value == "" {
				f.Tag.Value = fmt.Sprintf("`json:\"%s\"`", strings.ToLower(f.Names[0].Name))
			}
		}
	}
	return visitor(structTypes)
}

type visitor func(node ast.Node) (w ast.Visitor)

func (v visitor) Visit(node ast.Node) (w ast.Visitor) {
	return v(node)
}
