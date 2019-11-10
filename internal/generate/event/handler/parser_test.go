package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, expected, spec, "the parsed spec does not match the expected one")
}
