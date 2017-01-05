package models

import (
	"bytes"
	"crypto/sha1"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	sqlInsertObjectForm = "INSERT INTO object_forms (kind, attributes) VALUES ($1, $2) RETURNING *"
)

// ObjectForm is the central component in the object model. It is a flat
// collection of attributes. The key of each attribute is a hash of the
// attribute's value.
type ObjectForm struct {
	ID         uint
	Kind       string
	Attributes ObjectFormAttributes
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewObjectForm generates a new ObjectForm.
func NewObjectForm(kind string) *ObjectForm {
	return &ObjectForm{
		Kind:       kind,
		Attributes: map[string]interface{}{},
	}
}

// AddAttribute computes a SHA for the value and adds it to the ObjectForm.
// The SHA is returned as a result.
func (form *ObjectForm) AddAttribute(value interface{}) (string, error) {
	hasher := sha1.New()
	valBytes := new(bytes.Buffer)

	if err := json.NewEncoder(valBytes).Encode(value); err != nil {
		return "", err
	}

	checksum := hasher.Sum(valBytes.Bytes())
	hash := fmt.Sprintf("%x", checksum)

	form.Attributes[hash] = value
	return hash, nil
}

// Validate checks the properties on the ObjectForm and determines if they are
// all in a valid state.
func (form ObjectForm) Validate() error {
	if form.Kind == "" {
		return errors.New(errObjectFormMustHaveKind)
	}

	return nil
}

// Insert adds the ObjectForm to the database and returns a copy of the
// ObjectForm with values that were inserted.
func (form ObjectForm) Insert(db *sql.DB) (ObjectForm, error) {
	var newForm ObjectForm

	if err := form.Validate(); err != nil {
		return newForm, err
	}

	if form.ID != 0 {
		return newForm, fmt.Errorf(errNoInsertHasPrimaryKey, "ObjectForm")
	}

	stmt, err := db.Prepare(sqlInsertObjectForm)
	if err != nil {
		return newForm, err
	}

	var id uint
	var kind string
	var attributes ObjectFormAttributes
	var createdAt time.Time
	var updatedAt time.Time

	row := stmt.QueryRow(form.Kind, form.Attributes)
	if err := row.Scan(&id, &kind, &attributes, &createdAt, &updatedAt); err != nil {
		return newForm, err
	}

	return ObjectForm{
		ID:         id,
		Kind:       kind,
		Attributes: attributes,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}
