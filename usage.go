package main

const (
	genUsage = `
By default, target is a file, and tags is a comma-separated list of tag names.
The value for a field's tag is the lowercase of the field name (only exported
fields have tags generated). 

For example, 'graffiti gen foo.go yaml,json' will transform the following
struct:

type foo struct {
	ID string
	Name string
}

into this:

type foo struct {
	ID string ` + "`" + `yaml:"id" json:"id"` + "`" + `
	Name string ` + "`" + `yaml:"name" json:"name"` + "`" + `
}

`

	runUsage = `


`

	packages = ` By default, graffiti gen expects to be given a single filename
that contains the types to generate tags for.  If you pass the -p (--package)
flag, it instead expects to be given the name of a package. This should be the
same format as the go tool uses.  See 'go help packages' for full details.

Example:

	graffiti gen -p github.com/foo/bar/... json

This would generate json tags for all types in all packages on your local
machine with import paths that start with github.com/foo/bar/.
`
	mappings = `By default, graffiti creates a struct tag for each exported
field in the struct that is simple the lowercase of the field name.  For example
running this command:

 	graffiti gen foo.go json 

would produce output like this:

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

	graffiti gen foo.go bson -m ID=_id

which would produce output like this:

 	type foo struct {
 		ID string  ` + "`" + `bson:"_id"` + "`" + `
 		Nae string ` + "`" + `bson:"name"` + "`" + `
 	}

Note that all fields which don't match a key in the mapping default to the
normal behavior of simply being lowercased.

To make multiple translations, just separate them with a semicolon thusly:

	graffiti gen foo.go bson -m ID=_id;Subject=title

Note that the entire right hand side of a mapping is used as the output of the
tag value, so you can specify more than just the name, for example adding
",omitempty" to json or yaml fields.

	graffiti gen foo.go json -m Value=value,omitempty
`

	gotemplate = `By default, graffiti takes as its last argument a comma-separated list of tag
names which will be used to generate the tags. The value for the tags is the
lowercase of the field name (only exported fields have tags generated). For
example, 'graffiti file foo.go yaml,json' will transform the following struct:

	type foo struct {
		ID string
		Name string
	}

into this:

	type foo struct {
		ID string ` + "`" + `yaml:"id" json:"id"` + "`" + `
		Name string ` + "`" + `yaml:"name" json:"name"` + "`" + `
	}

If you want complete control of the output, you can use the -g (--template)
flag, in which case you pass the command a Go template that is used to generate
the tags for the field.  In the template, {{.F}} will be populated with the
lowercase name of the field, or the mapped value if using -map (see help
mappings).  The output of the template will be surrounded by backticks in the
output.  Tip: surround the template value in single quotes so that you can use
double quotes inside the value.

For example:

	graffiti file foo.go -template 'json:"{{.F}}"'

would produce the following from the same above struct:

type foo struct {
	ID string ` + "`" + `json:"id"` + "`" + `
	Name string ` + "`" + `json:"name"` + "`" + `
}
`
)
