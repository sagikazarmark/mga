package endpoint

import (
	"fmt"
	"strings"
)

// GeneratorOptions provides additional information to the generator for generating an endpoint from a service.
type GeneratorOptions struct {
	// BaseName specifies a base name for the service (other than the one automatically generated).
	//
	// When not specified falls back to base name created from the service name.
	BaseName string

	// ModuleName can be used instead of the package name as an operation name to uniquely identify a service call.
	//
	// Falls back to the package name.
	ModuleName string

	// WithOpenCensus enables generating a TraceEndpoint middleware.
	WithOpenCensus bool
}

// EndpointSetFromService creates an EndpointSet from a Service.
// nolint: golint
func EndpointSetFromService(svc Service, options GeneratorOptions) EndpointSet {
	baseName := options.BaseName
	withOpenCensus := options.WithOpenCensus
	moduleName := options.ModuleName

	if baseName == "" {
		baseName = strings.TrimSuffix(svc.Name, "Service")
	}

	endpointSet := EndpointSet{
		Name:           baseName,
		Service:        svc.TypeRef,
		Endpoints:      nil,
		WithOpenCensus: withOpenCensus,
	}

	if moduleName == "" {
		moduleName = svc.Package.Name
	}

	for _, method := range svc.Methods {
		var operationName string

		// if endpoint set name is empty, do not add it to the operation name
		if endpointSet.Name == "" {
			operationName = fmt.Sprintf("%s.%s", moduleName, method.Name)
		} else {
			operationName = fmt.Sprintf("%s.%s.%s", moduleName, endpointSet.Name, method.Name)
		}

		endpointSet.Endpoints = append(
			endpointSet.Endpoints,
			Endpoint{
				Name:          method.Name,
				OperationName: operationName,
			},
		)
	}

	return endpointSet
}
