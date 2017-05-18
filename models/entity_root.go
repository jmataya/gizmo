package models

import (
	"fmt"
	"time"
)

type EntityRoot struct {
	ID         int64
	Kind       string
	CreatedAt  time.Time
	ArchivedAt *time.Time
}

func (root EntityRoot) Validate() error {
	if root.Kind == "" {
		return fmt.Errorf(errFieldMustBeNonEmpty, "Kind")
	}
	return nil
}
