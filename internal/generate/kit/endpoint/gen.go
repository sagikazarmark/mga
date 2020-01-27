package endpoint

import (
	"fmt"
	"go/ast"
	"go/types"
	"io"
	"strings"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"

	"sagikazarmark.dev/mga/internal/generate/gentypes"
)

var (
	endpointMarker = markers.Must(markers.MakeDefinition("kit:endpoint", markers.DescribesType, Marker{}))
)

// +controllertools:marker:generateHelp:category=Kit

// Marker marker enables generating an endpoint for a service and provides information to the generator.
type Marker struct {
	// BaseName specifies a base name for the service (other than the one automatically generated).
	//
	// When not specified falls back to base name created from the service name.
	BaseName string `marker:"baseName,optional"`

	// ModuleName can be used instead of the package name as an operation name to uniquely identify a service call.
	//
	// Falls back to the package name.
	ModuleName string `marker:"moduleName,optional"`

	// WithOpenCensus enables generating a TraceEndpoint middleware.
	WithOpenCensus bool `marker:"withOpenCensus,optional"`
}

// Generator generates a Go kit Endpoint for a service.
type Generator struct {
	// HeaderFile specifies the header text (e.g. license) to prepend to generated files.
	HeaderFile string `marker:",optional"`

	// Year specifies the year to substitute for " YEAR" in the header file.
	Year string `marker:",optional"`
}

func (g Generator) RegisterMarkers(into *markers.Registry) error {
	if err := into.Register(endpointMarker); err != nil {
		return err
	}

	into.AddHelp(
		endpointMarker,
		markers.SimpleHelp("Kit", "enables endpoint generation for a service interface"),
	)

	return nil
}

func (g Generator) Generate(ctx *genall.GenerationContext) error {
	var headerText string

	if g.HeaderFile != "" {
		headerBytes, err := ctx.ReadFile(g.HeaderFile)
		if err != nil {
			return err
		}

		headerText = string(headerBytes)
	}

	headerText = strings.ReplaceAll(headerText, " YEAR", " "+g.Year)

	for _, root := range ctx.Roots {
		outContents := g.generatePackage(ctx, headerText, root)
		if outContents == nil {
			continue
		}

		writeOut(ctx, root, outContents)
	}

	return nil
}

func (g Generator) generatePackage(ctx *genall.GenerationContext, headerText string, root *loader.Package) []byte {
	ctx.Checker.Check(root, func(node ast.Node) bool {
		// ignore non-interfaces
		_, isIface := node.(*ast.InterfaceType)

		return isIface
	})

	root.NeedTypesInfo()

	var endpointSets []EndpointSet

	err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
		marker, ok := info.Markers.Get(endpointMarker.Name).(Marker)
		if !ok {
			return
		}

		typeInfo := root.TypesInfo.TypeOf(info.RawSpec.Name)
		if typeInfo == types.Typ[types.Invalid] {
			root.AddError(loader.ErrFromNode(fmt.Errorf("unknown type %s", info.Name), info.RawSpec))

			return
		}

		svc, err := parseInterface(root.TypesInfo.ObjectOf(info.RawSpec.Name))
		if err != nil {
			root.AddError(err)

			return
		}

		endpointSets = append(endpointSets, endpointSetFromService(svc, marker))
	})
	if err != nil {
		root.AddError(err)

		return nil
	}

	if len(endpointSets) == 0 {
		return nil
	}

	file := File{
		File: gentypes.File{
			Package: gentypes.PackageRef{
				Name: root.Name + "driver",
				Path: root.PkgPath + "/" + root.Name + "driver",
			},
			HeaderText: headerText,
		},
		EndpointSets: endpointSets,
	}

	outContents, err := Generate(file)
	if err != nil {
		root.AddError(err)

		return nil
	}

	return outContents
}

func endpointSetFromService(svc Service, marker Marker) EndpointSet {
	baseName := marker.BaseName
	withOpenCensus := marker.WithOpenCensus
	moduleName := marker.ModuleName

	if baseName == "" {
		baseName = strings.TrimSuffix(svc.Name, "Service")
	}

	endpointSet := EndpointSet{
		Name:           baseName,
		Service:        svc.TypeRef,
		Endpoints:      nil,
		WithOpenCensus: withOpenCensus,
	}

	if moduleName == "" {
		moduleName = svc.Package.Name
	}

	for _, method := range svc.Methods {
		var operationName string

		// if endpoint set name is empty, do not add it to the operation name
		if endpointSet.Name == "" {
			operationName = fmt.Sprintf("%s.%s", moduleName, method.Name)
		} else {
			operationName = fmt.Sprintf("%s.%s.%s", moduleName, endpointSet.Name, method.Name)
		}

		endpointSet.Endpoints = append(
			endpointSet.Endpoints,
			Endpoint{
				Name:          method.Name,
				OperationName: operationName,
			},
		)
	}

	return endpointSet
}

// writeOut outputs the given code.
func writeOut(ctx *genall.GenerationContext, root *loader.Package, outBytes []byte) {
	outputFile, err := ctx.Open(root, fmt.Sprintf("%sdriver/zz_generated.endpoint.go", root.Name))
	if err != nil {
		root.AddError(err)
		return
	}
	defer outputFile.Close()
	n, err := outputFile.Write(outBytes)
	if err != nil {
		root.AddError(err)
		return
	}
	if n < len(outBytes) {
		root.AddError(io.ErrShortWrite)
	}
}
