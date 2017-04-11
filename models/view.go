package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	sqlInsertView = "INSERT INTO views (name, attributes) VALUES ($1, $2) RETURNING *"
)

// View is an object that is used to define the different ways that an Entity,
// Content, or Taxonomy can render.
type View struct {
	ID         int64
	Name       string
	Attributes ViewAttributes
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Validate checks the properties on the View and determines if they
// are in a valid state.
func (view View) Validate() error {
	if view.Name == "" {
		return fmt.Errorf(errFieldMustBeNonEmpty, "Name")
	}

	return nil
}

// Insert adds the View to the database and returns a copy of the
// View with values that were inserted.
func (view View) Insert(db *sql.DB) (View, error) {
	if err := view.Validate(); err != nil {
		return view, err
	}

	if view.ID != 0 {
		return view, fmt.Errorf(errNoInsertHasPrimaryKey, "View")
	}

	stmt, err := db.Prepare(sqlInsertView)
	if err != nil {
		return view, err
	}

	var id int64
	var name string
	var attributes ViewAttributes
	var createdAt time.Time
	var updatedAt time.Time

	row := stmt.QueryRow(view.Name, view.Attributes)
	if err := row.Scan(&id, &name, &attributes, &createdAt, &updatedAt); err != nil {
		return view, err
	}

	return View{
		ID:         id,
		Name:       name,
		Attributes: attributes,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}
