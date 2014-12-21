package main

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/kballard/go-shellquote"
	"github.com/spf13/cobra"
)

const prefix = "graffiti: "

func run(target string) error {
	i, err := os.Stat(target)
	if err != nil {
		return err
	}
	if !i.IsDir() {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, target, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		return runfile(f, target)
	}

	p, err := build.Default.ImportDir(target, 0)
	if err != nil {
		return err
	}
	for _, file := range p.GoFiles {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		if err := runfile(f, file); err != nil {
			return err
		}
	}
	return nil
}

// runfile generates struct tags for the given file.
func runfile(f *ast.File, file string) error {
	cmd := makeCmd()
	for _, c := range cmd.Commands() {
		if c.Name() == "gen" {
			runner := c.Run
			// Default to appending the filename to the args, if a target isn't
			// specified.
			c.Run = func(cmd *cobra.Command, args []string) {
				if len(args) == 1 {
					args = append(args, file)
				}
				runner(cmd, args)
			}
		}
	}

	for _, c := range f.Comments {
		if strings.HasPrefix(c.Text(), prefix) {
			args, err := shellquote.Split(c.Text()[len(prefix):])
			if err != nil {
				return err
			}
			args = append([]string{"graffiti", "gen"}, args...)
			os.Args = args
			if err := cmd.Execute(); err != nil {
				return err
			}
		}
	}
	return nil
}
