package endpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sagikazarmark.dev/mga/pkg/gentypes"
)

func TestParse(t *testing.T) {
	expected := Service{
		TypeRef: gentypes.TypeRef{
			Name: "Service",
			Package: gentypes.PackageRef{
				Name: "parser",
				Path: "sagikazarmark.dev/mga/internal/generate/kit/endpoint/testdata/parser",
			},
		},
		Methods: []ServiceMethod{
			{
				Name: "Call",
				RequestParameters: []gentypes.Argument{
					{
						Name: "req",
						Type: gentypes.TypeRef{
							Name: "interface{}",
						},
					},
				},
				ResponseParameters: []gentypes.Argument{
					{
						Name: "",
						Type: gentypes.TypeRef{
							Name: "interface{}",
						},
					},
				},
			},
		},
	}

	actual, err := Parse("./testdata/parser", "Service")
	require.NoError(t, err)

	assert.Equal(t, expected, actual, "the parsed service does not match the expected one")
}
