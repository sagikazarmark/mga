package endpoint

import (
	"errors"
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"

	"sagikazarmark.dev/mga/pkg/gentypes"
)

// Service describes a service interface.
type Service struct {
	gentypes.TypeRef

	// Methods is the list of methods in the service.
	Methods []ServiceMethod
}

// ServiceMethod describes a method in a service.
type ServiceMethod struct {
	Name string
}

// Parse parses a given package, looks for an interface and returns it as a normalized structure.
func Parse(dir string, interfaceName string) (Service, error) {
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
		return Service{}, err
	}

	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(interfaceName)
		if obj == nil {
			continue
		}

		svc, err := ParseInterface(obj)
		if err != nil {
			return svc, err
		}

		return svc, nil
	}

	return Service{}, errors.New("interface not found")
}

// ParseInterface parses an object as a service interface.
func ParseInterface(obj types.Object) (Service, error) {
	iface, ok := obj.Type().Underlying().(*types.Interface)
	if !ok {
		return Service{}, fmt.Errorf("%q is not an interface", obj.Name())
	}

	svc := Service{
		TypeRef: gentypes.TypeRef{
			Name: obj.Name(),
			Package: gentypes.PackageRef{
				Name: obj.Pkg().Name(),
				Path: obj.Pkg().Path(),
			},
		},
		Methods: nil,
	}

	for i := 0; i < iface.NumMethods(); i++ {
		m := iface.Method(i)

		endpointSpec := ServiceMethod{
			Name: m.Name(),
		}

		svc.Methods = append(svc.Methods, endpointSpec)
	}

	return svc, nil
}
