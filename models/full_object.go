package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	sqlSelectFullObjectLatest = `
    SELECT f.*, s.*, c.* FROM object_heads AS h
    INNER JOIN object_commits AS c ON h.commit_id = c.id
    INNER JOIN object_forms AS f ON c.form_id = f.id
    INNER JOIN object_shadows AS s ON c.shadow_id = s.id
    WHERE c.form_id = $1 AND h.context_id = $2
  `

	sqlSelectFullObjectByCommit = `
		SELECT f.*, s.*, c.* FROM object_commits AS c
		INNER JOIN object_forms AS f ON c.form_id = f.id
		INNER JOIN object_shadows AS s ON c.shadow_id = s.id
		WHERE c.id = $1
	`
)

// FullObject is a convenience structure that aggregrates the main components
// of a Content object: form, shadow, and commit.
type FullObject struct {
	Form   ObjectForm
	Shadow ObjectShadow
	Commit ObjectCommit
}

// FindLatest retrieves the most recent FullObject in a given view.
func (f FullObject) FindLatest(db *sql.DB, id uint, viewID uint) (FullObject, error) {
	if id == 0 {
		return f, fmt.Errorf(errFieldMustBeGreaterThanZero, "id")
	} else if viewID == 0 {
		return f, fmt.Errorf(errFieldMustBeGreaterThanZero, "viewID")
	}

	row := db.QueryRow(sqlSelectFullObjectLatest, id, viewID)
	return f.findRow(row)
}

// FindByCommit retrieves a FullObject at a specific commit.
func (f FullObject) FindByCommit(db *sql.DB, commitID uint) (FullObject, error) {
	if commitID == 0 {
		return f, fmt.Errorf(errFieldMustBeGreaterThanZero, "commitID")
	}

	row := db.QueryRow(sqlSelectFullObjectByCommit, commitID)
	return f.findRow(row)
}

func (f FullObject) findRow(row *sql.Row) (FullObject, error) {
	var found FullObject

	var formID uint
	var formKind string
	var formAttributes ObjectFormAttributes
	var formCreatedAt time.Time
	var formUpdatedAt time.Time
	var shadowID uint
	var shadowFormID uint
	var shadowAttributes ObjectShadowAttributes
	var shadowCreatedAt time.Time
	var commitID uint
	var commitFormID uint
	var commitShadowID uint
	var commitPreviousID sql.NullInt64
	var commitCreatedAt time.Time

	err := row.Scan(&formID, &formKind, &formAttributes, &formCreatedAt, &formUpdatedAt,
		&shadowID, &shadowFormID, &shadowAttributes, &shadowCreatedAt,
		&commitID, &commitFormID, &commitShadowID, &commitPreviousID, &commitCreatedAt)

	if err != nil {
		return found, err
	}

	found.Form.ID = formID
	found.Form.Kind = formKind
	found.Form.Attributes = formAttributes
	found.Form.CreatedAt = formCreatedAt
	found.Form.UpdatedAt = formUpdatedAt
	found.Shadow.ID = shadowID
	found.Shadow.FormID = shadowFormID
	found.Shadow.Attributes = shadowAttributes
	found.Shadow.CreatedAt = shadowCreatedAt
	found.Commit.ID = commitID
	found.Commit.FormID = commitFormID
	found.Commit.ShadowID = commitShadowID
	found.Commit.PreviousID = commitPreviousID
	found.Commit.CreatedAt = commitCreatedAt

	return found, nil
}
