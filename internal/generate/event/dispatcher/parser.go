package dispatcher

import (
	"errors"
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// InterfaceSpec describes the event dispatcher interface.
type InterfaceSpec struct {
	Name    string
	Methods []MethodSpec
}

// MethodSpec describes a dispatcher method in an event dispatcher.
type MethodSpec struct {
	Name            string
	Event           TypeSpec
	ReceivesContext bool
	ReturnsError    bool
}

// TypeSpec represents a named type with import information.
type TypeSpec struct {
	Name    string
	Package PackageSpec
}

// PackageSpec contains import information.
type PackageSpec struct {
	Name string
	Path string
}

// Parse parses a given package, looks for an interface and returns it as a normalized structure.
func Parse(dir string, interfaceName string) (InterfaceSpec, error) {
	cfg := &packages.Config{
		Mode:  packages.LoadSyntax,
		Tests: false,
	}

	pkgs, err := packages.Load(cfg, dir)
	if err != nil {
		return InterfaceSpec{}, err
	}

	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(interfaceName)
		if obj == nil {
			continue
		}

		iface, ok := obj.Type().Underlying().(*types.Interface)
		if !ok {
			return InterfaceSpec{}, fmt.Errorf("%q is not an interface", interfaceName)
		}

		spec := InterfaceSpec{
			Name: interfaceName,
		}

		for i := 0; i < iface.NumMethods(); i++ {
			m := iface.Method(i)

			methodSpec := MethodSpec{
				Name: m.Name(),
			}

			sig := m.Type().(*types.Signature)

			if sig.Params().Len() < 1 || sig.Params().Len() > 2 {
				return spec, fmt.Errorf("dispatcher method %q can only have one or two parameters, but it has %d", m.Name(), sig.Params().Len())
			}

			firstParam := sig.Params().At(0)
			firstParamType, ok := firstParam.Type().(*types.Named)
			if !ok {
				return spec, fmt.Errorf("parameter %q in dispatcher method %q is not a valid type", firstParam.Name(), m.Name())
			}

			if sig.Params().Len() > 1 {
				if firstParamType.Obj().Pkg().Path() != "context" {
					return spec, fmt.Errorf("dispatcher method %q has two parameters, but the first one is not a context", m.Name())
				}

				methodSpec.ReceivesContext = true

				secondParam := sig.Params().At(1)
				secondParamType, ok := secondParam.Type().(*types.Named)
				if !ok {
					return spec, fmt.Errorf("parameter %q in dispatcher method %q is not a valid type", secondParam.Name(), m.Name())
				}

				methodSpec.Event.Name = secondParamType.Obj().Name()
				methodSpec.Event.Package.Name = secondParamType.Obj().Pkg().Name()
				methodSpec.Event.Package.Path = secondParamType.Obj().Pkg().Path()
			} else {
				methodSpec.Event.Name = firstParamType.Obj().Name()
				methodSpec.Event.Package.Name = firstParamType.Obj().Pkg().Name()
				methodSpec.Event.Package.Path = firstParamType.Obj().Pkg().Path()
			}

			if sig.Results().Len() > 1 {
				return spec, fmt.Errorf("dispatcher method %q can only have one or zero return values, but it has %d", m.Name(), sig.Results().Len())
			}

			if sig.Results().Len() == 1 {
				res := sig.Results().At(0)
				resType, ok := res.Type().(*types.Named)
				if !ok {
					return spec, fmt.Errorf("result %q in dispatcher method %q is not a valid type", res.Name(), m.Name())
				}

				if resType.Obj().Name() != "error" {
					return spec, fmt.Errorf("the return value in dispatcher method %q can only be error, but it is %q", m.Name(), resType.Obj().Name())
				}

				methodSpec.ReturnsError = true
			}

			spec.Methods = append(spec.Methods, methodSpec)
		}

		return spec, nil
	}

	return InterfaceSpec{}, errors.New("interface not found")
}
