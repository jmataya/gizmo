package models

// IlluminatedObject is the representation of an object when its form and
// shadow have been stitched together.
type IlluminatedObject struct {
	ContextID  uint                            `json:"contextId"`
	FormID     uint                            `json:"formId"`
	Kind       string                          `json:"kind"`
	Attributes map[string]IlluminatedAttribute `json:"attributes"`
}
