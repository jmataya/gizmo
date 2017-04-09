package gizmo

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq" // Needed to allow database/sql to use Postgres.
)

// ContentManager is the interface for creating, managing, and deleting Content.
type ContentManager interface {
	// Find retrieves the most recent version of a Content object within a View.
	// None of the parameters are modified, including the type hint.
	Find(id uint, viewID uint, typeHint Content) (Content, error)

	// FindByCommit retrieves a Content object at a specific commit. This will
	// retrieve the entire object, including all associated objects, as of that
	// commit. None of the parameters are modified, including the type hint.
	FindByCommit(commitID uint, typeHint Content) (Content, error)

	// Create saves a new Content object as a new entity and returns the created
	// version of the object back. If the ID, Content ID, or Commit ID of the
	// Content object have previously been set, they will be ignored.
	Create(toCreate Content, viewID uint) (Content, error)

	// Update modifies a previously saved Content object. The new version will be
	// branched from the Content's Commit ID, and will use the ID and View ID to
	// save in the appropriate format. If the object has not previously been saved
	// the method will error.
	Update(toUpdate Content) (Content, error)

	// Delete performs a soft-delete on a Content object. This must occur at the
	// most recent commit, so the Content is identified by the ID and View ID.
	Delete(id uint, viewID uint) error
}

// NewContentManager connects a PostgreSQL database with the supplied connection
// parameters and returns the created ContentManager.
func NewContentManager(host, dbName, user, password string, sslMode bool) (ContentManager, error) {
	if host == "" {
		return nil, errors.New("Database host must be non-empty")
	} else if dbName == "" {
		return nil, errors.New("Database name must be non-empty")
	} else if user == "" {
		return nil, errors.New("Database user must be non-empty")
	}

	var sslModeStr string
	if sslMode {
		sslModeStr = "enable"
	} else {
		sslModeStr = "disable"
	}

	dsn := fmt.Sprintf("host=%s dbname=%s sslMode=%s user=%s", host, dbName, sslModeStr, user)
	if password != "" {
		dsn = fmt.Sprintf("%s password=%s", dsn, password)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &defaultContentManager{db: db}, nil
}

type defaultContentManager struct {
	db *sql.DB
}

func (d *defaultContentManager) Find(id uint, viewID uint, typeHint Content) (Content, error) {
	return nil, errors.New("Not implemented")
}

func (d *defaultContentManager) FindByCommit(commitID uint, typeHint Content) (Content, error) {
	return nil, errors.New("Not implemented")
}

func (d *defaultContentManager) Create(toCreate Content, viewID uint) (Content, error) {
	return nil, errors.New("Not implemented")
}

func (d *defaultContentManager) Update(toUpdate Content) (Content, error) {
	return nil, errors.New("Not implemented")
}

func (d *defaultContentManager) Delete(id uint, viewID uint) error {
	return errors.New("Not implemented")
}
