package endpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		serviceName        string
		expectedDefinition PackageDefinition
	}{
		{
			serviceName: "Service",
			expectedDefinition: PackageDefinition{
				HeaderText:  "",
				PackageName: "parserdriver",
				EndpointSets: []EndpointSetDefinition{
					{
						BaseName: "",
						Service: ServiceDefinition{
							Name:        "Service",
							PackageName: "parser",
							PackagePath: "sagikazarmark.dev/mga/internal/generate/kit/endpoint/testdata/parser",
						},
						Endpoints: []EndpointDefinition{
							{
								Name:          "Call",
								OperationName: "parser.Call",
							},
						},
						WithOpenCensus: false,
					},
				},
			},
		},
		{
			serviceName: "OtherService",
			expectedDefinition: PackageDefinition{
				HeaderText:  "",
				PackageName: "parserdriver",
				EndpointSets: []EndpointSetDefinition{
					{
						BaseName: "Other",
						Service: ServiceDefinition{
							Name:        "OtherService",
							PackageName: "parser",
							PackagePath: "sagikazarmark.dev/mga/internal/generate/kit/endpoint/testdata/parser",
						},
						Endpoints: []EndpointDefinition{
							{
								Name:          "Call",
								OperationName: "parser.Other.Call",
							},
						},
						WithOpenCensus: false,
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.serviceName, func(t *testing.T) {
			def, err := Parse("./testdata/parser", test.serviceName)
			require.NoError(t, err)

			assert.Equal(t, test.expectedDefinition, def, "the parsed definition does not match the expected one")
		})
	}
}
