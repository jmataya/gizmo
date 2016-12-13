package lib

import "time"

// ObjectCommit represents an update to an object. It is an immutable object in
// the database and contains a reference to the commit that came before it.
type ObjectCommit struct {
	ID         uint
	FormID     uint
	ShadowID   uint
	PreviousID uint
	CreatedAt  time.Time
}
