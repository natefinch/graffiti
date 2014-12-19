graffiti
========

Graffiti is a tool to automatically add struct tags to fields in your go code.

This repo is still under heavy development and should not be used by anyone.

```
Usage: 
  graffiti gen <tags> [target] [flags]

Available Flags:
  -d, --dryrun=false: If set, changes are written to stdout instead of to the files.
  -g, --gotemplate=false: If set, tags is a go template (see help templates).
      --help=false: help for gen
  -m, --map="": Map field names to alternate tag names (see help mappings).
  -t, --types="": Generate tags only for these types (comma separated list).
```

Generates tags for a specific target (file or directory).

If no target is given, all go files in the current directory are processed. By
default tags is a comma-separated list of schema names like json or yaml. The
value for a field's tag is the lowercase of the field name. Only exported fields
have tags generated. 

For example, 

`graffiti gen json,yaml foo.go`

where foo.go looks like this:

```go
package foo

type foo struct {
	ID   string
	Name string
	mu   sync.Mutex
}
```

Will produce the following output:

```go
package foo

type foo struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"yaml"`
	mu   sync.Mutex
}
```

### TODO

Working on support for a `graffiti run` command that will parse a go file and
run graffiti commands embedded in comments in the file (much like go generate).
Then, run-on-save can work for your favorite editor, something like the
following:

//graffiti: json,yaml

Then whenever you save, graffiti will update the struct tags in this file.