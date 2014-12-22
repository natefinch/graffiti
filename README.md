graffiti
========

Graffiti is a tool to automatically add struct tags to fields in your go code.

### Gen Command
Generates struct tags for a specific target (file or directory).

```
Usage: 
  graffiti gen <tags> [target] [flags]

Available Flags:
  -d, --dryrun		If set, changes are written to stdout instead of to the files.
  -f, --format		If set, tags is a go template (see help templates).
  -m, --map=""		Map field names to alternate tag names (see help mappings).
  -t, --types=""	Generate tags only for these types (comma separated list).
```

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

### Run Command

Reads graffiti commands from a go file and executes them.

```
Usage: 
  graffiti run <file>
```

The run command parses go files, looking for comments with the prefix `graffiti:
`. All text following this prefix is used as the command line for graffiti.

The commands are always passed to the graffiti gen command, and the current file
is assumed to be the target if no target is given.  If there are multiple
graffiti commands in the file, all will be run in sequence.  Flags can be used
as usual, and the embedded CLI supports single and double quotes for arguments
similar to how /bin/sh works.

For example this command embedded in the foo.go file:

	// graffiti: -t MyStruct json,yaml

is equivalent to this command line invocation:

	graffiti gen -t MyStruct json,yaml foo.go