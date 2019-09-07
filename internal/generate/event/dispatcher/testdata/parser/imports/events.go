package imports

import (
	"sagikazarmark.dev/mga/internal/generate/event/dispatcher/testdata/parser/embeds"
)

// Events dispatches Event events.
type Events interface {
	embeds.Events
}
