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

// FailedTo returns an error indicating a failure to perform the given action.
// It wraps ErrFailedTo and optionally wraps the provided cause.
func FailedTo(action string, err error) error {
	failedToErr := fmt.Errorf("%w %s", ErrFailedTo, action)
	if err == nil {
		return failedToErr
	}
	return fmt.Errorf("%w: %w", failedToErr, err)
}
