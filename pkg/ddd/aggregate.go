package ddd

// Aggregate represents common functionality for aggregates mostly for unmarshaling
type Aggregate interface {
	Destroyed() bool
	Events() []Event
	UUID() string
}
