package endpoint

import (
	"go/types"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-tools/pkg/loader"

	"sagikazarmark.dev/mga/pkg/gentypes"
)

func TestGenerate_SimpleService(t *testing.T) {
	pkgs, err := loader.LoadRoots("./testdata/generator/simple_service")
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

	expected, err := ioutil.ReadFile("./testdata/generator/simple_service/endpoint/zz_generated.endpoint.go")
	require.NoError(t, err)

	actual, err := Generate(file)
	require.NoError(t, err)

	assert.Equal(t, string(expected), string(actual), "the generated code does not match the expected")
}

func TestGenerate_ServiceWithStruct(t *testing.T) {
	pkgs, err := loader.LoadRoots("./testdata/generator/service_with_struct")
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

	expected, err := ioutil.ReadFile("./testdata/generator/service_with_struct/endpoint/zz_generated.endpoint.go")
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

	expected, err := ioutil.ReadFile("./testdata/generator/multiple_services/endpoint/zz_generated.endpoint.go")
	require.NoError(t, err)

	actual, err := Generate(file)
	require.NoError(t, err)

	assert.Equal(t, string(expected), string(actual), "the generated code does not match the expected")
}

func TestGenerate_UnnamedParam(t *testing.T) {
	pkgs, err := loader.LoadRoots("./testdata/generator/unnamed_param")
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

	expected, err := ioutil.ReadFile("./testdata/generator/unnamed_param/endpoint/zz_generated.endpoint.go")
	require.NoError(t, err)

	actual, err := Generate(file)
	require.NoError(t, err)

	assert.Equal(t, string(expected), string(actual), "the generated code does not match the expected")
}
