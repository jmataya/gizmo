package lib

import "time"

// ObjectForm is the central component in the object model. It is a flat
// collection of attributes. The key of each attribute is a hash of the
// attribute's value.
type ObjectForm struct {
	ID         uint
	Kind       string
	Attributes map[string]interface{}
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
