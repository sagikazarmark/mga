package generics

import (
	"context"
)

type Optional[T any] struct{

}

type Service interface {
	CreateTodo(context.Context, string) (Optional[string, string], error)
}
