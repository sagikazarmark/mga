package dispatcher

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	input := InterfaceSpec{
		Name: "Events",
		Methods: []MethodSpec{
			{
				Name: "MarkedAsDone",
				Event: TypeSpec{
					Name: "MarkedAsDone",
					Package: PackageSpec{
						Path: "github.com/sagikazarmark/modern-go-application/internal/app/mga/todo",
						Name: "todo",
					},
				},
				ReceivesContext: true,
				ReturnsError:    true,
			},
			{
				Name: "MarkedAsDone2",
				Event: TypeSpec{
					Name: "MarkedAsDone2",
					Package: PackageSpec{
						Path: "github.com/sagikazarmark/modern-go-application/internal/app/mga/todo",
						Name: "todo",
					},
				},
				ReceivesContext: false,
				ReturnsError:    true,
			},
			{
				Name: "MarkedAsDone3",
				Event: TypeSpec{
					Name: "MarkedAsDone3",
					Package: PackageSpec{
						Path: "github.com/sagikazarmark/modern-go-application/internal/app/mga/todo",
						Name: "todo",
					},
				},
				ReceivesContext: true,
				ReturnsError:    false,
			},
		},
	}

	res, err := Generate("github.com/sagikazarmark/modern-go-application/internal/app/mga/todo/todogen", input)
	if err != nil {
		t.Fatal(err)
	}

	expected := `// Code generated with mga
package todogen

import (
	"context"
	errors "emperror.dev/errors"
	"github.com/sagikazarmark/modern-go-application/internal/app/mga/todo"
)

// EventBus is a generic event bus.
type EventBus interface {
	// Publish sends an event to the underlying message bus.
	Publish(ctx context.Context, event interface{}) error
}

// EventDispatcher dispatches events through the underlying generic event bus.
type EventDispatcher struct {
	bus EventBus
}

// NewEventDispatcher returns a new EventDispatcher instance.
func NewEventDispatcher(bus EventBus) EventDispatcher {
	return EventDispatcher{bus: bus}
}

// MarkedAsDone dispatches a(n) MarkedAsDone event.
func (d EventDispatcher) MarkedAsDone(ctx context.Context, event todo.MarkedAsDone) error {
	err := d.bus.Publish(ctx, event)
	if err != nil {
		return errors.WithDetails(errors.WithMessage(err, "failed to dispatch event"), "event", "MarkedAsDone")
	}

	return nil
}

// MarkedAsDone2 dispatches a(n) MarkedAsDone2 event.
func (d EventDispatcher) MarkedAsDone2(event todo.MarkedAsDone2) error {
	ctx := context.Background()
	err := d.bus.Publish(ctx, event)
	if err != nil {
		return errors.WithDetails(errors.WithMessage(err, "failed to dispatch event"), "event", "MarkedAsDone2")
	}

	return nil
}

// MarkedAsDone3 dispatches a(n) MarkedAsDone3 event.
func (d EventDispatcher) MarkedAsDone3(ctx context.Context, event todo.MarkedAsDone3) {
	_ = d.bus.Publish(ctx, event)
}
`

	if res != expected {
		t.Error("the generated code does not match the expected one")
	}
}
