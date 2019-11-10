// +build gengo

package dispatcher

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"

	"sagikazarmark.dev/mga/internal/utils/packageutil"
)

// InterfaceSpec describes the event dispatcher interface.
type InterfaceSpec struct {
	Name    string
	Package PackageSpec
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
	pkgMap, err := packageutil.DirsToPackageMap(dir)
	if err != nil {
		return InterfaceSpec{}, err
	}

	builder := parser.New()
	if err := builder.AddDir(dir); err != nil {
		return InterfaceSpec{}, err
	}

	typs, err := builder.FindTypes()
	if err != nil {
		return InterfaceSpec{}, err
	}

	for _, pkgName := range builder.FindPackages() {
		pkg := typs[pkgName]

		if !pkg.Has(interfaceName) {
			continue
		}

		typ := pkg.Type(interfaceName)

		if typ.Kind != types.Interface {
			return InterfaceSpec{}, fmt.Errorf("%q is not an interface", interfaceName)
		}

		spec := InterfaceSpec{
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
			method := typ.Methods[name]

			methodSpec := MethodSpec{
				Name: name,
			}

			if len(method.Signature.Parameters) < 1 || len(method.Signature.Parameters) > 2 {
				return spec, fmt.Errorf(
					"dispatcher method %q can only have one or two parameters, but it has %d",
					name,
					len(method.Signature.Parameters),
				)
			}

			firstParam := method.Signature.Parameters[0]

			if len(method.Signature.Parameters) > 1 {
				if firstParam.Name.Package != "context" {
					return spec, fmt.Errorf("dispatcher method %q has two parameters, but the first one is not a context", name)
				}

				methodSpec.ReceivesContext = true

				secondParam := method.Signature.Parameters[1]

				pkgName := filepath.Base(secondParam.Name.Package)
				pkgPath := secondParam.Name.Package
				if secondParam.Name.Path != "" {
					pkgName = secondParam.Name.Package
					pkgPath = secondParam.Name.Path
				}

				if p, ok := pkgMap[pkgName]; ok {
					pkgPath = p
				}

				methodSpec.Event.Name = secondParam.Name.Name
				methodSpec.Event.Package.Name = pkgName
				methodSpec.Event.Package.Path = pkgPath
			} else {
				pkgName := filepath.Base(firstParam.Name.Package)
				pkgPath := firstParam.Name.Package
				if firstParam.Name.Path != "" {
					pkgName = firstParam.Name.Package
					pkgPath = firstParam.Name.Path
				}

				if p, ok := pkgMap[pkgName]; ok {
					pkgPath = p
				}

				methodSpec.Event.Name = firstParam.Name.Name
				methodSpec.Event.Package.Name = pkgName
				methodSpec.Event.Package.Path = pkgPath
			}

			if len(method.Signature.Results) > 1 {
				return spec, fmt.Errorf(
					"dispatcher method %q can only have one or zero return values, but it has %d",
					name,
					len(method.Signature.Results),
				)
			}

			if len(method.Signature.Results) == 1 {
				res := method.Signature.Results[0]

				if res.Name.Name != "error" {
					return spec, fmt.Errorf(
						"the return value in dispatcher method %q can only be error, but it is %q",
						name,
						res.Name.Name,
					)
				}

				methodSpec.ReturnsError = true
			}

			spec.Methods = append(spec.Methods, methodSpec)
		}

		return spec, nil
	}

	return InterfaceSpec{}, errors.New("interface not found")
}
