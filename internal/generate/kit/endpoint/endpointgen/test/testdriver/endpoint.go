package testdriver

import (
	"github.com/go-kit/kit/endpoint"

	test2 "sagikazarmark.dev/mga/internal/generate/kit/endpoint/endpointgen/test"
)

func MakeCallEndpoint(_ test2.Service) endpoint.Endpoint {
	return endpoint.Nop
}
func MakeOtherCallEndpoint(_ test2.Service) endpoint.Endpoint {
	return endpoint.Nop
}

func MakeAnotherCallEndpoint(_ test2.Service) endpoint.Endpoint {
	return endpoint.Nop
}

func MakeCallAnotherEndpoint(_ test2.Service) endpoint.Endpoint {
	return endpoint.Nop
}
func MakeOtherCallAnotherEndpoint(_ test2.Service) endpoint.Endpoint {
	return endpoint.Nop
}

func MakeAnotherCallAnotherEndpoint(_ test2.Service) endpoint.Endpoint {
	return endpoint.Nop
}
