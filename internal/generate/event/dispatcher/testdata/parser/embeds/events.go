package embeds

import (
	"context"
)

// Event is something that happened at a given point in time.
type Event struct {
	ID string
}

// Events dispatches Event events.
type Events interface {
	OtherEvents
	moreEvents

	// Event dispatches an Event event.
	Event(ctx context.Context, event Event) error
}

// OtherEvents is an exported interface that will be embedded in the main Events interface.
type OtherEvents interface {
	// OtherEvent dispatches an Event event from another interface.
	OtherEvent(ctx context.Context, event Event) error
}

// moreEvents is an unexported interface that will be embedded in the main Events interface.
type moreEvents interface {
	// MoreEvent dispatches an Event event from another interface.
	MoreEvent(ctx context.Context, event Event) error
}
