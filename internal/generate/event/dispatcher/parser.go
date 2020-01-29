package dispatcher

import (
	"errors"
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"

	"sagikazarmark.dev/mga/internal/generate/gentypes"
)

// Events describes the event dispatcher interface.
type Events struct {
	gentypes.TypeRef

	Methods []EventMethod
}

// EventMethod describes a dispatcher method in an event dispatcher.
type EventMethod struct {
	Name            string
	Event           gentypes.TypeRef
	ReceivesContext bool
	ReturnsError    bool
}

// Parse parses a given package, looks for an interface and returns it as a normalized structure.
func Parse(dir string, interfaceName string) (Events, error) {
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
		return Events{}, err
	}

	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(interfaceName)
		if obj == nil {
			continue
		}

		return ParseEvents(obj)
	}

	return Events{}, errors.New("interface not found")
}

// ParseEvents parses an object as an event dispatcher.
func ParseEvents(obj types.Object) (Events, error) {
	iface, ok := obj.Type().Underlying().(*types.Interface)
	if !ok {
		return Events{}, fmt.Errorf("%q is not an interface", obj.Name())
	}

	events := Events{
		TypeRef: gentypes.TypeRef{
			Name: obj.Name(),
			Package: gentypes.PackageRef{
				Name: obj.Pkg().Name(),
				Path: obj.Pkg().Path(),
			},
		},
	}

	for i := 0; i < iface.NumMethods(); i++ {
		m := iface.Method(i)

		method := EventMethod{
			Name: m.Name(),
		}

		sig := m.Type().(*types.Signature)

		if sig.Params().Len() < 1 || sig.Params().Len() > 2 {
			return events, fmt.Errorf(
				"dispatcher method %q can only have one or two parameters, but it has %d",
				m.Name(),
				sig.Params().Len(),
			)
		}

		firstParam := sig.Params().At(0)
		firstParamType, ok := firstParam.Type().(*types.Named)
		if !ok {
			return events, fmt.Errorf("parameter %q in dispatcher method %q is not a valid type", firstParam.Name(), m.Name())
		}

		if sig.Params().Len() > 1 {
			if firstParamType.Obj().Pkg().Path() != "context" {
				return events, fmt.Errorf("dispatcher method %q has two parameters, but the first one is not a context", m.Name())
			}

			method.ReceivesContext = true

			secondParam := sig.Params().At(1)
			secondParamType, ok := secondParam.Type().(*types.Named)
			if !ok {
				return events, fmt.Errorf("parameter %q in dispatcher method %q is not a valid type", secondParam.Name(), m.Name())
			}

			method.Event.Name = secondParamType.Obj().Name()
			method.Event.Package.Name = secondParamType.Obj().Pkg().Name()
			method.Event.Package.Path = secondParamType.Obj().Pkg().Path()
		} else {
			method.Event.Name = firstParamType.Obj().Name()
			method.Event.Package.Name = firstParamType.Obj().Pkg().Name()
			method.Event.Package.Path = firstParamType.Obj().Pkg().Path()
		}

		if sig.Results().Len() > 1 {
			return events, fmt.Errorf(
				"dispatcher method %q can only have one or zero return values, but it has %d",
				m.Name(),
				sig.Results().Len(),
			)
		}

		if sig.Results().Len() == 1 {
			res := sig.Results().At(0)
			resType, ok := res.Type().(*types.Named)
			if !ok {
				return events, fmt.Errorf("result %q in dispatcher method %q is not a valid type", res.Name(), m.Name())
			}

			if resType.Obj().Name() != "error" {
				return events, fmt.Errorf(
					"the return value in dispatcher method %q can only be error, but it is %q",
					m.Name(),
					resType.Obj().Name(),
				)
			}

			method.ReturnsError = true
		}

		events.Methods = append(events.Methods, method)
	}

	return events, nil
}
