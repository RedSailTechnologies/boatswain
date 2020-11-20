package ddd

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func TestNewTimestamp(t *testing.T) {
	assert.NotNil(t, time.Unix(NewTimestamp(), 0))
}

func TestNewGuid(t *testing.T) {
	_, err := uuid.Parse(NewUUID())
	assert.Nil(t, err)
}
