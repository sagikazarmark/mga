package test

import (
	"context"

	"github.com/golang-jwt/jwt/v5"

	oldtodo "sagikazarmark.dev/mga/internal/generate/testify/mock/mockgen/test/subpkg"
)

// nolint: godox
// Todo is a note describing a task to be done.
type Todo struct {
	ID   string
	Text string
	Done bool
}

// +testify:mock
//
//go:generate mga gen mockery --name Service
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

// +testify:mock:testOnly=true
//
//go:generate mga gen mockery --name Service2
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

// +testify:mock
// from: https://github.com/sagikazarmark/mga/issues/34
//
//go:generate mga gen mockery --name Service3
type Service3 interface {
	// nolint: lll
	Refresh(ctx context.Context, refreshToken string, deviceID string, userName string, jwtToken *jwt.Token) (string, string, error)
}

// +testify:mock
// from https://github.com/sagikazarmark/mga/pull/42 #1.
//
//go:generate mga gen mockery --name Service4UnnamedParametersAndResults
type Service4UnnamedParametersAndResults interface {
	NamedParametersAndResults(isEnabled bool, count int, name string) (values []string, owner string, err error)

	UnnamedParameter(bool) (values []string, owner string, err error)

	UnnamedParameters(bool, int, string) (values []string, owner string, err error)

	UnnamedParametersAndResults(bool, int, string) ([]string, string, error)

	UnnamedResult(isEnabled bool, count int, name string) error

	UnnamedResults(isEnabled bool, count int, name string) ([]string, string, error)
}

// +testify:mock
// from https://github.com/sagikazarmark/mga/pull/42 #2.
//
//go:generate mga gen mockery --name Service5VariadicParameters
type Service5VariadicParameters interface {
	Regular(id string, count int, arguments []interface{}) (err error)

	Variadic(id string, count int, arguments ...interface{}) (err error)
}

// +testify:mock
// from https://github.com/sagikazarmark/mga/pull/42 #3.
//
//go:generate mga gen mockery --name Service6FunctionParameters
type Service6FunctionParameters interface {
	FunctionParameter(id string, predicate func(id oldtodo.ID, todo oldtodo.OldTodo) bool, count int) (err error)

	FunctionParameters(
		id string,
		predicate func(id string, importedTodo oldtodo.OldTodo) bool,
		operation func(count int, importedID oldtodo.ID),
		count int,
	) (err error)
}

// +testify:mock:external=true
//
//go:generate mga gen mockery --name Service7
type Service7 interface {
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

// +testify:mock:external=true,testOnly=true
//
//go:generate mga gen mockery --name Service8
type Service8 interface {
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
