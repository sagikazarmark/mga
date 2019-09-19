package test

import (
	"context"
)

type Event struct {
	ID string
}

//go:generate go run sagikazarmark.dev/mga generate event dispatcher --from Events --outdir .
type Events interface {
	// Event dispatches an Event event.
	Event(event Event)

	// EventWithContext accepts a context too.
	EventWithContext(ctx context.Context, event Event)

	// EventWithContextAndError combines the features of EventWithContext and EventWithError.
	EventWithContextAndError(ctx context.Context, event Event) error

	// EventWithError returns an error when something goes wrong during dispatching the event.
	EventWithError(event Event) error
}
