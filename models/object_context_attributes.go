package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// ObjectContextAttributes is a single-level map containing all the values that
// make up the ObjectContext.
type ObjectContextAttributes map[string]interface{}

// Scan is an interface used for getting JSON out of the database
// and it turns it into a struct.
func (oca *ObjectContextAttributes) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	return json.Unmarshal(source, oca)
}

// Value converts the struct to a byte array of JSON.
func (oca ObjectContextAttributes) Value() (driver.Value, error) {
	j, err := json.Marshal(oca)
	return j, err
}
