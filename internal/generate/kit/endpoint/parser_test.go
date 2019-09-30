package endpoint

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	spec, err := Parse("./testdata/parser", "Service")
	if err != nil {
		t.Fatal(err)
	}

	expected := ServiceSpec{
		Name: "Service",
		Package: PackageSpec{
			Name: "parser",
			Path: "sagikazarmark.dev/mga/internal/generate/kit/endpoint/testdata/parser",
		},
		Endpoints: []EndpointSpec{
			{
				Name: "Call",
			},
		},
	}

	if !reflect.DeepEqual(expected, spec) {
		t.Error("the parsed spec does not match the expected one")
	}
}
