package ddd

// Event is the basic interface all events implement
type Event interface {
	IsEvent()
}
