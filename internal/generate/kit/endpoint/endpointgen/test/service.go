package test

import (
	"context"

	"github.com/dgrijalva/jwt-go"

	"sagikazarmark.dev/mga/internal/generate/kit/endpoint/endpointgen/test/svctypes"
)

// nolint: godox
// Todo is a note describing a task to be done.
type Todo struct {
	ID   string
	Text string
	Done bool
}

// +kit:endpoint
type Service interface {
	// CreateTodo adds a new todo to the todo list.
	CreateTodo(ctx context.Context, text string) (id string, err error)

	// ListTodos returns the list of todos.
	ListTodos(ctx context.Context) ([]Todo, error)

	// MarkAsDone marks a todo as done.
	MarkAsDone(ctx context.Context, id string) error
}

// +kit:endpoint
type Service2 interface {
	// CreateTodo adds a new todo to the todo list.
	CreateTodo(ctx context.Context, text svctypes.Text) (id svctypes.ID, err error)

	// ListTodos returns the list of todos.
	ListTodos(ctx context.Context) ([]svctypes.Todo, error)

	// MarkAsDone marks a todo as done.
	MarkAsDone(ctx context.Context, id svctypes.ID) error
}

// +kit:endpoint
// from: https://github.com/sagikazarmark/mga/issues/34
type Service3 interface {
	// nolint: lll
	Refresh(ctx context.Context, refreshToken string, deviceID string, userName string, jwtToken *jwt.Token) (string, string, error)
}
