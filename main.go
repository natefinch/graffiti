package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/natefinch/graffiti/tags"
)

func main() {
	if err := makeCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func makeCmd() *cobra.Command {
	base := &cobra.Command{
		Use:   "graffiti",
		Short: "generate struct tags",
		Long:  "Graffiti generates struct tags for your go code.",
	}

	// Order here determines order in help output.
	base.AddCommand(genCmd())
	base.AddCommand(runCmd())
	addtopics(base)
	return base
}

func genCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen <tags> [target]",
		Short: "Generate struct tags for go structs in a file or directory.",
		Long:  genUsage,
	}
	var types, mapping string
	var isTempl, dryRun bool
	cmd.Flags().StringVarP(&types, "types", "t", "", "Generate tags only for these types (comma separated list).")
	cmd.Flags().StringVarP(&mapping, "map", "m", "", "Map field names to alternate tag names (see help mappings).")
	cmd.Flags().BoolVarP(&isTempl, "gotemplate", "g", false, "If set, tags is a go template (see help templates).")
	cmd.Flags().BoolVarP(&dryRun, "dryrun", "d", false, "If set, changes are written to stdout instead of to the files.")

	addtopics(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		// Tags is required, target is optional
		if len(args) != 2 && len(args) != 1 {
			fmt.Printf("Wrong number of arguments, expected 1 or 2, got %d", len(args))
			os.Exit(-1)
		}

		opt := tags.Options{DryRun: dryRun}

		if len(args) > 1 {
			opt.Target = args[1]
		} else {
			opt.Target = "."
		}

		if mapping != "" {
			m, err := makeMap(mapping)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			opt.Mapping = m
		}

		if types != "" {
			opt.Types = strings.Split(types, ",")
		}

		if !isTempl {
			opt.Tags = strings.Split(args[0], ",")
		} else {
			t, err := template.New("tag template").Parse(args[1])
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			opt.Template = t
		}
		if err := tags.Generate(opt); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}

	return cmd
}

func makeMap(val string) (map[string]string, error) {
	maps := strings.Split(val, ";")
	mapping := map[string]string{}
	for _, m := range maps {
		parts := strings.SplitN(m, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("badly formatted mapping: " + m)
		}
		mapping[parts[0]] = parts[1]
	}
	return mapping, nil
}

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <file>",
		Short: "Run graffiti commands embedded in a go file.",
		Long:  runUsage,
	}
	cmd.Run = func(cmd *cobra.Command, args []string) {
		// File is required.
		if len(args) != 1 {
			fmt.Printf("Wrong number of arguments, expected 1, got %d", len(args))
			os.Exit(-1)
		}
		if err := run(args[0]); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	return cmd
}

func addtopics(base *cobra.Command) {
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
}
