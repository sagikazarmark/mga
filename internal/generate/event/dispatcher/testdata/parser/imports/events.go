package imports

import (
	"context"
)

// ImportedEvent is something that happened at a given point in time.
type ImportedEvent struct {
	ID string
}

// ImportedEvents dispatches ImportedEvent events and is imported in the main event parser test.
type ImportedEvents interface {
	// ImportedEvent dispatches an ImportedEvent event.
	ImportedEvent(ctx context.Context, event ImportedEvent) error
}
