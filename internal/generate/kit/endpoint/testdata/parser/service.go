package parser

import (
	"context"
)

// Service is an interface for an application use case.
type Service interface {
	Call(ctx context.Context, req interface{}) (interface{}, error)
}
