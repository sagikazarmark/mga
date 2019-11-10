// +build gengo

package handler

import (
	"errors"
	"fmt"

	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"

	"sagikazarmark.dev/mga/internal/utils/packageutil"
)

// EventSpec describes the event struct.
type EventSpec struct {
	Name    string
	Package PackageSpec
}

// PackageSpec contains import information.
type PackageSpec struct {
	Name string
	Path string
}

// Parse parses a given package, looks for a struct and returns it as a normalized structure.
func Parse(dir string, eventName string) (EventSpec, error) {
	pkgMap, err := packageutil.DirsToPackageMap(dir)
	if err != nil {
		return EventSpec{}, err
	}

	builder := parser.New()
	if err := builder.AddDir(dir); err != nil {
		return EventSpec{}, err
	}

	typs, err := builder.FindTypes()
	if err != nil {
		return EventSpec{}, err
	}

	for _, pkgName := range builder.FindPackages() {
		pkg := typs[pkgName]

		if !pkg.Has(eventName) {
			continue
		}

		typ := pkg.Type(eventName)

		if typ.Kind != types.Struct {
			return EventSpec{}, fmt.Errorf("%q is not an struct", eventName)
		}

		spec := EventSpec{
			Name: eventName,
			Package: PackageSpec{
				Name: pkg.Name,
				Path: pkgMap[pkg.Name],
			},
		}

		return spec, nil
	}

	return EventSpec{}, errors.New("event not found")
}
