package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// ObjectFormAttributes is a single-level map containing all the values that
// make up the ObjectForm.
type ObjectFormAttributes map[string]interface{}

// Scan is an interface used for getting JSON out of the database
// and it turns it into a struct.
func (ofa *ObjectFormAttributes) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	return json.Unmarshal(source, ofa)
}

// Value converts the struct to a byte array of JSON.
func (ofa ObjectFormAttributes) Value() (driver.Value, error) {
	j, err := json.Marshal(ofa)
	return j, err
}
