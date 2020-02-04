package multiple_services

import (
	"context"
)

type Service interface {
	Call(ctx context.Context, param string) (id string, err error)
}

type OtherService interface {
	Call(ctx context.Context) (err error)
}

type Another interface {
	Call(ctx context.Context) (err error)
}
