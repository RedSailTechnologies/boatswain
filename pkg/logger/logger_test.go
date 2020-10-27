package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitCreatesSingleton(t *testing.T) {
	// the logger should be initialized already, on package init
	assert.NotNil(t, l)
	initial := l
	initializeLogger()
	after := l
	assert.Equal(t, &initial, &after)
}

func TestDebug(t *testing.T) {
	Debug("message")
	Debug("message", "with", "params")
}

func TestInfo(t *testing.T) {
	Info("message")
	Info("message", "with", "params")
}

func TestWarn(t *testing.T) {
	Warn("message")
	Warn("message", "with", "params")
}

func TestError(t *testing.T) {
	Error("message")
	Error("message", "with", "params")
}

func TestPanic(t *testing.T) {
	assert.Panics(t, func() { Panic("message") })
	assert.Panics(t, func() { Panic("message", "with", "params") })
}
