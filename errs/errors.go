package errs

import (
	"errors"
	"fmt"
)

var (
	ErrInvalid  = errors.New("invalid")
	ErrNil      = errors.New("cannot be nil")
	ErrFailedTo = errors.New("failed to")
)

// Invalid returns an error indicating that value is not valid for the given label.
// The error wraps ErrInvalid, allowing it to be checked with errors.Is.
func Invalid(label string) error {
	return fmt.Errorf("%w %s", ErrInvalid, label)
}

// Nil returns an error indicating that the value identified by label is null.
// The error wraps ErrNil, allowing it to be checked with errors.Is.
func Nil(label string) error {
	return fmt.Errorf("%s %w", label, ErrNil)
}

// Failed returns a generic failure error associated with the given label.
// The error wraps ErrFailedTo, allowing it to be checked with errors.Is.
func FailedTo(label string) error {
	return fmt.Errorf("%w %s", ErrFailedTo, label)
}

// WithContext returns a new error that adds context to err.
// The original error is wrapped using %w, so it can be
// inspected with errors.Is and errors.As.
func WithContext(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
