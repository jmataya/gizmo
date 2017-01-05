package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	sqlInsertObjectCommit = "INSERT INTO object_commits (form_id, shadow_id) VALUES ($1, $2) RETURNING id, form_id, shadow_id, created_at"
)

// ObjectCommit represents an update to an object. It is an immutable object in
// the database and contains a reference to the commit that came before it.
type ObjectCommit struct {
	ID         uint
	FormID     uint
	ShadowID   uint
	PreviousID uint
	CreatedAt  time.Time
}

// Validate checks the properties on the ObjectCommit and determines if they
// are all in a valid state.
func (commit ObjectCommit) Validate() error {
	if commit.FormID == 0 {
		return errors.New(errObjectCommitMustHaveFormID)
	} else if commit.ShadowID == 0 {
		return errors.New(errObjectCommitMustHaveShadowID)
	}

	return nil
}

// Insert adds the ObjectCommit to the database and returns a copy of the
// ObjectCommits with values that were inserted.
func (commit ObjectCommit) Insert(db *sql.DB) (ObjectCommit, error) {
	var newCommit ObjectCommit

	if err := commit.Validate(); err != nil {
		return newCommit, err
	}

	if commit.ID != 0 {
		return newCommit, fmt.Errorf(errNoInsertHasPrimaryKey, "ObjectCommit")
	}

	stmt, err := db.Prepare(sqlInsertObjectCommit)
	if err != nil {
		return newCommit, err
	}

	var id uint
	var formID uint
	var shadowID uint
	var createdAt time.Time

	row := stmt.QueryRow(commit.FormID, commit.ShadowID)
	if err := row.Scan(&id, &formID, &shadowID, &createdAt); err != nil {
		return newCommit, err
	}

	return ObjectCommit{
		ID:        id,
		FormID:    formID,
		ShadowID:  shadowID,
		CreatedAt: createdAt,
	}, nil
}
