package dispatcher

import (
	"errors"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func Parse(dir string, ifaceName string) (Spec, error) {
	cfg := &packages.Config{
		Mode:  packages.LoadSyntax,
		Tests: false,
	}

	pkgs, err := packages.Load(cfg, dir)
	if err != nil {
		return Spec{}, err
	}

	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(ifaceName)
		if obj == nil {
			continue
		}

		iface, ok := obj.Type().Underlying().(*types.Interface)
		if !ok {
			return Spec{}, errors.New("the given name is not an interface")
		}

		spec := Spec{
			Package: pkg.PkgPath,
			Name:    ifaceName,
		}

		for i := 0; i < iface.NumMethods(); i++ {
			m := iface.Method(i)

			eventSpec := EventSpec{
				DispatchName: m.Name(),
			}

			sig := m.Type().(*types.Signature)

			if sig.Params().Len() < 1 || sig.Params().Len() > 2 {
				return spec, errors.New("method can only have one or two parameters")
			}

			firstParam := sig.Params().At(0).Type().(*types.Named)

			if sig.Params().Len() > 1 {
				if firstParam.Obj().Pkg().Path() != "context" {
					return spec, errors.New("first parameter must be a context when there are two parameters")
				}

				eventSpec.ReceivesContext = true

				secondParam := sig.Params().At(1).Type().(*types.Named)

				eventSpec.EventName = secondParam.Obj().Name()
				eventSpec.EventPackage = secondParam.Obj().Pkg().Path()
			} else {
				eventSpec.EventName = firstParam.Obj().Name()
				eventSpec.EventPackage = firstParam.Obj().Pkg().Path()
			}

			if sig.Results().Len() > 1 {
				return spec, errors.New("there can only be zero or one return values")
			}

			if sig.Results().Len() == 1 {
				res := sig.Results().At(0).Type().(*types.Named)
				if res.Obj().Name() != "error" {
					return spec, errors.New("the return value can only be error")
				}

				eventSpec.ReturnsError = true
			}

			spec.Events = append(spec.Events, eventSpec)
		}

		return spec, nil
	}

	return Spec{}, errors.New("interface not found")
}
