package models

const (
	// General error messages.
	errNoInsertHasPrimaryKey = "%s has a primary key and cannot be inserted"

	// Query error messages.
	errFieldMustBeGreaterThanZero = "%s must be greater than zero"
	errFieldMustBeZero            = "%s must be zero"
	errFieldMustBeNonEmpty        = "%s must be non-empty"
)
