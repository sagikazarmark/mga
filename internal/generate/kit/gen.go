package kit

import (
	"go/ast"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

// EndpointGenerator generates a Go kit Endpoint for a service.
type EndpointGenerator struct{}

func (g EndpointGenerator) RegisterMarkers(into *markers.Registry) error {
	if err := into.Register(endpointMarker); err != nil {
		return err
	}

	into.AddHelp(
		endpointMarker,
		markers.SimpleHelp("Kit", "enables endpoint generation for a service interface"),
	)

	return nil
}

func (g EndpointGenerator) Generate(ctx *genall.GenerationContext) error {
	for _, root := range ctx.Roots {
		g.parsePackage(ctx, root)
	}

	return nil
}

func (g EndpointGenerator) parsePackage(ctx *genall.GenerationContext, root *loader.Package) {
	ctx.Checker.Check(root, func(node ast.Node) bool {
		// ignore non-interfaces
		_, isIface := node.(*ast.InterfaceType)

		return isIface
	})

	err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
		typeMarker := info.Markers.Get(endpointMarker.Name)
		if typeMarker != nil {

		}
	})
	if err != nil {
		root.AddError(err)

		return
	}
}
