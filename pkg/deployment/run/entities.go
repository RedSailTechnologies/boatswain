package run

// A Step represents the outcome of a step being executed
type Step struct {
	Name   string
	Status Status
	Start  int64
	Stop   int64
	Logs   []Log
}

func (s *Step) log(m string, l LogLevel, t int64) {
	s.Logs = append(s.Logs, Log{
		Timestamp: t,
		Level:     l,
		Message:   m,
	})
}

// A Link is a reference attached to a run
type Link struct {
	Name string
	URL  string
}

// A Log is part of a step's result
type Log struct {
	Timestamp int64
	Level     LogLevel
	Message   string
}

// RuntimeError is thrown when some runtime execution rule is broken
type RuntimeError struct {
	m string
}

func (e RuntimeError) Error() string {
	return e.m
}
