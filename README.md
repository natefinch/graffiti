graffiti
========

Graffiti is a tool to automatically add struct tags to fields in your go code.

This repo is still under heavy development and should not be used by anyone.

By default, when given a list of schemas, graffiti will populate struct tags for
all exported fields of structs in the given file, with the lowercase of the
field name.


For example, 

`graffiti json,yaml foo.go`

where foo.go looks like this:

```
package foo

type foo struct {
	ID   string
	Name string
	mu   sync.Mutex
}
```

Will produce the following output:

```
package foo

type foo struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"yaml"`
	mu   sync.Mutex
}
```

The idea is to support arbitrary tags (no need to bake-in the tag types), and
allow for some simple rules, like tagging Id as _id for bson.

Also, the idea is to support in-file tags with the appropriate command-line so
that run-on-save can work for your favorite editor, something like the following:

//graffiti: json,yaml

Then whenever you save, graffiti will update the struct tags in this file.



