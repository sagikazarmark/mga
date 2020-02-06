package simple_service

import (
	"context"
)

type Service interface {
	CreateTodo(ctx context.Context, text string) (id string, err error)
}
