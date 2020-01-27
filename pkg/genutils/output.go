package genutils

import (
	"fmt"
	"io"
	"path/filepath"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
)

// +controllertools:marker:generateHelp:category=""

// OutputArtifacts outputs artifacts to different locations, depending on
// whether they're package-associated or not.
type OutputArtifacts struct{}

func (o OutputArtifacts) Open(pkg *loader.Package, itemPath string) (io.WriteCloser, error) {
	if len(pkg.CompiledGoFiles) == 0 {
		return nil, fmt.Errorf("cannot output to a package with no path on disk")
	}

	outDir := filepath.Join(filepath.Dir(pkg.CompiledGoFiles[0]), filepath.Dir(itemPath))

	return genall.OutputToDirectory(outDir).Open(pkg, filepath.Base(itemPath))
}
