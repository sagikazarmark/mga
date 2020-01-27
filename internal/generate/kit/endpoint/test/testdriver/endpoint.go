package testdriver

import (
	"github.com/go-kit/kit/endpoint"

	"sagikazarmark.dev/mga/internal/generate/kit/endpoint/test"
)

func MakeCallEndpoint(_ test.Service) endpoint.Endpoint {
	return endpoint.Nop
}
func MakeOtherCallEndpoint(_ test.Service) endpoint.Endpoint {
	return endpoint.Nop
}

func MakeAnotherCallEndpoint(_ test.Service) endpoint.Endpoint {
	return endpoint.Nop
}

func MakeCallAnotherEndpoint(_ test.Service) endpoint.Endpoint {
	return endpoint.Nop
}
func MakeOtherCallAnotherEndpoint(_ test.Service) endpoint.Endpoint {
	return endpoint.Nop
}

func MakeAnotherCallAnotherEndpoint(_ test.Service) endpoint.Endpoint {
	return endpoint.Nop
}
