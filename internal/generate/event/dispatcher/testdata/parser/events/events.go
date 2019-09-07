package events

import (
	"context"
)

// Event is something that happened at a given point in time.
type Event struct {
	ID string
}

// Events dispatches Event events.
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
