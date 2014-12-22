package main

import (
	"reflect"
	"testing"
	"text/template"

	"github.com/natefinch/graffiti/tags"
)

func TestMakeOptions(t *testing.T) {
	isTempl := false
	dryRun := true
	types := "foo,bar"
	mapping := "ID=_id;Bar=notbar"
	args := []string{"json,yaml", "foo.go"}
	o, err := makeOptions(types, mapping, isTempl, dryRun, args)
	if err != nil {
		t.Fatalf("Unexpected error from makeOptions: %s", err)
	}
	expected := tags.Options{
		Target:  "foo.go",
		Tags:    []string{"json", "yaml"},
		Mapping: map[string]string{"ID": "_id", "Bar": "notbar"},
		Types:   []string{"foo", "bar"},
		DryRun:  dryRun,
	}

	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected:\n%#v\ngot:\n%#v", expected, o)
	}
}

func TestMakeOptionsTemplate(t *testing.T) {
	isTempl := true
	dryRun := false
	types := "foo,bar"
	mapping := "ID=_id;Bar=notbar"
	args := []string{`json:"{{.F}}"`}
	o, err := makeOptions(types, mapping, isTempl, dryRun, args)
	if err != nil {
		t.Fatalf("Unexpected error from makeOptions: %s", err)
	}
	expected := tags.Options{
		Target:  ".",
		Mapping: map[string]string{"ID": "_id", "Bar": "notbar"},
		Types:   []string{"foo", "bar"},
		DryRun:  dryRun,
	}

	expectedT := template.Must(template.New("tag template").Parse(`json:"{{.F}}"`))
	gotT := o.Template
	o.Template = nil

	if !reflect.DeepEqual(o, expected) {
		t.Errorf("Expected:\n%#v\ngot:\n%#v", expected, o)
	}

	if !reflect.DeepEqual(*gotT, *expectedT) {
		t.Errorf("Expected:\n%#v\ngot:\n%#v", expected, o)
	}
}
