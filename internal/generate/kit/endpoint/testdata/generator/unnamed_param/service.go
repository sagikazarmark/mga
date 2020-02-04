package unnamed_param

import (
	"context"
)

type Service interface {
	Call(context.Context, string) (string, error)
}
