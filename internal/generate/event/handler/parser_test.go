package handler

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	spec, err := Parse("./testdata/parser", "Event")
	if err != nil {
		t.Fatal(err)
	}

	expected := EventSpec{
		Name: "Event",
		Package: PackageSpec{
			Name: "parser",
			Path: "sagikazarmark.dev/mga/internal/generate/event/handler/testdata/parser",
		},
	}

	if !reflect.DeepEqual(expected, spec) {
		t.Error("the parsed spec does not match the expected one")
	}
}
