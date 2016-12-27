package models

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

// NewObjectShadow generates a new ObjectShadow.
func NewObjectShadow() *ObjectShadow {
	return &ObjectShadow{
		Attributes: map[string]attribute{},
	}
}

// AddAttribute adds an attribute to the ObjectShadow.
func (shadow *ObjectShadow) AddAttribute(attrName, attrType, attrRef string) {
	attr := attribute{Type: attrType, Ref: attrRef}
	shadow.Attributes[attrName] = attr
}

type attribute struct {
	Type string
	Ref  string
}
