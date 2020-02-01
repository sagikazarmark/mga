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

	RequestParameters  []gentypes.Argument
	ResponseParameters []gentypes.Argument
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

		method := ServiceMethod{
			Name: m.Name(),
		}

		sig := m.Type().(*types.Signature)

		if sig.Params().Len() < 1 {
			return svc, fmt.Errorf(
				"service method %q needs to have at least a context argument",
				m.Name(),
			)
		}

		if sig.Variadic() {
			return svc, fmt.Errorf(
				"variadic service method %q is not supported at the moment",
				m.Name(),
			)
		}

		firstParam := sig.Params().At(0)
		firstParamType, ok := firstParam.Type().(*types.Named)
		if !ok {
			return svc, fmt.Errorf("parameter %q in service method %q is not a valid type", firstParam.Name(), m.Name())
		}

		if firstParamType.Obj().Pkg().Path() != "context" {
			return svc, fmt.Errorf("first parameter of service method %q must be a context", m.Name())
		}

		numParams := sig.Params().Len()
		for i := 1; i < numParams; i++ {
			param := sig.Params().At(i)

			argument := gentypes.Argument{
				Name: param.Name(),
				Type: gentypes.TypeRef{
					Name: param.Type().String(),
				},
			}

			paramType, ok := param.Type().(*types.Named)
			if ok {
				argument.Type.Name = paramType.Obj().Name()
				argument.Type.Package.Name = paramType.Obj().Pkg().Name()
				argument.Type.Package.Path = paramType.Obj().Pkg().Path()
			}

			method.RequestParameters = append(method.RequestParameters, argument)
		}

		if sig.Results().Len() < 1 {
			return svc, fmt.Errorf(
				"service method %q needs to have at least one result",
				m.Name(),
			)
		}

		lastParam := sig.Results().At(sig.Results().Len() - 1)
		lastParamType, ok := lastParam.Type().(*types.Named)
		if !ok {
			return svc, fmt.Errorf("parameter %q in service method %q is not a valid type", lastParam.Name(), m.Name())
		}

		if lastParamType.Obj().Name() != "error" {
			return svc, fmt.Errorf("last return parameter of service method %q must be an error", m.Name())
		}

		numResults := sig.Results().Len()
		for i := 0; i < numResults-1; i++ {
			param := sig.Results().At(i)

			argument := gentypes.Argument{
				Name: param.Name(),
				Type: gentypes.TypeRef{
					Name: param.Type().String(),
				},
			}

			paramType, ok := param.Type().(*types.Named)
			if ok {
				argument.Type.Name = paramType.Obj().Name()
				argument.Type.Package.Name = paramType.Obj().Pkg().Name()
				argument.Type.Package.Path = paramType.Obj().Pkg().Path()
			}

			method.ResponseParameters = append(method.ResponseParameters, argument)
		}

		svc.Methods = append(svc.Methods, method)
	}

	return svc, nil
}
