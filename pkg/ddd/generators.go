package ddd

import (
	"time"

	"github.com/google/uuid"
)

// NewTimestamp gets the current timestamp
func NewTimestamp() int64 {
	return time.Now().Unix()
}

// NewUUID gets a new UUID
func NewUUID() string {
	return uuid.New().String()
}
