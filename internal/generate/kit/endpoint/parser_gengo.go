// +build gengo

package endpoint

import (
	"errors"
	"fmt"
	"sort"

	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"

	"sagikazarmark.dev/mga/internal/utils/packageutil"
)

// ServiceSpec describes the service interface.
type ServiceSpec struct {
	Name      string
	Package   PackageSpec
	Endpoints []EndpointSpec
}

// EndpointSpec describes a dispatcher method in an event dispatcher.
// nolint: golint
type EndpointSpec struct {
	Name string
}

// PackageSpec contains import information.
type PackageSpec struct {
	Name string
	Path string
}

// Parse parses a given package, looks for an interface and returns it as a normalized structure.
func Parse(dir string, interfaceName string) (ServiceSpec, error) {
	pkgMap, err := packageutil.DirsToPackageMap(dir)
	if err != nil {
		return ServiceSpec{}, err
	}

	builder := parser.New()
	if err := builder.AddDir(dir); err != nil {
		return ServiceSpec{}, err
	}

	typs, err := builder.FindTypes()
	if err != nil {
		return ServiceSpec{}, err
	}

	for _, pkgName := range builder.FindPackages() {
		pkg := typs[pkgName]

		if !pkg.Has(interfaceName) {
			continue
		}

		typ := pkg.Type(interfaceName)

		if typ.Kind != types.Interface {
			return ServiceSpec{}, fmt.Errorf("%q is not an interface", interfaceName)
		}

		spec := ServiceSpec{
			Name: interfaceName,
			Package: PackageSpec{
				Name: pkg.Name,
				Path: pkgMap[pkg.Name],
			},
		}

		methodNames := make([]string, 0, len(typ.Methods))
		for name := range typ.Methods {
			methodNames = append(methodNames, name)
		}
		sort.Strings(methodNames)

		for _, name := range methodNames {
			endpointSpec := EndpointSpec{
				Name: name,
			}

			spec.Endpoints = append(spec.Endpoints, endpointSpec)
		}

		return spec, nil
	}

	return ServiceSpec{}, errors.New("interface not found")
}
