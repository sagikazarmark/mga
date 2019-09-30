package test

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func MakeCallEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}
func MakeOtherCallEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}

func MakeAnotherCallEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}
