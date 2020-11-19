package ddd

import "fmt"

// DestroyedError represents an error when subsequent commands are called on a destroyed cluster
type DestroyedError struct {
	Entity string
}

func (err DestroyedError) Error() string {
	return fmt.Sprintf("%s has been destroyed", err.Entity)
}

// IDError represents an error with id mismatches, not set, etc.
type IDError struct{}

func (err IDError) Error() string {
	return "id error"
}

// InvalidArgumentError represents an error with an argument out of range, etc.
type InvalidArgumentError struct {
	Arg string
	Val string
}

func (err InvalidArgumentError) Error() string {
	return err.Val
}

// NotFoundError represents when a cluster cannot be found
type NotFoundError struct {
	Entity string
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", err.Entity)
}

// RequiredArgumentError represents an invalid argument passed to a command
type RequiredArgumentError struct {
	Arg string
}

func (err RequiredArgumentError) Error() string {
	return fmt.Sprintf("%s is required", err.Arg)
}
