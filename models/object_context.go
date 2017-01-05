package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	sqlInsertObjectContext = "INSERT INTO object_contexts (name, attributes) VALUES ($1, $2) RETURNING *"
)

// ObjectContext is an object that is used to define which ObjectShadow to use
// for an object.
type ObjectContext struct {
	ID         uint
	Name       string
	Attributes ObjectContextAttributes
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Validate checks the properties on the ObjectContext and determines if they
// are in a valid state.
func (context ObjectContext) Validate() error {
	if context.Name == "" {
		return errors.New(errObjectContextMustHaveName)
	}

	return nil
}

// Insert adds the ObjectContext to the database and returns a copy of the
// ObjectContext with values that were inserted.
func (context ObjectContext) Insert(db *sql.DB) (ObjectContext, error) {
	var newContext ObjectContext

	if err := context.Validate(); err != nil {
		return newContext, err
	}

	if context.ID != 0 {
		return newContext, fmt.Errorf(errNoInsertHasPrimaryKey, "ObjectContext")
	}

	stmt, err := db.Prepare(sqlInsertObjectContext)
	if err != nil {
		return newContext, err
	}

	var id uint
	var name string
	var attributes ObjectContextAttributes
	var createdAt time.Time
	var updatedAt time.Time

	row := stmt.QueryRow(context.Name, context.Attributes)
	if err := row.Scan(&id, &name, &attributes, &createdAt, &updatedAt); err != nil {
		return newContext, err
	}

	return ObjectContext{
		ID:         id,
		Name:       name,
		Attributes: attributes,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}
