package main

import (
	//"errors"
	//"flag"
	// "fmt"
	// "go/ast"
	// "go/parser"
	// "go/printer"
	// "go/token"
	"os"
	//"strings"
	//"unicode"

	"github.com/spf13/cobra"
)

func main() {
	base := &cobra.Command{
		Use:   "graffiti",
		Short: "generate struct tags",
		Long:  "Graffiti generates struct tags for your go code.",
	}

	// Order here determines order in help output.
	genCmd(base)
	runCmd(base)
	topics(base)

	if err := base.Execute(); err != nil {
		os.Exit(1)
	}
}

func genCmd(base *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "gen <target> <tags>",
		Short: "Generate struct tags for structs in a file.",
		Long:  genUsage,
	}
	var types, mapping string
	var isTempl, isPkg bool
	cmd.Flags().StringVarP(&types, "types", "t", "", "Generate tags only for these types (comma separated list).")
	cmd.Flags().StringVarP(&mapping, "map", "m", "", "Map field names to alternate tag names (see help mappings).")
	cmd.Flags().BoolVarP(&isTempl, "gotemplate", "g", false, "If set, tags is a go template (see help templates).")
	cmd.Flags().BoolVarP(&isPkg, "pkg", "p", false, "If set, target is a package (see help packages).")

	base.AddCommand(cmd)
	cmd.Run = func(cmd *cobra.Command, args []string) {
		// Target & tags are required.
		if len(args) != 2 {
			cmd.Usage()
			return
		}
	}

	topics(cmd)
}

func runCmd(base *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "run <file>",
		Short: "Run graffiti commands embedded in a go file.",
		Long:  runUsage,
	}
	base.AddCommand(cmd)
	cmd.Run = func(cmd *cobra.Command, args []string) {
		// File is required.
		if len(args) != 1 {
			cmd.Usage()
			return
		}
	}
}

func topics(base *cobra.Command) {
	base.AddCommand(&cobra.Command{
		Use:   "mappings",
		Short: "description of field name mappings",
		Long:  mappings,
	})
	base.AddCommand(&cobra.Command{
		Use:   "templates",
		Short: "how to use templated output",
		Long:  gotemplate,
	})
	base.AddCommand(&cobra.Command{
		Use:   "packages",
		Short: "how to generate tags for package(s)",
		Long:  packages,
	})
}

/*

func old() {
	// get command line from file
	// -from foo.go

	// mapping
	// -map "ID=_id"

	// use a template
	// -template 'json:"{{.N}},omitempty" bson:"{{.N}}"'
	// not compatible w/ basic version

	// only make tags for these types
	// -types MyStruct,MyStruct2

	// run over all files in the package
	// -pkg <package>

	// basic (default to lowercase, run over all types in foo.go)
	// graffiti -tags yaml,bson,json foo.go

	info, err := parseArgs(os.Args[1:])

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing command line: %s\n", err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, info.filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file: %s\n", err)
		os.Exit(-1)
	}

	ast.Walk(visitor(tagger), f)
	c := printer.Config{Mode: printer.RawFormat}
	if err := c.Fprint(os.Stdout, fset, f); err != nil {
		fmt.Fprintf(os.Stderr, "error printing: %s", err)
		os.Exit(-1)
	}
}


// cliInfo represents the CLI flags
type cliInfo struct {
	filename string
	tags     []tag
}

// tag represents a single schema with an optional list of field name mappings.
type tag struct {
	Name string
	Map  map[string]string
}


// parse args parses the command line args (not including the command name).
func parseArgs(args []string) (cliInfo, error) {
	set := flag.NewFlagSet("", flag.ExitOnError)
	types := set.String("types", "", "comma separated list of types to generate tags for")
	mapping := set.String("map", "", "semicolon separated map of fieldName=tagName")
	templ := set.Bool("template", false, "use a go template to generate tags")

	set.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
		set.PrintDefaults()
	}

	if err := set.Parse(args); err != nil {
		return cliInfo{}, err
	}

	// now, the only things not parsed should either be just the filename
	info := cliInfo{}
	switch len(set.Args()) {
	case 1:
		// just the filename
		info.filename = set.Arg(0)
	default:
		return cliInfo{}, fmt.Errorf("wrong number of arguments after parsing flags: %d", len(set.Args()))
	}

	// now go through every flagged value
	errs := []string{}
	set.Visit(func(f *flag.Flag) {
		t := tag{Name: f.Name, Map: map[string]string{}}
		if f.Value.String() == "true" {
			// bool flag
			info.tags = append(info.tags, t)
			return
		}
		mappings := strings.Split(f.Value.String(), ";")
		for _, m := range mappings {
			vals := strings.SplitN(m, "=", 2)
			if len(vals) != 2 {
				errs = append(errs, fmt.Sprintf("invalid mapping for %s: %s", f.Name, m))
				continue
			}
			t.Map[vals[0]] = vals[1]
		}
		info.tags = append(info.tags, t)
	})
	if len(errs) > 0 {
		return cliInfo{}, errors.New(strings.Join(errs, "\n"))
	}
	return info, nil
}

func tagger(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	if s, ok := n.(*ast.StructType); ok {
		for _, f := range s.Fields.List {
			if len(f.Names) > 1 {
				// skip fields declared as a, b, c int
				continue
			}
			name := f.Names[0].Name
			if !unicode.IsUpper([]rune(name)[0]) {
				// only tag exported fields
				continue
			}
			if f.Tag == nil {
				f.Tag = &ast.BasicLit{}
			}
			if f.Tag.Value == "" {
				f.Tag.Value = fmt.Sprintf("`json:\"%s\"`", strings.ToLower(name))
			}
		}
	}
	return visitor(tagger)
}

type visitor func(node ast.Node) (w ast.Visitor)

func (v visitor) Visit(node ast.Node) (w ast.Visitor) {
	return v(node)
}
*/
