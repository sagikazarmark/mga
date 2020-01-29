package dispatcher

import (
	"strings"
)

// EventDispatcherFromEvents creates an EventDispatcher from Events.
// nolint: golint
func EventDispatcherFromEvents(events Events) EventDispatcher {
	return EventDispatcher{
		Name:              cleanEventDispatcherName(events.Name),
		DispatcherMethods: events.Methods,
	}
}

func cleanEventDispatcherName(name string) string {
	name = strings.TrimSuffix(name, "Events")
	name = strings.TrimSuffix(name, "EventBus")
	name = strings.TrimSuffix(name, "EventDispatcher")

	return name
}
