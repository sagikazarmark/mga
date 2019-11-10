// +build !gengo

package handler

import (
	"errors"
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
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
		return EventSpec{}, err
	}

	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(eventName)
		if obj == nil {
			continue
		}

		_, ok := obj.Type().Underlying().(*types.Struct)
		if !ok {
			return EventSpec{}, fmt.Errorf("%q is not a struct", eventName)
		}

		spec := EventSpec{
			Name: eventName,
			Package: PackageSpec{
				Name: obj.Pkg().Name(),
				Path: obj.Pkg().Path(),
			},
		}

		return spec, nil
	}

	return EventSpec{}, errors.New("event not found")
}
