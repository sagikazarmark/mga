package test

// Run: go:generate go run sagikazarmark.dev/mga generate event

// +mga:event:handler
type Event struct {
	ID string
}
