package lib

import "time"

// ObjectContext is an object that is used to define which ObjectShadow to use
// for an object.
type ObjectContext struct {
	ID         uint
	Name       string
	Attributes map[string]interface{}
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
