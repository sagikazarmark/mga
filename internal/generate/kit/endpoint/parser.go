package endpoint

import (
	"errors"
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
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
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedTypesSizes |
			packages.NeedSyntax |
			packages.NeedTypesInfo,
		Tests: false,
	}

	pkgs, err := packages.Load(cfg, dir)
	if err != nil {
		return ServiceSpec{}, err
	}

	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(interfaceName)
		if obj == nil {
			continue
		}

		iface, ok := obj.Type().Underlying().(*types.Interface)
		if !ok {
			return ServiceSpec{}, fmt.Errorf("%q is not an interface", interfaceName)
		}

		spec := ServiceSpec{
			Name: interfaceName,
			Package: PackageSpec{
				Name: obj.Pkg().Name(),
				Path: obj.Pkg().Path(),
			},
		}

		for i := 0; i < iface.NumMethods(); i++ {
			m := iface.Method(i)

			endpointSpec := EndpointSpec{
				Name: m.Name(),
			}

			spec.Endpoints = append(spec.Endpoints, endpointSpec)
		}

		return spec, nil
	}

	return ServiceSpec{}, errors.New("interface not found")
}
