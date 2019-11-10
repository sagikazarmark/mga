package packageutil

import (
	"golang.org/x/tools/go/packages"
)

// DirsToPackageMap parses packages in a list of directories and returns a map of them.
func DirsToPackageMap(dir string) (map[string]string, error) {
	pkgs, err := packages.Load(nil, dir)
	if err != nil {
		return nil, err
	}

	pkgMap := make(map[string]string, len(pkgs))

	for _, pkg := range pkgs {
		pkgMap[pkg.Name] = pkg.PkgPath
	}

	return pkgMap, nil
}
