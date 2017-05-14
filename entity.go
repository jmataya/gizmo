package gizmo

import (
	"errors"
	"fmt"
)

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

	// Kind is an identifier of the type of Entity.
	Kind() string

	// SetKind sets the type of Entity.
	SetKind(str string) error

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

	// Relations gets all of the relations to other Entities.
	Relations() map[string][]int64

	// RelationsByEntity gets all of the associations between this Entity and an
	// Entity of a specified type.
	RelationsByEntity(entityType string) ([]int64, error)

	// SetRelation creates a mapping between this Entity and another existing
	// Entity. If the mapping already exists, nothing is changed.
	SetRelation(entityType string, entityID int64) error

	// UpdateRelation updates the ID of a mapping between this Entity and another
	// previously associated Entity.
	UpdateRelation(entityType string, oldID int64, newID int64) error

	// RemoveRelation deletes a mapping between this Entity and another Entity.
	// If the mapping did not previously exist, an error is thrown.
	RemoveRelation(entityType string, entityID int64) error
}

// EntityObject is the default implementation of the Entity interface.
// Its general purpose is to be embedded in the other Entity objects.
type EntityObject struct {
	id         int64
	commitID   int64
	viewID     int64
	kind       string
	attributes map[string]interface{}
	relations  map[string][]int64
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

// Kind is an identifier of the type of Entity.
func (c *EntityObject) Kind() string {
	return c.kind
}

// SetKind sets the type of Entity.
func (c *EntityObject) SetKind(kind string) error {
	if kind == "" {
		return errors.New("Kind must be non-empty")
	}

	c.kind = kind
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

// Relations gets all of the relations to other Entities.
func (c *EntityObject) Relations() map[string][]int64 {
	if c.relations == nil {
		return map[string][]int64{}
	}

	return c.relations
}

// RelationsByEntity gets all of the associations between this Entity and an
// Entity of a specified type.
func (c *EntityObject) RelationsByEntity(entityType string) ([]int64, error) {
	if entityType == "" {
		return nil, errors.New("Entity type must be non-empty")
	}

	if c.relations == nil {
		return []int64{}, nil
	}

	ids, ok := c.relations[entityType]
	if !ok {
		return []int64{}, nil
	}

	return ids, nil
}

// SetRelation creates a mapping between this Entity and another existing
// Entity. If the mapping already exists, nothign is changed.
func (c *EntityObject) SetRelation(entityType string, entityID int64) error {
	if entityType == "" {
		return errors.New("Entity type must be non-empty")
	}

	if c.relations == nil {
		c.relations = map[string][]int64{}
	}

	ids, ok := c.relations[entityType]
	if !ok {
		ids = []int64{}
	}

	for _, id := range ids {
		if id == entityID {
			// The ID is already in the list, meaning the mapping exists.
			// Return with no error and no change.
			return nil
		}
	}

	c.relations[entityType] = append(ids, entityID)
	return nil
}

// UpdateRelation updates the ID of a mapping between this Entity and another
// previously associated Entity.
func (c *EntityObject) UpdateRelation(entityType string, oldID int64, newID int64) error {
	if entityType == "" {
		return errors.New("Entity type must be non-empty")
	} else if c.relations == nil {
		return fmt.Errorf("Mapping to %d of type %s is not found", oldID, entityType)
	}

	ids, ok := c.relations[entityType]
	if !ok {
		return fmt.Errorf("Mapping to %d of type %s is not found", oldID, entityType)
	}

	for index, id := range ids {
		if id == oldID {
			ids[index] = newID
			c.relations[entityType] = ids
			return nil
		}
	}

	return fmt.Errorf("Mapping to %d of type %s is not found", oldID, entityType)
}

// RemoveRelation deletes a mapping between this Entity and another Entity.
// If the mapping did not previously exist, an error is thrown.
func (c *EntityObject) RemoveRelation(entityType string, entityID int64) error {
	if entityType == "" {
		return errors.New("Entity type must be non-empty")
	} else if c.relations == nil {
		return fmt.Errorf("Mapping to %d of type %s is not found", entityID, entityType)
	}

	ids, ok := c.relations[entityType]
	if !ok {
		return fmt.Errorf("Mapping to %d of type %s is not found", entityID, entityType)
	}

	for index, id := range ids {
		if id == entityID {
			copy(ids[index:], ids[index+1:])
			ids = ids[:len(ids)-1]

			c.relations[entityType] = ids
			return nil
		}
	}

	return fmt.Errorf("Mapping to %d of type %s is not found", entityID, entityType)
}
