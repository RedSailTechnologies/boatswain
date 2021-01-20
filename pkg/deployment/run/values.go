package run

// Status is an alias to give us the status of a Step
type Status string

const (
	// NotStarted signifies a step not yet executed
	NotStarted Status = "NotStarted"

	// InProgress signifies a step that has been started
	InProgress Status = "InProgress"

	// Failed signifies a step that failed
	Failed Status = "Failed"

	// Succeeded signifies a step that was successful
	Succeeded Status = "Succeeded"

	// Skipped signifies a step that was not run
	Skipped Status = "Skipped"
)

// LogLevel is an alias for different log levels for run messages
// NOTE: this is not for internal logging, but rather how we represent deployment logs
type LogLevel string

const (
	// Debug represents a log level not normally shown, for debug purposes
	Debug LogLevel = "DEBUG"

	// Info represents informational logs
	Info LogLevel = "INFO"

	// Warn represents warning logs for conditions not normally seen, but which aren't errors
	Warn LogLevel = "WARN"

	// Error represents the log for an error
	Error LogLevel = "ERROR"
)
