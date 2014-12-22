package main

const (
	genUsage = `
Generates struct tags for a specific target (file or directory).

If no target is given, all go files in the current directory are processed. By
default tags is a comma-separated list of schema names like json or yaml. The
value for a field's tag is the lowercase of the field name. Only exported fields
have tags generated. `

	runUsage = `
Reads graffiti commands from a go file and executes them.

For example:

	// graffiti: <command line>

The commands are always passed to the graffiti gen command, and the current file
is assumed to be the target if no target is given.  If there are multiple
graffiti commands in the file, all will be run in sequence.  Flags can be used
as usual, and the embedded CLI supports single and double quotes for arguments
similar to how /bin/sh works. 
`

	mappings = `By default, graffiti creates a struct tag for each exported
field in the struct that is simple the lowercase of the field name.  For example
running this command:

 	graffiti gen json foo.go 

If foo.go contains this struct:

 	type foo struct {
 		Name string
 	}

The command would produce output like this:

 	type foo struct {
 		Name string ` + "`" + `json:"name"` + "`" + `
 	}

However, there can be times when the correct tag name is not the same as the
lowercase of the field name.  For these cases, you may provide a map of field
name to tag name using the -m (--map) flag on graffiti gen.  The value for the
flag is a semicolon delimited list of mappings.  Each mapping take the form of
FieldName=tagvalue.

For example, mongodb's bson requires each document to have an _id field.  This
obviously cannot be directly translated from a field name, so you need a struct
tag to change the name of the serialized field to _id.  To tell graffiti to make
this kind of a tag, you'd use a command line like this:

	graffiti gen bson foo.go -m ID=_id

which would produce output like this:

 	type foo struct {
 		ID string  ` + "`" + `bson:"_id"` + "`" + `
 		Name string ` + "`" + `bson:"name"` + "`" + `
 	}

Note that all fields which don't match a key in the mapping default to the
normal behavior of simply being lowercased.

To make multiple translations, just separate them with a semicolon thusly:

	graffiti gen bson foo.go -m ID=_id;Subject=title

Note that the entire right hand side of a mapping is used as the output of the
tag value, so you can specify more than just the name, for example adding
",omitempty" to json or yaml fields.

	graffiti gen json foo.go -m Value=value,omitempty
`

	gotemplate = `By default, graffiti takes as its last argument a comma-separated list of tag
names which will be used to generate the tags. The value for the tags is the
lowercase of the field name (only exported fields have tags generated). For
example, 'graffiti gen yaml,json foo.go' will transform the following struct:

	type foo struct {
		ID string
		Name string
	}

into this:

	type foo struct {
		ID string ` + "`" + `yaml:"id" json:"id"` + "`" + `
		Name string ` + "`" + `yaml:"name" json:"name"` + "`" + `
	}

If you want complete control of the output, you can use the -g (--gotemplate)
flag, in which case you pass the command a Go template that is used to generate
the tags for the field.  In the template, {{.F}} will be populated with the
lowercase name of the field, or the mapped value if using -map (see help
mappings).  The output of the template will be surrounded by backticks in the
output.  Tip: surround the template value in single quotes so that you can use
double quotes inside the value.

For example:

	graffiti gen --g 'json:"{{.F}}"' foo.go

would produce the following from the same above struct:

type foo struct {
	ID string ` + "`" + `json:"id"` + "`" + `
	Name string ` + "`" + `json:"name"` + "`" + `
}
`
)
