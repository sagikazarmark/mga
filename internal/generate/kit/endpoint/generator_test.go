package endpoint

import (
	"fmt"
	"go/types"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
	"sigs.k8s.io/controller-tools/pkg/loader"

	"sagikazarmark.dev/mga/pkg/gentypes"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "service_with_struct",
		},
		{
			name: "simple_service",
		},
		{
			name: "todo",
		},
		{
			name: "unnamed_param",
		},
		{
			name: "pointer_message",
		},
		{
			name: "different_package",
		},
		{
			name: "generics",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			pkgs, err := loader.LoadRootsWithConfig(
				&packages.Config{
					Mode: packages.NeedDeps | packages.NeedTypes,
				},
				fmt.Sprintf("./testdata/generator/%s", test.name),
			)
			require.NoError(t, err)

			pkg := pkgs[0]

			pkg.NeedTypesInfo()

			service := pkg.Types.Scope().Lookup("Service").Type().(*types.Named)

			file := File{
				File: gentypes.File{
					HeaderText: `// Copyright 2020 Acme Inc.
// All rights reserved.
//
// Licensed under "Only for testing purposes" license.
`,
					Package: gentypes.PackageRef{
						Name: "pkgdriver",
						Path: "app.dev/pkg/pkdriver",
					},
				},
				EndpointSets: []EndpointSet{
					{
						Service: Service{
							Object: service.Obj(),
							Type:   service.Underlying().(*types.Interface),
						},
						WithOpenCensus: true,
					},
				},
			}

			expected, err := os.ReadFile(fmt.Sprintf("./testdata/generator/%s/endpoint/zz_generated.endpoint.go", test.name))
			require.NoError(t, err)

			actual, err := Generate(file)
			require.NoError(t, err)

			assert.Equal(t, string(expected), string(actual), "the generated code does not match the expected")
		})
	}
}

func TestGenerate_CustomModule(t *testing.T) {
	pkgs, err := loader.LoadRoots("./testdata/generator/custom_module")
	require.NoError(t, err)

	pkg := pkgs[0]

	pkg.NeedTypesInfo()

	service := pkg.Types.Scope().Lookup("Service").Type().(*types.Named)

	file := File{
		File: gentypes.File{
			HeaderText: `// Copyright 2020 Acme Inc.
// All rights reserved.
//
// Licensed under "Only for testing purposes" license.
`,
			Package: gentypes.PackageRef{
				Name: "pkgdriver",
				Path: "app.dev/pkg/pkdriver",
			},
		},
		EndpointSets: []EndpointSet{
			{
				Service: Service{
					Object: service.Obj(),
					Type:   service.Underlying().(*types.Interface),
				},
				ModuleName:     "path.to.custom_module",
				WithOpenCensus: true,
			},
		},
	}

	expected, err := os.ReadFile("./testdata/generator/custom_module/endpoint/zz_generated.endpoint.go")
	require.NoError(t, err)

	actual, err := Generate(file)
	require.NoError(t, err)

	assert.Equal(t, string(expected), string(actual), "the generated code does not match the expected")
}

func TestGenerate_MultipleServices(t *testing.T) {
	pkgs, err := loader.LoadRoots("./testdata/generator/multiple_services")
	require.NoError(t, err)

	pkg := pkgs[0]

	pkg.NeedTypesInfo()

	service := pkg.Types.Scope().Lookup("Service").Type().(*types.Named)
	otherService := pkg.Types.Scope().Lookup("OtherService").Type().(*types.Named)
	another := pkg.Types.Scope().Lookup("Another").Type().(*types.Named)

	file := File{
		File: gentypes.File{
			HeaderText: `// Copyright 2020 Acme Inc.
// All rights reserved.
//
// Licensed under "Only for testing purposes" license.
`,
			Package: gentypes.PackageRef{
				Name: "pkgdriver",
				Path: "app.dev/pkg/pkdriver",
			},
		},
		EndpointSets: []EndpointSet{
			{
				Service: Service{
					Object: service.Obj(),
					Type:   service.Underlying().(*types.Interface),
				},
				WithOpenCensus: true,
			},
			{
				Service: Service{
					Object: otherService.Obj(),
					Type:   otherService.Underlying().(*types.Interface),
				},
				WithOpenCensus: true,
			},
			{
				Service: Service{
					Object: another.Obj(),
					Type:   another.Underlying().(*types.Interface),
				},
				WithOpenCensus: true,
			},
		},
	}

	expected, err := os.ReadFile("./testdata/generator/multiple_services/endpoint/zz_generated.endpoint.go")
	require.NoError(t, err)

	actual, err := Generate(file)
	require.NoError(t, err)

	assert.Equal(t, string(expected), string(actual), "the generated code does not match the expected")
}

func TestGenerate_ServiceError(t *testing.T) {
	pkgs, err := loader.LoadRoots("./testdata/generator/service_error")
	require.NoError(t, err)

	pkg := pkgs[0]

	pkg.NeedTypesInfo()

	service := pkg.Types.Scope().Lookup("Service").Type().(*types.Named)

	file := File{
		File: gentypes.File{
			HeaderText: `// Copyright 2020 Acme Inc.
// All rights reserved.
//
// Licensed under "Only for testing purposes" license.
`,
			Package: gentypes.PackageRef{
				Name: "pkgdriver",
				Path: "app.dev/pkg/pkdriver",
			},
		},
		EndpointSets: []EndpointSet{
			{
				Service: Service{
					Object: service.Obj(),
					Type:   service.Underlying().(*types.Interface),
				},
				WithOpenCensus: true,
				ErrorStrategy:  "service",
			},
		},
	}

	expected, err := os.ReadFile("./testdata/generator/service_error/endpoint/zz_generated.endpoint.go")
	require.NoError(t, err)

	actual, err := Generate(file)
	require.NoError(t, err)

	assert.Equal(t, string(expected), string(actual), "the generated code does not match the expected")
}
