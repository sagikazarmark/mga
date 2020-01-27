package test

import (
	"context"
)

// go:generate go run sagikazarmark.dev/mga generate kit endpoint --outdir . --with-oc Service
// +kit:endpoint
type Service interface {
	Call(ctx context.Context, req interface{}) (interface{}, error)
	OtherCall(ctx context.Context, req interface{}) (interface{}, error)
	AnotherCall(ctx context.Context, req interface{}) (interface{}, error)
}

// +kit:endpoint
type AnotherService interface {
	Call(ctx context.Context, req interface{}) (interface{}, error)
	OtherCall(ctx context.Context, req interface{}) (interface{}, error)
	AnotherCall(ctx context.Context, req interface{}) (interface{}, error)
}
