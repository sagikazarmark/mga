package service_with_struct

import (
	"context"
)

type NewTodo struct {
	Text string
}

type CreatedTodo struct {
	ID string
}

type Service interface {
	CreateTodo(ctx context.Context, newTodo NewTodo) (response CreatedTodo, err error)
}
