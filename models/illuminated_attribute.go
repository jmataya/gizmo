package models

// IlluminatedAttributes is the type and value of a single piece of data
// in an IlluminatedObject.
type IlluminatedAttribute struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
