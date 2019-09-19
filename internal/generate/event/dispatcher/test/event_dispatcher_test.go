package test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/subscriber"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setUpPublisher(t *testing.T) (*cqrs.EventBus, <-chan *message.Message) {
	publisher := gochannel.NewGoChannel(gochannel.Config{}, watermill.NopLogger{})
	const topic = "event"
	eventBus, err := cqrs.NewEventBus(publisher, func(_ string) string { return topic }, &cqrs.JSONMarshaler{})
	require.NoError(t, err)

	messages, err := publisher.Subscribe(context.Background(), topic)
	require.NoError(t, err)

	return eventBus, messages
}

func TestEventDispatcher_Event(t *testing.T) {
	eventBus, messages := setUpPublisher(t)

	events := NewEventDispatcher(eventBus)

	event := Event{
		ID: "id",
	}

	events.Event(event)

	received, all := subscriber.BulkRead(messages, 1, time.Second)
	if !all {
		t.Fatal("no message received")
	}

	var receivedEvent Event

	err := json.Unmarshal(received[0].Payload, &receivedEvent)
	require.NoError(t, err)

	assert.Equal(t, event, receivedEvent)
}

func TestEventDispatcher_EventWithContext(t *testing.T) {
	eventBus, messages := setUpPublisher(t)

	events := NewEventDispatcher(eventBus)

	event := Event{
		ID: "id",
	}

	events.EventWithContext(context.Background(), event)

	received, all := subscriber.BulkRead(messages, 1, time.Second)
	if !all {
		t.Fatal("no message received")
	}

	var receivedEvent Event

	err := json.Unmarshal(received[0].Payload, &receivedEvent)
	require.NoError(t, err)

	assert.Equal(t, event, receivedEvent)
}

func TestEventDispatcher_EventWithContextAndError(t *testing.T) {
	eventBus, messages := setUpPublisher(t)

	events := NewEventDispatcher(eventBus)

	event := Event{
		ID: "id",
	}

	err := events.EventWithContextAndError(context.Background(), event)
	require.NoError(t, err)

	received, all := subscriber.BulkRead(messages, 1, time.Second)
	if !all {
		t.Fatal("no message received")
	}

	var receivedEvent Event

	err = json.Unmarshal(received[0].Payload, &receivedEvent)
	require.NoError(t, err)

	assert.Equal(t, event, receivedEvent)
}

func TestEventDispatcher_EventWithError(t *testing.T) {
	eventBus, messages := setUpPublisher(t)

	events := NewEventDispatcher(eventBus)

	event := Event{
		ID: "id",
	}

	err := events.EventWithError(event)
	require.NoError(t, err)

	received, all := subscriber.BulkRead(messages, 1, time.Second)
	if !all {
		t.Fatal("no message received")
	}

	var receivedEvent Event

	err = json.Unmarshal(received[0].Payload, &receivedEvent)
	require.NoError(t, err)

	assert.Equal(t, event, receivedEvent)
}

type failingEventBus struct {
	err error
}

func (f failingEventBus) Publish(ctx context.Context, event interface{}) error {
	return f.err
}

func TestEventDispatcher_EventWithContextAndError_Error(t *testing.T) {
	failure := errors.NewPlain("error")
	eventBus := failingEventBus{err: failure}

	events := NewEventDispatcher(eventBus)

	event := Event{
		ID: "id",
	}

	err := events.EventWithContextAndError(context.Background(), event)

	assert.Equal(t, failure, errors.Cause(err))
}

func TestEventDispatcher_EventWithError_Error(t *testing.T) {
	failure := errors.NewPlain("error")
	eventBus := failingEventBus{err: failure}

	events := NewEventDispatcher(eventBus)

	event := Event{
		ID: "id",
	}

	err := events.EventWithError(event)

	assert.Equal(t, failure, errors.Cause(err))
}
