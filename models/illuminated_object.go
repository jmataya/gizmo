package lib

// IlluminatedObject is the representation of an object when its form and
// shadow have been stitched together.
type IlluminatedObject struct {
	ContextID  uint
	FormID     uint
	Attributes map[string]illuminatedAttribute
}

type illuminatedAttribute struct {
	Type  string
	Value interface{}
}
