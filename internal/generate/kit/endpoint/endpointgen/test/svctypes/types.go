package svctypes

type ID string

type Text string

// nolint: godox
// Todo is a note describing a task to be done.
type Todo struct {
	ID   string
	Text string
	Done bool
}
