package handler

// EventHandlerFromEvent creates an EventHandler from an Event.
// nolint: golint
func EventHandlerFromEvent(event Event) EventHandler {
	return EventHandler{
		Name:  event.Name,
		Event: event,
	}
}
