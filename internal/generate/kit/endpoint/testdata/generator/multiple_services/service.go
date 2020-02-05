package multiple_services

import (
	"context"
)

type Service interface {
	CreateTodo(ctx context.Context, text string) (id string, err error)
}

type OtherService interface {
	CreateTodo(ctx context.Context) (err error)
}

type Another interface {
	CreateTodo(ctx context.Context) (err error)
}
