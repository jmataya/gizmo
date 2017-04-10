package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/FoxComm/gizmo/common"
)

const (
	sqlInsertEntityRoot = "INSERT INTO entity_roots (kind) VALUES ($1) RETURNING *"
)

type EntityRoot struct {
	ID         int64
	Kind       string
	CreatedAt  time.Time
	ArchivedAt *time.Time
}

func (root EntityRoot) Validate() error {
	if root.Kind == "" {
		return fmt.Errorf(errFieldMustBeNonEmpty, "Kind")
	}
	return nil
}

func (root EntityRoot) Insert(db common.DB) (EntityRoot, error) {
	if err := root.Validate(); err != nil {
		return root, err
	}

	if root.ID != 0 {
		return root, fmt.Errorf(errNoInsertHasPrimaryKey, "EntityRoot")
	}

	stmt, err := db.Prepare(sqlInsertEntityRoot)
	if err != nil {
		return root, err
	}

	var newRoot EntityRoot
	row := stmt.QueryRow(strings.ToLower(root.Kind))
	err = row.Scan(
		&newRoot.ID,
		&newRoot.Kind,
		&newRoot.CreatedAt,
		&newRoot.ArchivedAt)

	return newRoot, err
}
