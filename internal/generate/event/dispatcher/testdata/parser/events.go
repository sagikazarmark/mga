package parser

import (
	"context"

	"sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/imports"
)

// Event is something that happened at a given point in time.
type Event struct {
	ID string
}

// ImportedAliasedEvent is a type aliased event.
type ImportedAliasedEvent = imports.ImportedEvent

// Events dispatches Event events.
type Events interface {
	EmbeddedEvents
	embeddedEvents
	imports.ImportedEvents

	// Event dispatches an Event event.
	Event(event Event)

	// EventWithContext accepts a context too.
	EventWithContext(ctx context.Context, event Event)

	// EventWithContextAndError combines the features of EventWithContext and EventWithError.
	EventWithContextAndError(ctx context.Context, event Event) error

	// EventWithError returns an error when something goes wrong during dispatching the event.
	EventWithError(event Event) error

	// ImportedAliasedEvent dispatches an aliased event type.
	ImportedAliasedEvent(ctx context.Context, event ImportedAliasedEvent) error

	// ImportedEventDispatch dispatches an imported event type.
	ImportedEventDispatch(ctx context.Context, event imports.ImportedEvent) error
}

// EmbeddedEvents is an exported interface that will be embedded in the main Events interface.
type EmbeddedEvents interface {
	// EventEmbedded dispatches an Event event from another interface.
	EventEmbedded(ctx context.Context, event Event) error
}

// embeddedEvents is an unexported interface that will be embedded in the main Events interface.
type embeddedEvents interface {
	// EventStillEmbedded dispatches an Event event from another interface.
	EventEmbeddedFromUnexportedInterface(ctx context.Context, event Event) error
}
