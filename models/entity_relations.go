package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// EntityRelations is a map graphs the linkage of an entity (the EntityVersion
// that includes this value) to the commit IDs of other entity versions. The
// string in the map represents the version kind of the destination of the
// mapping.
type EntityRelations map[string][]int64

// NewEntityRelations initializes a new EntityRelations object based on a map
// of relations.
func NewEntityRelations(relations map[string][]int64) EntityRelations {
	entityRelations := EntityRelations(relations)
	return entityRelations
}

// Scan is an interface for getting JSON out of the database and turns it into a
// struct.
func (er *EntityRelations) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	return json.Unmarshal(source, er)
}

// Value converts the struct to a byte array of JSON.
func (er *EntityRelations) Value() (driver.Value, error) {
	return json.Marshal(er)
}
