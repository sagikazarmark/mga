package test

// go:generate go run sagikazarmark.dev/mga generate event handler --outdir . Event
// +mga:event:handler
type Event struct {
	ID string
}
