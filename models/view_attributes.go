package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// ViewAttributes is a single-level map containing all the values that
// make up the View.
type ViewAttributes map[string]interface{}

// Scan is an interface used for getting JSON out of the database
// and it turns it into a struct.
func (viewA *ViewAttributes) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	return json.Unmarshal(source, viewA)
}

// Value converts the struct to a byte array of JSON.
func (viewA ViewAttributes) Value() (driver.Value, error) {
	j, err := json.Marshal(viewA)
	return j, err
}
