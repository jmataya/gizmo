package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/FoxComm/gizmo/common"
)

const (
	sqlInsertObjectShadow = "INSERT INTO object_shadows (form_id, attributes) VALUES ($1, $2) RETURNING *"
)

// ObjectShadow is a view of data on a form. It is an immutable record in the
// database that defines which attributes should be visible on the illuminated
// object.
type ObjectShadow struct {
	ID         int64
	FormID     int64
	Attributes ObjectShadowAttributes
	CreatedAt  time.Time
}

// NewObjectShadow generates a new ObjectShadow.
func NewObjectShadow() *ObjectShadow {
	return &ObjectShadow{
		Attributes: map[string]attribute{},
	}
}

// AddAttribute adds an attribute to the ObjectShadow.
func (shadow *ObjectShadow) AddAttribute(attrName, attrType, attrRef string) error {
	if attrName == "" {
		return fmt.Errorf(errFieldMustBeNonEmpty, "attrName")
	} else if attrType == "" {
		return fmt.Errorf(errFieldMustBeNonEmpty, "attrType")
	} else if attrRef == "" {
		return fmt.Errorf(errFieldMustBeNonEmpty, "attrRef")
	}

	attr := attribute{Type: attrType, Ref: attrRef}
	shadow.Attributes[attrName] = attr
	return nil
}

// Validate checks the properties on the ObjectShadow and determines if they
// are all in a valid state.
func (shadow ObjectShadow) Validate() error {
	if shadow.FormID == 0 {
		return errors.New(errObjectShadowMustHaveFormID)
	}

	return nil
}

// Insert adds the ObjectShadow to the database and returns a copy of the
// ObjectShadow with the values that were inserted.
func (shadow ObjectShadow) Insert(db common.DB) (ObjectShadow, error) {
	var newShadow ObjectShadow

	if err := shadow.Validate(); err != nil {
		return newShadow, err
	}

	if shadow.ID != 0 {
		return newShadow, fmt.Errorf(errNoInsertHasPrimaryKey, "ObjectShadow")
	}

	stmt, err := db.Prepare(sqlInsertObjectShadow)
	if err != nil {
		return newShadow, err
	}

	var id int64
	var formID int64
	var attributes ObjectShadowAttributes
	var createdAt time.Time

	row := stmt.QueryRow(shadow.FormID, shadow.Attributes)
	if err := row.Scan(&id, &formID, &attributes, &createdAt); err != nil {
		return newShadow, err
	}

	return ObjectShadow{
		ID:         id,
		FormID:     formID,
		Attributes: attributes,
		CreatedAt:  createdAt,
	}, nil
}
