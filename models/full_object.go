package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/FoxComm/gizmo/common"
	log "github.com/sirupsen/logrus"
)

const (
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

// Find retrieves a FullObject at a specific commit.
func (f FullObject) Find(db *sql.DB, commitID int64) (FullObject, error) {
	if commitID == 0 {
		return f, fmt.Errorf(errFieldMustBeGreaterThanZero, "commitID")
	}

	row := db.QueryRow(sqlSelectFullObjectByCommit, commitID)
	return f.findRow(row)
}

// Insert adds the FullObject to the database.
func (f FullObject) Insert(db common.DB) (FullObject, error) {
	if f.Form.ID != 0 {
		return f, fmt.Errorf(errFieldMustBeZero, "Form.ID")
	} else if f.Shadow.ID != 0 {
		return f, fmt.Errorf(errFieldMustBeZero, "Shadow.ID")
	} else if f.Commit.ID != 0 {
		return f, fmt.Errorf(errFieldMustBeZero, "Shadow.ID")
	}

	log.Debugln("Inserting Form")
	newForm, err := f.Form.Insert(db)
	if err != nil {
		return f, err
	}
	log.Debugf("Inserted Form with ID=%d", newForm.ID)

	log.Debugln("Inserting Shadow")
	f.Shadow.FormID = newForm.ID
	newShadow, err := f.Shadow.Insert(db)
	if err != nil {
		return f, err
	}
	log.Debugf("Inserted Shadow with ID=%d", newShadow.ID)

	log.Debugln("Inserting Commit")
	f.Commit.FormID = newForm.ID
	f.Commit.ShadowID = newShadow.ID
	newCommit, err := f.Commit.Insert(db)
	if err != nil {
		return f, err
	}
	log.Debugf("Inserted Commit with ID=%d", newCommit.ID)

	return FullObject{
		Form:   newForm,
		Shadow: newShadow,
		Commit: newCommit,
	}, nil
}

func (f FullObject) findRow(row *sql.Row) (FullObject, error) {
	var found FullObject

	var formID int64
	var formKind string
	var formAttributes ObjectFormAttributes
	var formCreatedAt time.Time
	var formUpdatedAt time.Time
	var shadowID int64
	var shadowFormID int64
	var shadowAttributes ObjectShadowAttributes
	var shadowCreatedAt time.Time
	var commitID int64
	var commitFormID int64
	var commitShadowID int64
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
