package gentypes

// TypeRef is a reference to a type.
// It can be used generate code referring to a type, possibly in another package.
type TypeRef struct {
	// Name of the type.
	Name string

	// PackageName is a reference to the package where the type can be found.
	Package PackageRef
}

// PackageRef is a reference to a package.
// It can be used generate code referring to a package.
type PackageRef struct {
	// Name of the package.
	Name string

	// Path to the package.
	Path string
}

// File is an input to source code generation.
type File struct {
	// Package where the file belongs.
	Package PackageRef

	// HeaderText is added as a comment to the top of the generated file, above any package comments.
	//
	// It is useful for adding license information to generated files.
	HeaderText string
}
