package test

//go:generate go run sagikazarmark.dev/mga generate event handler --outdir . Event
type Event struct {
	ID string
}
