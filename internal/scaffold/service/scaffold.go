package service

import (
	"bytes"
	"go/format"

	"github.com/dave/jennifer/jen"
)

// Scaffold creates a new empty service interface.
func Scaffold(pkg string) ([]byte, error) {
	code := jen.NewFile(pkg)

	code.Comment("+kit:endpoint")
	code.Line()

	code.Comment("Service <insert your description>.")
	code.Type().Id("Service").Interface(jen.Comment("Insert your operations here"))

	code.Comment("NewService returns a new Service.")
	code.Func().Id("NewService").Params().Params(jen.Id("Service")).Block(
		jen.Return(jen.Id("service").Values()),
	)

	code.Type().Id("service").Struct()

	var buf bytes.Buffer

	err := code.Render(&buf)
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}
