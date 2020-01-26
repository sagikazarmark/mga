package endpoint

import (
	"errors"
	"fmt"
	"go/types"
	"strings"

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
func Parse(dir string, interfaceName string) (PackageDefinition, error) {
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
		return PackageDefinition{}, err
	}

	for _, pkg := range pkgs {
		obj := pkg.Types.Scope().Lookup(interfaceName)
		if obj == nil {
			continue
		}

		def := PackageDefinition{
			HeaderText:   "",
			PackageName:  pkg.Name+"driver",
			LogicalName:  pkg.Name,
			EndpointSets: nil,
		}

		setDef, err := parseInterface(obj)
		if err != nil {
			return def, err
		}

		def.EndpointSets = append(def.EndpointSets, setDef)

		return def, nil
	}

	return PackageDefinition{}, errors.New("interface not found")
}

func parseInterface(obj types.Object) (SetDefinition, error) {
	iface, ok := obj.Type().Underlying().(*types.Interface)
	if !ok {
		return SetDefinition{}, fmt.Errorf("%q is not an interface", obj.Name())
	}

	def := SetDefinition{
		BaseName: strings.TrimSuffix(obj.Name(), "Service"),
		Service: ServiceDefinition{
			Name:        obj.Name(),
			PackageName: obj.Pkg().Name(),
			PackagePath: obj.Pkg().Path(),
		},
		Endpoints:      nil,
		WithOpenCensus: false,
	}

	for i := 0; i < iface.NumMethods(); i++ {
		m := iface.Method(i)

		endpointSpec := EndpointDefinition{
			Name: m.Name(),
		}

		def.Endpoints = append(def.Endpoints, endpointSpec)
	}

	return def, nil
}
