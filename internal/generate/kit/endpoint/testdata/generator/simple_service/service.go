package simple_service

import (
	"context"
)

type Service interface {
	Call(ctx context.Context, param string) (id string, err error)
}
