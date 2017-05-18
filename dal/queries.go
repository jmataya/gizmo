package dal

import (
	"fmt"
	"strings"

	"github.com/FoxComm/gizmo/common"
	"github.com/FoxComm/gizmo/models"
)

const (
	sqlInsertEntityRoot = "INSERT INTO entity_roots (kind) VALUES ($1) RETURNING *"

	// General error messages.
	errNoInsertHasPrimaryKey = "%s has a primary key and cannot be inserted"
)

func InsertEntityHaed(db common.DB, model models.EntityHead) (models.EntityHead, error) {
	if err := model.Validate(); err != nil {
		return model, err
	}

	if model.ID != 0 {
		return model, fmt.Errorf(errNoInsertHasPrimaryKey, "EntityHead")
	}
}

// InsertEntityRoot inserts a new EntityRoot object into the database.
func InsertEntityRoot(db common.DB, model models.EntityRoot) (models.EntityRoot, error) {
	if err := model.Validate(); err != nil {
		return model, err
	}

	if model.ID != 0 {
		return model, fmt.Errorf(errNoInsertHasPrimaryKey, "EntityRoot")
	}

	var root models.EntityRoot

	d := NewDataAccessLayer(db)
	row := d.Query(sqlInsertEntityRoot, strings.ToLower(model.Kind))
	d.Scan(row, &root.ID, &root.Kind, &root.CreatedAt, &root.ArchivedAt)

	return root, d.Result()
}
