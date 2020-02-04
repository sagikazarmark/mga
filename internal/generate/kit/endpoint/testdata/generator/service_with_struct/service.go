package service_with_struct

import (
	"context"
)

type Request struct {
	Param string
}

type Response struct {
	ID string
}

type Service interface {
	Call(ctx context.Context, req Request) (response Response, err error)
}
