package gizmo

import "errors"

// Content is the most basic structure in the library. It represents an
// arbitrary piece of data that can be represented in different Views and is
// automatically versioned on each save.
type Content interface {
	// Identifier is the unique ID of the Content object across all Contexts.
	Identifier() uint

	// SetIdentifier sets the unique ID for the Content object.
	SetIdentifier(id uint) error

	// CommitID is the ID of the commit with the specific data in this object.
	CommitID() uint

	// SetCommitID sets the commit ID for the Content object.
	SetCommitID(id uint) error

	// ViewID is the ID of the Content that this exists in.
	ViewID() uint

	// SetViewID sets the ViewID for the Content object.
	SetViewID(id uint) error

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

// ContentObject is the default implementation of the Content interface.
// Its general purpose is to be embedded in the other content objects.
type ContentObject struct {
	id         uint
	commitID   uint
	viewID     uint
	attributes map[string]interface{}
}

// Identifier is the unique ID of the Content object across all Contexts.
func (c *ContentObject) Identifier() uint {
	return c.id
}

// SetIdentifier sets the unique ID for the Content object.
func (c *ContentObject) SetIdentifier(id uint) error {
	if id == 0 {
		return errors.New("Identifier must be greater than 0")
	}

	c.id = id
	return nil
}

// CommitID is the ID of the commit with the specific data in this object.
func (c *ContentObject) CommitID() uint {
	return c.commitID
}

// SetCommitID sets the commit ID for the Content object.
func (c *ContentObject) SetCommitID(commitID uint) error {
	if commitID == 0 {
		return errors.New("CommitID must be greater than 0")
	}

	c.commitID = commitID
	return nil
}

// ViewID is the ID of the Content that this exists in.
func (c *ContentObject) ViewID() uint {
	return c.viewID
}

// SetViewID sets the ViewID for the Content object.
func (c *ContentObject) SetViewID(viewID uint) error {
	if viewID == 0 {
		return errors.New("ViewID must be greater than 0")
	}

	c.viewID = viewID
	return nil
}

// Attribute gets the value of a custom attribute. If the attribute is not
// found, nil is returned.
func (c *ContentObject) Attribute(key string) (interface{}, error) {
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
func (c *ContentObject) SetAttribute(key string, value interface{}) error {
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
func (c *ContentObject) RemoveAttribute(key string) error {
	if key == "" {
		return errors.New("Attribute key must be non-empty")
	}

	if c.attributes == nil {
		return nil
	}

	delete(c.attributes, key)
	return nil
}
