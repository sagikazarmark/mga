package test

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type eventHandlerStub struct {
	ctx   context.Context
	event Event
}

func (s *eventHandlerStub) Event(ctx context.Context, event Event) error {
	s.ctx = ctx
	s.event = event

	return nil
}

func TestEventEventHandler(t *testing.T) {
	handler := NewEventEventHandler(&eventHandlerStub{}, "event_handler")

	assert.Implements(t, (*cqrs.EventHandler)(nil), handler)
}

func TestEventEventHandler_HandlerName(t *testing.T) {
	handler := NewEventEventHandler(&eventHandlerStub{}, "event_handler")

	name := handler.HandlerName()

	assert.Equal(t, "event_handler", name)
}

func TestEventEventHandler_NewEvent(t *testing.T) {
	handler := NewEventEventHandler(&eventHandlerStub{}, "event_handler")

	event := handler.NewEvent()

	assert.IsType(t, &Event{}, event)
}

func TestEventEventHandler_Handle(t *testing.T) {
	h := &eventHandlerStub{}
	handler := NewEventEventHandler(h, "event_handler")

	ctx := context.Background()
	event := Event{
		ID: "1234",
	}

	err := handler.Handle(ctx, &event)
	require.NoError(t, err)

	assert.Equal(t, h.ctx, ctx)
	assert.Equal(t, h.event, event)
}
