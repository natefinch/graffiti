package main

import (
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	args := []string{"-yaml, -json", "foo.go"}
	info, err := parseArgs(args)
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
	expected := cliInfo{
		filename: "foo.go",
		tags: []tag{
			tag{
				Name: "yaml",
			},
			tag{
				Name: "json",
			},
		},
	}

	if !reflect.DeepEqual(expected, info) {
		t.Errorf("expected %#v, got: %#v", expected, info)
	}
}

func TestParseMoreArgs(t *testing.T) {
	args := []string{"-bson", "ID=_id;Bar=notbar", "-yaml", "foo.go"}
	info, err := parseArgs(args)
	if err != nil {
		t.Errorf("unexpected error: %#v", err)
	}
	expected := cliInfo{
		filename: "foo.go",
		tags: []tag{
			tag{
				Name: "bson",
				Map: map[string]string{
					"ID":  "_id",
					"Bar": "notbar",
				},
			},
			tag{
				Name: "yaml",
			},
		},
	}

	if !reflect.DeepEqual(expected, info) {
		t.Errorf("expected %#v, got: %#v", expected, info)
	}
}
