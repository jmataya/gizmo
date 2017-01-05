package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// ObjectShadowAttributes is a single-level map containing all the values that
// make up the ObjectShadow.
type ObjectShadowAttributes map[string]attribute

// Scan is an interface used for getting JSON out of the database
// and it turns it into a struct.
func (osha *ObjectShadowAttributes) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	return json.Unmarshal(source, osha)
}

// Value converts the struct to a byte array of JSON.
func (osha ObjectShadowAttributes) Value() (driver.Value, error) {
	j, err := json.Marshal(osha)
	return j, err
}

type attribute struct {
	Type string `json:"type"`
	Ref  string `json:"ref"`
}
