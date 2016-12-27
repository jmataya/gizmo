package models

// SimpleObject is the interface that a model must implement to be used as part
// of the illuminated services.
type SimpleObject interface {
	// Identifier accesses the object's primary key.
	Identifier() uint

	// SetIdentifier sets the object's primary key.
	SetIdentifier(id uint)
}
