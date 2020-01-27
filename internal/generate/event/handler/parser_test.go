package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"sagikazarmark.dev/mga/internal/generate/gentypes"
)

func TestParse(t *testing.T) {
	expected := Event{
		Name: "Event",
		Package: gentypes.PackageRef{
			Name: "parser",
			Path: "sagikazarmark.dev/mga/internal/generate/event/handler/testdata/parser",
		},
	}

	spec, err := Parse("./testdata/parser", "Event")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, spec, "the parsed event does not match the expected one")
}
