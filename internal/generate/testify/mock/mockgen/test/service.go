package test

import (
	"context"
)

// nolint: godox
// Todo is a note describing a task to be done.
type Todo struct {
	ID   string
	Text string
	Done bool
}

//go:generate mga gen mockery --name Service --inpkg
// +testify:mock
type Service interface {
	// CreateTodo adds a new todo to the todo list.
	CreateTodo(ctx context.Context, text string) (id string, err error)

	// ListTodos returns the list of todos.
	ListTodos(ctx context.Context) ([]Todo, error)

	// MarkAsDone marks a todo as done.
	MarkAsDone(ctx context.Context, id string) error

	// TouchTodo records work on a todo.
	TouchTodo(ctx context.Context, id string)
}
