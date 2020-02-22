package different_package

import (
	"context"

	"sagikazarmark.dev/mga/internal/generate/kit/endpoint/testdata/generator/different_package/svctype"
)

type Service interface {
	CreateTodo(ctx context.Context, text svctype.Text) (todo svctype.Todo, err error)
}
