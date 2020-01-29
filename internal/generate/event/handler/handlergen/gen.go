package handlergen

import (
	"fmt"
	"go/ast"
	"go/types"
	"io"
	"strings"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"

	"sagikazarmark.dev/mga/internal/generate/event/handler"
	"sagikazarmark.dev/mga/internal/generate/gentypes"
)

// nolint: gochecknoglobals
var (
	handlerMarker = markers.Must(markers.MakeDefinition("mga:event:handler", markers.DescribesType, struct{}{}))
)

// Generator generates a Go kit Endpoint for a service.
type Generator struct {
	// HeaderFile specifies the header text (e.g. license) to prepend to generated files.
	HeaderFile string `marker:",optional"`

	// Year specifies the year to substitute for " YEAR" in the header file.
	Year string `marker:",optional"`
}

func (g Generator) RegisterMarkers(into *markers.Registry) error {
	if err := into.Register(handlerMarker); err != nil {
		return err
	}

	into.AddHelp(
		handlerMarker,
		markers.SimpleHelp("Kit", "enables event handler generation for an event"),
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
		_, isStruct := node.(*ast.StructType)

		return isStruct
	})

	root.NeedTypesInfo()

	var eventHandlers []handler.EventHandler

	err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
		if marker := info.Markers.Get(handlerMarker.Name); marker == nil {
			return
		}

		typeInfo := root.TypesInfo.TypeOf(info.RawSpec.Name)
		if typeInfo == types.Typ[types.Invalid] {
			root.AddError(loader.ErrFromNode(fmt.Errorf("unknown type %s", info.Name), info.RawSpec))

			return
		}

		event, err := handler.ParseEvent(root.TypesInfo.ObjectOf(info.RawSpec.Name))
		if err != nil {
			root.AddError(err)

			return
		}

		eventHandlers = append(eventHandlers, handler.EventHandlerFromEvent(event))
	})
	if err != nil {
		root.AddError(err)

		return nil
	}

	if len(eventHandlers) == 0 {
		return nil
	}

	file := handler.File{
		File: gentypes.File{
			Package: gentypes.PackageRef{
				Name: root.Name + "gen",
				Path: root.PkgPath + "/" + root.Name + "gen",
			},
			HeaderText: headerText,
		},
		EventHandlers: eventHandlers,
	}

	outContents, err := handler.Generate(file)
	if err != nil {
		root.AddError(err)

		return nil
	}

	return outContents
}

// writeOut outputs the given code.
func writeOut(ctx *genall.GenerationContext, root *loader.Package, outBytes []byte) {
	outputFile, err := ctx.Open(root, fmt.Sprintf("%sgen/zz_generated.event_handler.go", root.Name))
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
