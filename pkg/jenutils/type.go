package jenutils

import (
	"go/types"
	"strings"

	"github.com/dave/jennifer/jen"
)

// nolint: gochecknoglobals
var basicKindMap = map[types.BasicKind]func(statement *jen.Statement) *jen.Statement{
	types.Bool:       (*jen.Statement).Bool,
	types.Int:        (*jen.Statement).Int,
	types.Int8:       (*jen.Statement).Int8,
	types.Int16:      (*jen.Statement).Int16,
	types.Int32:      (*jen.Statement).Int32,
	types.Int64:      (*jen.Statement).Int64,
	types.Uint:       (*jen.Statement).Uint,
	types.Uint8:      (*jen.Statement).Uint8,
	types.Uint16:     (*jen.Statement).Uint16,
	types.Uint32:     (*jen.Statement).Uint32,
	types.Uint64:     (*jen.Statement).Uint64,
	types.Uintptr:    (*jen.Statement).Uintptr,
	types.Float32:    (*jen.Statement).Float32,
	types.Float64:    (*jen.Statement).Float64,
	types.Complex64:  (*jen.Statement).Complex64,
	types.Complex128: (*jen.Statement).Complex128,
	types.String:     (*jen.Statement).String,
}

// Type attaches a type to a statement based on a parsed type.
func Type(stmt *jen.Statement, t types.Type) jen.Code {
	switch t := t.(type) {
	case *types.Basic:
		f, ok := basicKindMap[t.Kind()]
		if !ok {
			panic("invalid basic kind: " + t.String())
		}

		return f(stmt)

	case *types.Array:
		return Type(stmt.Index(jen.Lit(t.Len())), t.Elem())

	case *types.Slice:
		return Type(stmt.Index(), t.Elem())

	case *types.Chan:
		if t.Dir() == types.RecvOnly {
			stmt = stmt.Op("<-")
		}
		stmt = stmt.Chan()
		if t.Dir() == types.SendOnly {
			stmt = stmt.Op("<-")
		}

		return Type(stmt, t.Elem())

	case *types.Map:
		return Type(stmt.Map(Type(&jen.Statement{}, t.Key())), t.Elem())

	case *types.Interface:
		return stmt.Interface()

	case *types.Named:
		if pkg := t.Obj().Pkg(); pkg != nil {
			return stmt.Qual(pkg.Path(), t.Obj().Name())
		}

		// builtin interfaces (eg. error) have no package
		return stmt.Id(t.Obj().Name())

	case *types.Pointer:
		return Type(stmt.Op("*"), t.Elem())

	case *types.Signature:
		stmt := stmt.Func()

		if t.Recv() != nil {
			stmt.Params(Type(jen.Id(t.Recv().Name()), t.Recv().Type()))
		}

		if t.Params() != nil {
			params := make([]jen.Code, t.Params().Len())
			for paramIndex := 0; paramIndex < t.Params().Len(); paramIndex++ {
				params[paramIndex] = Type(jen.Id(t.Params().At(paramIndex).Name()), t.Params().At(paramIndex).Type())
			}
			stmt.Params(params...)
		}

		if t.Results() != nil {
			results := make([]jen.Code, t.Results().Len())
			for resultIndex := 0; resultIndex < t.Results().Len(); resultIndex++ {
				results[resultIndex] = Type(jen.Id(t.Results().At(resultIndex).Name()), t.Results().At(resultIndex).Type())
			}
			stmt.Params(results...)
		}

		return stmt
	case *types.Struct:
		var fields []jen.Code

		for i := 0; i < t.NumFields()-1; i++ {
			field := t.Field(i)

			fields = append(fields, Type(jen.Id(field.Name()), field.Type()))
		}

		return stmt.Struct(fields...)
	}

	panic("unknown type: " + t.String())
}

// Import adds an import to the generated file when the type is a named type.
func Import(file *jen.File, typ types.Type) {
	switch t := typ.(type) {
	case *types.Basic, *types.Interface, *types.Struct:
		return

	case *types.Array:
		Import(file, t.Elem())

	case *types.Slice:
		Import(file, t.Elem())

	case *types.Chan:
		Import(file, t.Elem())

	case *types.Map:
		Import(file, t.Key())
		Import(file, t.Elem())

	case *types.Named:
		if pkg := t.Obj().Pkg(); pkg != nil {
			if pkg.Path() != pkg.Name() && // Internal packages always have the same path as the package name
				!strings.HasSuffix(pkg.Path(), "/"+pkg.Name()) { // Package name is different
				file.ImportAlias(pkg.Path(), pkg.Name())
			} else {
				file.ImportName(pkg.Path(), pkg.Name())
			}
		}

		// builtin interfaces (eg. error) have no package

	case *types.Pointer:
		Import(file, t.Elem())

	case *types.Signature:
		if t.Recv() != nil {
			Import(file, t.Recv().Type())
		}

		if t.Params() != nil {
			for paramIndex := 0; paramIndex < t.Params().Len(); paramIndex++ {
				Import(file, t.Params().At(paramIndex).Type())
			}
		}

		if t.Results() != nil {
			for resultIndex := 0; resultIndex < t.Results().Len(); resultIndex++ {
				Import(file, t.Results().At(resultIndex).Type())
			}
		}
	}
}

// IsNillable checks if a type is nillable. Useful for guarding type conversions.
func IsNillable(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Pointer, *types.Array, *types.Map, *types.Interface, *types.Signature, *types.Chan, *types.Slice:
		return true
	case *types.Named:
		return IsNillable(t.Underlying())
	}

	return false
}
