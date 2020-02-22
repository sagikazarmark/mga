package mock

import (
	"fmt"
	"go/types"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-tools/pkg/loader"

	"sagikazarmark.dev/mga/pkg/gentypes"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name string
	}{

	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			pkgs, err := loader.LoadRoots(fmt.Sprintf("./testdata/generator/%s", test.name))
			require.NoError(t, err)

			pkg := pkgs[0]

			pkg.NeedTypesInfo()

			iface := pkg.Types.Scope().Lookup("Service").Type().(*types.Named)

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
				Interfaces: []Interface{
					{
						Object: iface.Obj(),
						Type:   iface.Underlying().(*types.Interface),
					},
				},
			}

			expected, err := ioutil.ReadFile(fmt.Sprintf("./testdata/generator/%s/endpoint/zz_generated.mock.go", test.name))
			require.NoError(t, err)

			actual, err := Generate(file)
			require.NoError(t, err)

			assert.Equal(t, string(expected), string(actual), "the generated code does not match the expected")
		})
	}
}
