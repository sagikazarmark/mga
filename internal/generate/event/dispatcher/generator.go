package dispatcher

import (
	"bytes"
	"strings"

	"github.com/dave/jennifer/jen"
)

func cleanEventDispatcherName(name string) string {
	name = strings.TrimSuffix(name, "Events")
	name = strings.TrimSuffix(name, "EventBus")
	name = strings.TrimSuffix(name, "EventDispatcher")

	return name
}

// Generate generates an event dispatcher.
func Generate(pkg string, spec InterfaceSpec) (string, error) {
	name := cleanEventDispatcherName(spec.Name)

	eventBusTypeName := name + "EventBus"
	eventDispatcherTypeName := name + "EventDispatcher"

	const (
		eventBusVarName = "bus"
	)

	file := jen.NewFilePath(pkg)

	// TODO: better comment
	// TODO: add version
	file.PackageComment("Code generated with mga")

	for _, method := range spec.Methods {
		file.ImportName(method.Event.Package.Path, method.Event.Package.Name)
	}

	file.Commentf("%s is a generic event bus.", eventBusTypeName)
	file.Type().Id(eventBusTypeName).Interface(
		jen.Comment("Publish sends an event to the underlying message bus."),
		jen.Id("Publish").Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("event").Interface(),
		).Error(),
	).Line()

	file.Commentf("%s dispatches events through the underlying generic event bus.", eventDispatcherTypeName)
	file.Type().Id(eventDispatcherTypeName).Struct(
		jen.Id(eventBusVarName).Id(eventBusTypeName),
	).Line()

	file.Commentf("New%s returns a new %s instance.", eventDispatcherTypeName, eventDispatcherTypeName)
	file.Func().
		Id("New" + eventDispatcherTypeName).
		Params(jen.Id(eventBusVarName).Id(eventBusTypeName)).
		Id(eventDispatcherTypeName).
		Block(
			jen.Return(
				jen.Op("&").Id(eventDispatcherTypeName).Values(jen.Dict{
					jen.Id(eventBusVarName): jen.Id(eventBusVarName),
				}),
			),
		).
		Line()

	for _, method := range spec.Methods {
		var params []jen.Code

		if method.ReceivesContext {
			params = append(params, jen.Id("ctx").Qual("context", "Context"))
		}

		params = append(params, jen.Id("event").Qual(method.Event.Package.Path, method.Event.Name))

		file.Commentf("%s dispatches a(n) %s event.", method.Name, method.Event.Name)
		fn := file.Func().Params(
			jen.Id("d").Id(eventDispatcherTypeName),
		).Id(method.Name).Params(params...)

		if method.ReturnsError {
			fn = fn.Error()
		}

		var block []jen.Code

		if !method.ReceivesContext {
			block = append(block, jen.Id("ctx").Op(":=").Qual("context", "Background").Call())
		}

		if method.ReturnsError {
			block = append(
				block,
				jen.Err().Op(":=").Id("d").Dot(eventBusVarName).Dot("Publish").Call(
					jen.Id("ctx"),
					jen.Id("event"),
				),
				jen.If(
					jen.Err().Op("!=").Nil(),
				).Block(
					jen.Return(jen.Qual("emperror.dev/errors", "WithDetails").Call(
						jen.Qual("emperror.dev/errors", "WithMessage").Call(
							jen.Err(),
							jen.Lit("failed to dispatch event"),
						),
						jen.Lit("event"), jen.Lit(method.Event.Name),
					)),
				),
				jen.Line(),
				jen.Return(jen.Nil()),
			)
		} else {
			block = append(block, jen.Id("d").Dot("eventBus").Dot("Publish").Call(
				jen.Id("ctx"),
				jen.Id("event"),
			))
		}

		fn = fn.Block(block...).Line()
	}

	var buf bytes.Buffer

	err := file.Render(&buf)

	return buf.String(), err
}
