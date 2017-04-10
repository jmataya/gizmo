package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	sqlInsertObjectHead = "INSERT INTO object_heads (context_id, commit_id) VALUES ($1, $2) RETURNING *"
)

// ObjectHead is a pointer to the most up-to-date commit of a given object.
type ObjectHead struct {
	ID         int64
	ContextID  int64
	CommitID   int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ArchivedAt *time.Time
}

// Validate checks the properties on the ObjectCommit and determines if they
// are all in a valid state.
func (head ObjectHead) Validate() error {
	if head.ContextID == 0 {
		return errors.New(errObjectHeadMustHaveContextID)
	} else if head.CommitID == 0 {
		return errors.New(errObjectHeadMustHaveCommitID)
	}

	return nil
}

// Insert adds the ObjectHead to the database and returns a copy of the
// ObjectHead with the values that were inserted.
func (head ObjectHead) Insert(db *sql.DB) (ObjectHead, error) {
	var newHead ObjectHead

	if err := head.Validate(); err != nil {
		return newHead, err
	}

	if head.ID != 0 {
		return newHead, fmt.Errorf(errNoInsertHasPrimaryKey, "ObjectHead")
	}

	stmt, err := db.Prepare(sqlInsertObjectHead)
	if err != nil {
		return newHead, err
	}

	var id int64
	var contextID int64
	var commitID int64
	var createdAt time.Time
	var updatedAt time.Time
	var archivedAt *time.Time

	row := stmt.QueryRow(head.ContextID, head.CommitID)
	err = row.Scan(&id, &contextID, &commitID, &createdAt, &updatedAt, &archivedAt)
	if err != nil {
		return newHead, err
	}

	return ObjectHead{
		ID:        id,
		ContextID: contextID,
		CommitID:  commitID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
