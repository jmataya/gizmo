package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/FoxComm/gizmo/common"
)

const (
	sqlInsertEntityVersion = "INSERT INTO entity_versions (content_commit_id) VALUES ($1) RETURNING *"
)

// EntityVersion is a snapshot in time of the full structure of an Entity. It
// contains all the references to the Content and any other dependent Entities.
// Once inserted into the database it is completely immutable.
type EntityVersion struct {
	ID              int64
	ParentID        sql.NullInt64
	ContentCommitID int64
	CreatedAt       time.Time
}

// Validate checks all the properties on the EntityVersion and determines if
// they are all in a valid state.
func (version EntityVersion) Validate() error {
	if version.ContentCommitID == 0 {
		return fmt.Errorf(errFieldMustBeNonEmpty, "ContentCommitID")
	}

	return nil
}

// Insert adds the EntityVersion to the database and returns a copy of the
// EntityVersion with the values that were inserted.
func (version EntityVersion) Insert(db common.DB) (EntityVersion, error) {
	var newVersion EntityVersion

	if err := version.Validate(); err != nil {
		return newVersion, err
	}

	if version.ID != 0 {
		return newVersion, fmt.Errorf(errNoInsertHasPrimaryKey, "EntityVersion")
	}

	stmt, err := db.Prepare(sqlInsertEntityVersion)
	if err != nil {
		return newVersion, err
	}

	var id int64
	var parentID sql.NullInt64
	var contentCommitID int64
	var createdAt time.Time

	row := stmt.QueryRow(version.ContentCommitID)
	if err := row.Scan(&id, &parentID, &contentCommitID, &createdAt); err != nil {
		return newVersion, err
	}

	newVersion.ID = id
	newVersion.ParentID = parentID
	newVersion.ContentCommitID = contentCommitID
	newVersion.CreatedAt = createdAt

	return newVersion, nil
}