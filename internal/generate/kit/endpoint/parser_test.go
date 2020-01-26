package endpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	def, err := Parse("./testdata/parser", "Service")
	if err != nil {
		t.Fatal(err)
	}

	expected := PackageDefinition{
		HeaderText:  "",
		PackageName: "parserdriver",
		LogicalName: "parser",
		EndpointSets: []SetDefinition{
			{
				BaseName: "",
				Service: ServiceDefinition{
					Name:        "Service",
					PackageName: "parser",
					PackagePath: "sagikazarmark.dev/mga/internal/generate/kit/endpoint/testdata/parser",
				},
				Endpoints: []EndpointDefinition{
					{
						Name: "Call",
					},
				},
				WithOpenCensus: false,
			},
		},
	}

	assert.Equal(t, expected, def, "the parsed definition does not match the expected one")
}
