package mockgen

import (
	"fmt"
	"go/ast"
	"go/types"
	"io"
	"strings"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"

	"sagikazarmark.dev/mga/internal/generate/testify/mock"
	"sagikazarmark.dev/mga/pkg/gentypes"
	"sagikazarmark.dev/mga/pkg/genutils"
)

// nolint: gochecknoglobals
var (
	mockMarker = markers.Must(markers.MakeDefinition("testify:mock", markers.DescribesType, struct{}{}))
)

// Generator generates a Go kit Endpoint for a service.
type Generator struct {
	// HeaderFile specifies the header text (e.g. license) to prepend to generated files.
	HeaderFile string `marker:",optional"`

	// Year specifies the year to substitute for " YEAR" in the header file.
	Year string `marker:",optional"`

	// TestOnly tells the generator to put the generated mocks in a _test package.
	TestOnly bool `marker:",optional"`
}

func (g Generator) RegisterMarkers(into *markers.Registry) error {
	if err := into.Register(mockMarker); err != nil {
		return err
	}

	into.AddHelp(
		mockMarker,
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
		// _, isIface := node.(*ast.InterfaceType)

		return true
	})

	root.NeedTypesInfo()

	var interfaces []mock.Interface

	err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
		marker := info.Markers.Get(mockMarker.Name)
		if marker == nil {
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

		interfaces = append(
			interfaces,
			mock.Interface{
				Object: named.Obj(),
				Type:   named.Underlying().(*types.Interface),
			},
		)
	})
	if err != nil {
		root.AddError(err)

		return nil
	}

	if len(interfaces) == 0 {
		return nil
	}

	packageName, packagePath := root.Name, root.PkgPath
	if pkgrefer, ok := ctx.OutputRule.(genutils.PackageRefer); ok {
		packageName, packagePath = pkgrefer.PackageRef(root)
	}

	if g.TestOnly && !strings.HasSuffix(packageName, "_test") {
		packageName += "_test"
	}

	file := mock.File{
		File: gentypes.File{
			Package: gentypes.PackageRef{
				Name: packageName,
				Path: packagePath,
			},
			HeaderText: headerText,
		},
		Interfaces: interfaces,
	}

	outContents, err := mock.Generate(file)
	if err != nil {
		root.AddError(err)

		return nil
	}

	return outContents
}

// writeOut outputs the given code.
func writeOut(ctx *genall.GenerationContext, root *loader.Package, outBytes []byte) {
	outputFile, err := ctx.Open(root, "zz_generated.mock.go")
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
