package lib

import "time"

// ObjectShadow is a view of data on a form. It is an immutable record in the
// database that defines which attributes should be visible on the illuminated
// object.
type ObjectShadow struct {
	ID         uint
	FormID     uint
	Attributes map[string]attribute
	CreatedAt  time.Time
}

type attribute struct {
	Type string
	Ref  string
}
