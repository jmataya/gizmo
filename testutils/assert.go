package testutils

import "testing"

const (
	defaultEqualError = "Assertion failed: expected=%v, actual=%v"
)

// Assert is a convenience structure to easy testing correctness.
type Assert struct {
	t *testing.T
}

// NewAssert creates a new assertion object.
func NewAssert(t *testing.T) Assert {
	return Assert{t}
}

// Equal checks for equality between two values.
func (a Assert) Equal(expected interface{}, actual interface{}) bool {
	return a.Equalf(expected, actual, defaultEqualError, expected, actual)
}

// Equalf checks for equality between two values and lets the callers set a
// descriptive message upon error.
func (a Assert) Equalf(expected interface{}, actual interface{}, format string, args ...interface{}) bool {
	if expected != actual {
		a.t.Errorf(format, args...)
		return false
	}

	return true
}
