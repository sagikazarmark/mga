package test

import (
	"context"

	"github.com/dgrijalva/jwt-go"

	oldtodo "sagikazarmark.dev/mga/internal/generate/testify/mock/mockgen/test/subpkg"
)

// nolint: godox
// Todo is a note describing a task to be done.
type Todo struct {
	ID   string
	Text string
	Done bool
}

//go:generate mga gen mockery --name Service
// +testify:mock
type Service interface {
	// CreateTodo adds a new todo to the todo list.
	CreateTodo(ctx context.Context, text string) (id string, err error)

	// ListTodos returns the list of todos.
	ListTodos(ctx context.Context) ([]Todo, error)

	// MarkAsDone marks a todo as done.
	MarkAsDone(ctx context.Context, id string) error

	// TouchTodo records work on a todo.
	TouchTodo(ctx context.Context, id oldtodo.ID)

	// ImportOldTodo imports a todo from the old format.
	ImportOldTodo(ctx context.Context, oldTodo oldtodo.OldTodo) (id string, err error)
}

//go:generate mga gen mockery --name Service2
// +testify:mock:testOnly=true
type Service2 interface {
	// CreateTodo adds a new todo to the todo list.
	CreateTodo(ctx context.Context, text string) (id string, err error)

	// ListTodos returns the list of todos.
	ListTodos(ctx context.Context) ([]Todo, error)

	// MarkAsDone marks a todo as done.
	MarkAsDone(ctx context.Context, id string) error

	// TouchTodo records work on a todo.
	TouchTodo(ctx context.Context, id string)

	// ImportOldTodo imports a todo from the old format.
	ImportOldTodo(ctx context.Context, oldTodo oldtodo.OldTodo) (id string, err error)
}

//go:generate mga gen mockery --name Service3
// +testify:mock
// from: https://github.com/sagikazarmark/mga/issues/34
type Service3 interface {
	// nolint: lll
	Refresh(ctx context.Context, refreshToken string, deviceID string, userName string, jwtToken *jwt.Token) (string, string, error)
}

//go:generate mga gen mockery --name Service4UnnamedParametersAndResults
// +testify:mock
// from https://github.com/sagikazarmark/mga/pull/42 #1.
type Service4UnnamedParametersAndResults interface {
	NamedParametersAndResults(isEnabled bool, count int, name string) (values []string, owner string, err error)

	UnnamedParameter(bool) (values []string, owner string, err error)

	UnnamedParameters(bool, int, string) (values []string, owner string, err error)

	UnnamedParametersAndResults(bool, int, string) ([]string, string, error)

	UnnamedResult(isEnabled bool, count int, name string) error

	UnnamedResults(isEnabled bool, count int, name string) ([]string, string, error)
}

//go:generate mga gen mockery --name Service5VariadicParameters
// +testify:mock
// from https://github.com/sagikazarmark/mga/pull/42 #2.
type Service5VariadicParameters interface {
	Regular(id string, count int, arguments []interface{}) (err error)

	Variadic(id string, count int, arguments ...interface{}) (err error)
}
