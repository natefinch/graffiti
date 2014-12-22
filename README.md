graffiti
========

Graffiti is a tool to automatically add struct tags to fields in your go code.

This repo is still under heavy development and should not be used by anyone.

### Gen

```
Usage: 
  graffiti gen <tags> [target] [flags]

Available Flags:
  -d, --dryrun		If set, changes are written to stdout instead of to the files.
  -f, --format		If set, tags is a go template (see help templates).
  -m, --map=""		Map field names to alternate tag names (see help mappings).
  -t, --types=""	Generate tags only for these types (comma separated list).
```

Generates struct tags for a specific target (file or directory).

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

### Run

Graffiti can parse go source files and look for embedded graffiti commands in comments.  

Running 

	graffiti run foo.go

Will parse foo.go and look for comments of the form 

	// graffiti: <command line>

These commands will be run in the order they are found in the file.  Commands in the file are always run with the `graffiti gen` command, and if no target is given, the current file is the default target. Thus, this command embedded in the foo.go file:

	// graffiti: -t MyStruct json,yaml

is equivalent to this command line invocation:

	graffiti gen -t MyStruct json,yaml foo.go

