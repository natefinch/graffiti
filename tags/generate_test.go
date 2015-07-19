package tags

import (
	"go/parser"
	"go/token"
	"testing"
	"text/template"
)

func TestGenBasic(t *testing.T) {
	o := Options{Tags: []string{"json", "yaml"}}
	fset := token.NewFileSet()

	code := `
package foo

type bar struct {
	ID       string
	Name     string
	LastName string
	private  string
}
`[1:]

	n, err := parser.ParseFile(fset, "foo.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Unexpected setup error from parsefile: %v", err)
	}
	b, err := gen(o, fset, n)
	if err != nil {
		t.Fatalf("Unexpected error from gen: %v", err)
	}

	expected := `package foo

type bar struct {
	ID       string ` + "`" + `json:"id" yaml:"id"` + "`" + `
	Name     string ` + "`" + `json:"name" yaml:"name"` + "`" + `
	LastName string ` + "`" + `json:"lastName" yaml:"lastName"` + "`" + `
	private  string
}
`

	if expected != string(b) {
		t.Errorf("expected: \n%q\ngot:\n%q", expected, string(b))
	}
}

func TestGenMap(t *testing.T) {
	o := Options{Tags: []string{"json", "yaml"}, Mapping: map[string]string{"ID": "_id", "Name": "title"}}
	fset := token.NewFileSet()

	code := `
package foo

type bar struct {
	ID      string
	Name    string
	private string
}
`[1:]

	n, err := parser.ParseFile(fset, "foo.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Unexpected setup error from parsefile: %v", err)
	}
	b, err := gen(o, fset, n)
	if err != nil {
		t.Fatalf("Unexpected error from gen: %v", err)
	}

	expected := `package foo

type bar struct {
	ID      string ` + "`" + `json:"_id" yaml:"_id"` + "`" + `
	Name    string ` + "`" + `json:"title" yaml:"title"` + "`" + `
	private string
}
`

	if expected != string(b) {
		t.Errorf("expected: \n%q\ngot:\n%q", expected, string(b))
	}
}

func TestGenType(t *testing.T) {
	o := Options{Tags: []string{"json"}, Types: []string{"bar"}}
	fset := token.NewFileSet()

	code := `package foo

type bar struct {
	ID string
}

type not struct {
	ID string
}
`

	n, err := parser.ParseFile(fset, "foo.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Unexpected setup error from parsefile: %v", err)
	}
	b, err := gen(o, fset, n)
	if err != nil {
		t.Fatalf("Unexpected error from gen: %v", err)
	}

	expected := `package foo

type bar struct {
	ID string ` + "`" + `json:"id"` + "`" + `
}

type not struct {
	ID string
}
`

	if expected != string(b) {
		t.Errorf("expected: \n%q\ngot:\n%q", expected, string(b))
	}
}

func TestGenFormat(t *testing.T) {
	temp, err := template.New("tags").Parse(`foo:"{{.F}}" bar:"{{.F}}"`)
	if err != nil {
		t.Fatalf("Unexpected setup error parsing template: %v", err)
	}

	o := Options{Template: temp}
	fset := token.NewFileSet()

	code := `package foo

type bar struct {
	ID string
}
`

	n, err := parser.ParseFile(fset, "foo.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Unexpected setup error from parsefile: %v", err)
	}
	b, err := gen(o, fset, n)
	if err != nil {
		t.Fatalf("Unexpected error from gen: %v", err)
	}

	expected := `package foo

type bar struct {
	ID string ` + "`" + `foo:"id" bar:"id"` + "`" + `
}
`

	if expected != string(b) {
		t.Errorf("expected: \n%q\ngot:\n%q", expected, string(b))
	}
}

func TestGenFormatMapping(t *testing.T) {
	temp, err := template.New("tags").Parse(`foo:"{{.F}}" bar:"{{.F}}"`)
	if err != nil {
		t.Fatalf("Unexpected setup error parsing template: %v", err)
	}

	o := Options{
		Template: temp,
		Mapping:  map[string]string{"ID": "_id"},
	}
	fset := token.NewFileSet()

	code := `package foo

type bar struct {
	ID string
}
`

	n, err := parser.ParseFile(fset, "foo.go", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Unexpected setup error from parsefile: %v", err)
	}
	b, err := gen(o, fset, n)
	if err != nil {
		t.Fatalf("Unexpected error from gen: %v", err)
	}

	expected := `package foo

type bar struct {
	ID string ` + "`" + `foo:"_id" bar:"_id"` + "`" + `
}
`

	if expected != string(b) {
		t.Errorf("expected: \n%q\ngot:\n%q", expected, string(b))
	}
}
