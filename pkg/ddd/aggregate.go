package ddd

type Aggregate interface {
	Destroyed() bool
	Events() []Event
	UUID() string
}
