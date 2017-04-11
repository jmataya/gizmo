package gizmo

import "errors"

// Entity is the most basic structure in the library. It represents an
// arbitrary piece of data that can be represented in different Views and is
// automatically versioned on each save.
type Entity interface {
	// Identifier is the unique ID of the Entity object across all Views.
	Identifier() int64

	// SetIdentifier sets the unique ID for the Entity object.
	SetIdentifier(id int64) error

	// CommitID is the ID of the commit with the specific data in this object.
	CommitID() int64

	// SetCommitID sets the commit ID for the Entity object.
	SetCommitID(id int64) error

	// ViewID is the ID of the Entity that this exists in.
	ViewID() int64

	// SetViewID sets the ViewID for the Entity object.
	SetViewID(id int64) error

	// Attributes gets all of the custom attributes.
	Attributes() map[string]interface{}

	// Attribute gets the value of a custom attribute. If the attribute is not
	// found, nil is returned.
	Attribute(key string) (interface{}, error)

	// SetAttribute sets the value of a custom attribute. It can be used to
	// either create or update the attribute.
	SetAttribute(key string, value interface{}) error

	// RemoveAttribute removes the attribute from the list of custom attributes.
	// If the key does not exist, the function is a no-op.
	RemoveAttribute(key string) error
}

// EntityObject is the default implementation of the Entity interface.
// Its general purpose is to be embedded in the other Entity objects.
type EntityObject struct {
	id         int64
	commitID   int64
	viewID     int64
	attributes map[string]interface{}
}

// Identifier is the unique ID of the Entity object across all Views.
func (c *EntityObject) Identifier() int64 {
	return c.id
}

// SetIdentifier sets the unique ID for the Entity object.
func (c *EntityObject) SetIdentifier(id int64) error {
	if id == 0 {
		return errors.New("Identifier must be greater than 0")
	}

	c.id = id
	return nil
}

// CommitID is the ID of the commit with the specific data in this object.
func (c *EntityObject) CommitID() int64 {
	return c.commitID
}

// SetCommitID sets the commit ID for the Entity object.
func (c *EntityObject) SetCommitID(commitID int64) error {
	if commitID == 0 {
		return errors.New("CommitID must be greater than 0")
	}

	c.commitID = commitID
	return nil
}

// ViewID is the ID of the Entity that this exists in.
func (c *EntityObject) ViewID() int64 {
	return c.viewID
}

// SetViewID sets the ViewID for the Entity object.
func (c *EntityObject) SetViewID(viewID int64) error {
	if viewID == 0 {
		return errors.New("ViewID must be greater than 0")
	}

	c.viewID = viewID
	return nil
}

// Attributes gets the set of custom attributes.
func (c *EntityObject) Attributes() map[string]interface{} {
	return c.attributes
}

// Attribute gets the value of a custom attribute. If the attribute is not
// found, nil is returned.
func (c *EntityObject) Attribute(key string) (interface{}, error) {
	if key == "" {
		return nil, errors.New("Attribute key must be non-empty")
	}

	if c.attributes == nil {
		return nil, nil
	}

	value, found := c.attributes[key]
	if !found {
		return nil, nil
	}

	return value, nil
}

// SetAttribute sets the value of a custom attribute. It can be used to
// either create or update the attribute.
func (c *EntityObject) SetAttribute(key string, value interface{}) error {
	if key == "" {
		return errors.New("Attribute key must be non-empty")
	}

	if c.attributes == nil {
		c.attributes = map[string]interface{}{}
	}

	c.attributes[key] = value
	return nil
}

// RemoveAttribute removes the attribute from the list of custom attributes.
// If the key does not exist, the function is a no-op.
func (c *EntityObject) RemoveAttribute(key string) error {
	if key == "" {
		return errors.New("Attribute key must be non-empty")
	}

	if c.attributes == nil {
		return nil
	}

	delete(c.attributes, key)
	return nil
}
