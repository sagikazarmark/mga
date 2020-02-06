package unnamed_param

import (
	"context"
)

type Service interface {
	CreateTodo(context.Context, string) (string, error)
}
