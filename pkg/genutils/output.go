package genutils

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

// nolint: gochecknoinits
func init() {
	for ruleName, rule := range allOutputRules {
		ruleMarker := markers.Must(markers.MakeDefinition("output:"+ruleName, markers.DescribesPackage, rule))
		if err := optionsRegistry.Register(ruleMarker); err != nil {
			panic(err)
		}
		if helpGiver, hasHelp := rule.(genall.HasHelp); hasHelp {
			if help := helpGiver.Help(); help != nil {
				optionsRegistry.AddHelp(ruleMarker, help)
			}
		}
	}
}

// nolint: gochecknoglobals
var (
	allOutputRules = map[string]genall.OutputRule{
		"dir":    genall.OutputToDirectory(""),
		"none":   genall.OutputToNothing,
		"stdout": genall.OutputToStdout,
		"pkg":    OutputPackage{},
		"subpkg": OutputSubpackage{},
	}

	optionsRegistry = &markers.Registry{}
)

// LookupOutput looks up an output rule based on an output string.
func LookupOutput(output string) (genall.OutputRule, error) {
	outDef := optionsRegistry.Lookup("+output:"+output, markers.DescribesPackage)
	if outDef == nil {
		return nil, errors.New("invalid output option")
	}

	val, err := outDef.Parse("+output:" + output)
	if err != nil {
		return nil, err
	}

	outputRule, ok := val.(genall.OutputRule)
	if !ok {
		return nil, errors.New("no output rule found for this output option")
	}

	return outputRule, nil
}

// PackageRefer returns package reference (name and path) based on the output rule.
type PackageRefer interface {
	PackageRef(pkg *loader.Package) (string, string)
}

// +controllertools:marker:generateHelp:category=""

// OutputPackage outputs artifacts to the original package location.
type OutputPackage struct{}

func (o OutputPackage) Open(pkg *loader.Package, itemPath string) (io.WriteCloser, error) {
	if len(pkg.CompiledGoFiles) == 0 {
		return nil, fmt.Errorf("cannot output to a package with no path on disk")
	}

	outDir := filepath.Join(filepath.Dir(pkg.CompiledGoFiles[0]), filepath.Dir(itemPath))

	return genall.OutputToDirectory(outDir).Open(pkg, filepath.Base(itemPath))
}

func (o OutputPackage) PackageRef(pkg *loader.Package) (string, string) {
	return pkg.Name, pkg.PkgPath
}

// +controllertools:marker:generateHelp:category=""

// OutputSubpackage outputs artifacts to a subpackage of the original package.
type OutputSubpackage struct {
	// Prefix added to the package name.
	Prefix string `marker:",optional"`

	// Package is the main package name. When empty, the original package name is used.
	Package string `marker:",optional"`

	// Suffix added to the package name.
	Suffix string `marker:",optional"`
}

func (o OutputSubpackage) Open(pkg *loader.Package, itemPath string) (io.WriteCloser, error) {
	if len(pkg.CompiledGoFiles) == 0 {
		return nil, fmt.Errorf("cannot output to a package with no path on disk")
	}

	outDir := filepath.Join(filepath.Dir(pkg.CompiledGoFiles[0]), o.packageName(pkg))

	return genall.OutputToDirectory(outDir).Open(pkg, filepath.Base(itemPath))
}

func (o OutputSubpackage) PackageRef(pkg *loader.Package) (string, string) {
	pkgName := o.packageName(pkg)

	return pkgName, filepath.Join(pkg.PkgPath, pkgName)
}

func (o OutputSubpackage) packageName(pkg *loader.Package) string {
	pkgName := o.Package
	if pkgName == "" {
		pkgName = pkg.Name
	}

	return o.Prefix + pkgName + o.Suffix
}
