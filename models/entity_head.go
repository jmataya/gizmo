package models

import (
	"fmt"
	"time"

	"github.com/jmataya/gizmo/common"
)

const (
	sqlInsertEntityHead = "INSERT INTO entity_heads (root_id, view_id, version_id) VALUES ($1, $2, $3) RETURNING *"
)

type EntityHead struct {
	ID        int64
	RootID    int64
	ViewID    int64
	VersionID int64

	CreatedAt  time.Time
	UpdatedAt  time.Time
	ArchivedAt *time.Time
}

func (head EntityHead) Validate() error {
	if head.ViewID == 0 {
		return fmt.Errorf(errFieldMustBeNonEmpty, "ViewID")
	} else if head.VersionID == 0 {
		return fmt.Errorf(errFieldMustBeNonEmpty, "VersionID")
	} else if head.RootID == 0 {
		return fmt.Errorf(errFieldMustBeNonEmpty, "RootID")
	}

	return nil
}

func (head EntityHead) Insert(db common.DB) (EntityHead, error) {
	if err := head.Validate(); err != nil {
		return head, err
	}

	if head.ID != 0 {
		return head, fmt.Errorf(errNoInsertHasPrimaryKey, "EntityHead")
	}

	stmt, err := db.Prepare(sqlInsertEntityHead)
	if err != nil {
		return head, err
	}

	var newHead EntityHead
	row := stmt.QueryRow(head.RootID, head.ViewID, head.VersionID)
	err = row.Scan(
		&newHead.ID,
		&newHead.RootID,
		&newHead.ViewID,
		&newHead.VersionID,
		&newHead.CreatedAt,
		&newHead.UpdatedAt,
		&newHead.ArchivedAt)

	return newHead, err
}
