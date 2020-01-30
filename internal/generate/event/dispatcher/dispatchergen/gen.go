package dispatchergen

import (
	"fmt"
	"go/ast"
	"go/types"
	"io"
	"strings"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"

	"sagikazarmark.dev/mga/internal/generate/event/dispatcher"
	"sagikazarmark.dev/mga/pkg/gentypes"
	"sagikazarmark.dev/mga/pkg/genutils"
)

// nolint: gochecknoglobals
var (
	dispatcherMarker = markers.Must(markers.MakeDefinition("mga:event:dispatcher", markers.DescribesType, struct{}{}))
)

// Generator generates a Go kit Endpoint for a service.
type Generator struct {
	// HeaderFile specifies the header text (e.g. license) to prepend to generated files.
	HeaderFile string `marker:",optional"`

	// Year specifies the year to substitute for " YEAR" in the header file.
	Year string `marker:",optional"`
}

func (g Generator) RegisterMarkers(into *markers.Registry) error {
	if err := into.Register(dispatcherMarker); err != nil {
		return err
	}

	into.AddHelp(
		dispatcherMarker,
		markers.SimpleHelp("Kit", "enables event dispatcher generation for events"),
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
		// _, isIface := node.(*ast.InterfaceType)

		return true
	})

	root.NeedTypesInfo()

	var eventDispatchers []dispatcher.EventDispatcher

	err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
		if marker := info.Markers.Get(dispatcherMarker.Name); marker == nil {
			return
		}

		typeInfo := root.TypesInfo.TypeOf(info.RawSpec.Name)
		if typeInfo == types.Typ[types.Invalid] {
			root.AddError(loader.ErrFromNode(fmt.Errorf("unknown type %s", info.Name), info.RawSpec))

			return
		}

		events, err := dispatcher.ParseEvents(root.TypesInfo.ObjectOf(info.RawSpec.Name))
		if err != nil {
			root.AddError(err)

			return
		}

		eventDispatchers = append(eventDispatchers, dispatcher.EventDispatcherFromEvents(events))
	})
	if err != nil {
		root.AddError(err)

		return nil
	}

	if len(eventDispatchers) == 0 {
		return nil
	}

	packageName, packagePath := root.Name, root.PkgPath
	if pkgrefer, ok := ctx.OutputRule.(genutils.PackageRefer); ok {
		packageName, packagePath = pkgrefer.PackageRef(root)
	}

	file := dispatcher.File{
		File: gentypes.File{
			Package: gentypes.PackageRef{
				Name: packageName,
				Path: packagePath,
			},
			HeaderText: headerText,
		},
		EventDispatchers: eventDispatchers,
	}

	outContents, err := dispatcher.Generate(file)
	if err != nil {
		root.AddError(err)

		return nil
	}

	return outContents
}

// writeOut outputs the given code.
func writeOut(ctx *genall.GenerationContext, root *loader.Package, outBytes []byte) {
	outputFile, err := ctx.Open(root, "zz_generated.event_dispatcher.go")
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
