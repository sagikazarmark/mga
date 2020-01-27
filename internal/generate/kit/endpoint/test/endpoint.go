package test

import (
	"github.com/go-kit/kit/endpoint"
)

func MakeCallEndpoint(_ Service) endpoint.Endpoint {
	return endpoint.Nop
}
func MakeOtherCallEndpoint(_ Service) endpoint.Endpoint {
	return endpoint.Nop
}

func MakeAnotherCallEndpoint(_ Service) endpoint.Endpoint {
	return endpoint.Nop
}
