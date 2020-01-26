package endpoint

// PackageDefinition is a Go kit service driver package.
// It represents one or more services and provides information for generating endpoints for these services.
type PackageDefinition struct {
	// HeaderText is added as a comment to the top of the generated file, above any package comments.
	//
	// It is useful for adding license information to generated files.
	HeaderText string

	// PackageName is the name of the target Go package.
	PackageName string

	// PackagePath is the path of the target Go package.
	PackagePath string

	// EndpointSets represents endpoints to be generated for each service in the module.
	EndpointSets []EndpointSetDefinition
}

// EndpointSetDefinition represents endpoints for a single service.
// nolint: golint
type EndpointSetDefinition struct {
	// BaseName is a base name for the endpoint set.
	BaseName string

	// Service provides information about the service the endpoint set is generated for.
	Service ServiceDefinition

	// Endpoints is the list of endpoints represented by the set.
	Endpoints []EndpointDefinition

	// WithOpenCensus enables generating a trace middleware for the endpoint set.
	WithOpenCensus bool
}

// ServiceDefinition represents the service details for an endpoint set.
type ServiceDefinition struct {
	// Name is the original service name.
	Name string

	// PackageName is the name of the service package.
	PackageName string

	// PackagePath is the path to the service package.
	PackagePath string
}

// EndpointDefinition represents an endpoint for a single service call.
// nolint: golint
type EndpointDefinition struct {
	// Name identifies a call within a service.
	Name string

	// OperationName uniquely identifies a service call.
	OperationName string
}
