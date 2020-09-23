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
	mockMarker = markers.Must(markers.MakeDefinition("testify:mock", markers.DescribesType, Marker{}))
)

// Marker enables generating a mock for an interface and provides information to the generator.
type Marker struct {
	// TestOnly tells the generator to write the generated mock in a test file.
	TestOnly bool `marker:"testOnly,optional"`

	// External tells the generator to write the generated mock in an "external" test package
	// (package name suffixed with _test).
	// External also implies TestOnly, but setting both values will generate the mock twice:
	// once in the same package in a test file, once in an external package.
	External bool `marker:"external,optional"`
}

// Generator generates a Go kit Endpoint for a service.
type Generator struct {
	// HeaderFile specifies the header text (e.g. license) to prepend to generated files.
	HeaderFile string `marker:",optional"`

	// Year specifies the year to substitute for " YEAR" in the header file.
	Year string `marker:",optional"`
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
		g.generatePackage(ctx, headerText, root)
	}

	return nil
}

func (g Generator) generatePackage(ctx *genall.GenerationContext, headerText string, root *loader.Package) {
	ctx.Checker.Check(root, func(node ast.Node) bool {
		// ignore non-interfaces
		// _, isIface := node.(*ast.InterfaceType)

		return true
	})

	root.NeedTypesInfo()

	var interfaces []mock.Interface
	var testOnlyInterfaces []mock.Interface
	var externalInterfaces []mock.Interface

	err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
		marker, ok := info.Markers.Get(mockMarker.Name).(Marker)
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

		iface := mock.Interface{
			Object: named.Obj(),
			Type:   named.Underlying().(*types.Interface),
		}

		switch {
		case marker.External:
			externalInterfaces = append(externalInterfaces, iface)

			if marker.TestOnly {
				testOnlyInterfaces = append(testOnlyInterfaces, iface)
			}

		case marker.TestOnly:
			testOnlyInterfaces = append(testOnlyInterfaces, iface)

		default:
			interfaces = append(interfaces, iface)
		}
	})
	if err != nil {
		root.AddError(err)

		return
	}

	packageName, packagePath := root.Name, root.PkgPath
	if pkgrefer, ok := ctx.OutputRule.(genutils.PackageRefer); ok {
		packageName, packagePath = pkgrefer.PackageRef(root)
	}

	if len(externalInterfaces) > 0 {
		file := mock.File{
			File: gentypes.File{
				Package: gentypes.PackageRef{
					Name: packageName + "_test",
					Path: packagePath + "_test", // See https://github.com/dave/jennifer/issues/73
				},
				HeaderText: headerText,
			},
			Interfaces: externalInterfaces,
		}

		outContents, err := mock.Generate(file)
		if err != nil {
			root.AddError(err)

			return
		}

		if outContents != nil {
			writeOut(ctx, root, outContents, "zz_generated.mock_external_test.go")
		}
	}

	if len(testOnlyInterfaces) > 0 {
		file := mock.File{
			File: gentypes.File{
				Package: gentypes.PackageRef{
					Name: packageName,
					Path: packagePath,
				},
				HeaderText: headerText,
			},
			Interfaces: testOnlyInterfaces,
		}

		outContents, err := mock.Generate(file)
		if err != nil {
			root.AddError(err)

			return
		}

		if outContents != nil {
			writeOut(ctx, root, outContents, "zz_generated.mock_test.go")
		}
	}

	if len(interfaces) > 0 {
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

			return
		}

		if outContents != nil {
			writeOut(ctx, root, outContents, "zz_generated.mock.go")
		}
	}
}

// writeOut outputs the given code.
func writeOut(ctx *genall.GenerationContext, root *loader.Package, outBytes []byte, fileName string) {
	outputFile, err := ctx.Open(root, fileName)
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
