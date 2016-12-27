package models

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"time"
)

// ObjectForm is the central component in the object model. It is a flat
// collection of attributes. The key of each attribute is a hash of the
// attribute's value.
type ObjectForm struct {
	ID         uint
	Kind       string
	Attributes map[string]interface{}
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewObjectForm generates a new ObjectForm.
func NewObjectForm(kind string) *ObjectForm {
	return &ObjectForm{
		Kind:       kind,
		Attributes: map[string]interface{}{},
	}
}

// AddAttribute computes a SHA for the value and adds it to the ObjectForm.
// The SHA is returned as a result.
func (form *ObjectForm) AddAttribute(value interface{}) (string, error) {
	hasher := sha1.New()
	valBytes := new(bytes.Buffer)

	if err := json.NewEncoder(valBytes).Encode(value); err != nil {
		return "", err
	}

	checksum := hasher.Sum(valBytes.Bytes())
	hash := fmt.Sprintf("%x", checksum)

	form.Attributes[hash] = value
	return hash, nil
}
