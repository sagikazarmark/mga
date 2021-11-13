package endpointgen

import (
	"fmt"
	"go/ast"
	"go/types"
	"io"
	"strings"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"

	"sagikazarmark.dev/mga/internal/generate/kit/endpoint"
	"sagikazarmark.dev/mga/pkg/gentypes"
	"sagikazarmark.dev/mga/pkg/genutils"
)

// nolint: gochecknoglobals
var (
	endpointMarker = markers.Must(markers.MakeDefinition("kit:endpoint", markers.DescribesType, Marker{}))
)

// +controllertools:marker:generateHelp:category=Kit

// Marker enables generating an endpoint for a service and provides information to the generator.
type Marker struct {
	// BaseName specifies a base name for the service (other than the one automatically generated).
	//
	// When not specified falls back to base name created from the service name.
	BaseName string `marker:"baseName,optional"`

	// ModuleName can be used instead of the package name in an operation name to uniquely identify a service call.
	//
	// Falls back to the package name.
	ModuleName string `marker:"moduleName,optional"`

	// WithOpenCensus enables generating a TraceEndpoint middleware.
	WithOpenCensus bool `marker:"withOpenCensus,optional"`

	// ErrorStrategy decides whether returned errors are checked for being endpoint or service errors.
	ErrorStrategy string `marker:"errorStrategy,optional"`
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

func (Generator) CheckFilter() loader.NodeFilter {
	return func(node ast.Node) bool {
		// ignore non-interfaces
		_, isIface := node.(*ast.InterfaceType)

		return isIface
	}
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
	ctx.Checker.Check(root)

	root.NeedTypesInfo()

	var endpointSets []endpoint.EndpointSet

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

		if !types.IsInterface(typeInfo) {
			root.AddError(loader.ErrFromNode(fmt.Errorf("%s is not an interface", info.Name), info.RawSpec))

			return
		}

		named, ok := typeInfo.(*types.Named)
		if !ok {
			root.AddError(loader.ErrFromNode(fmt.Errorf("%s is not a named type", info.Name), info.RawSpec))

			return
		}

		endpointSets = append(
			endpointSets,
			endpoint.EndpointSet{
				Service: endpoint.Service{
					Object: named.Obj(),
					Type:   named.Underlying().(*types.Interface),
				},
				ModuleName:     marker.ModuleName,
				WithOpenCensus: marker.WithOpenCensus,
				ErrorStrategy:  marker.ErrorStrategy,
			},
		)
	})
	if err != nil {
		root.AddError(err)

		return nil
	}

	if len(endpointSets) == 0 {
		return nil
	}

	packageName, packagePath := root.Name, root.PkgPath
	if pkgrefer, ok := ctx.OutputRule.(genutils.PackageRefer); ok {
		packageName, packagePath = pkgrefer.PackageRef(root)
	}

	file := endpoint.File{
		File: gentypes.File{
			Package: gentypes.PackageRef{
				Name: packageName,
				Path: packagePath,
			},
			HeaderText: headerText,
		},
		EndpointSets: endpointSets,
	}

	outContents, err := endpoint.Generate(file)
	if err != nil {
		root.AddError(err)

		return nil
	}

	return outContents
}

// writeOut outputs the given code.
func writeOut(ctx *genall.GenerationContext, root *loader.Package, outBytes []byte) {
	outputFile, err := ctx.Open(root, "zz_generated.endpoint.go")
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
