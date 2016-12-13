package lib

import "time"

// ObjectHead is a pointer to the most up-to-date commit of a given object.
type ObjectHead struct {
	ID         uint
	ContextID  uint
	CommitID   uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ArchivedAt time.Time
}
