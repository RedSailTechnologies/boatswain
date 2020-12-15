package health

// ReadyService is the interface services must implement to be part of
// a normal ready check when running in kubernetes
type ReadyService interface {
	Ready() error
}
