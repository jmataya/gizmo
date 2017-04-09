package testutils

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	defaultDB   = "gizmo_test"
	defaultUser = "gizmo"
)

// InitDB initializes a new connection to the default testing database.
func InitDB(t *testing.T) *sql.DB {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = defaultDB
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = defaultUser
	}

	dsn := fmt.Sprintf("user=%s dbname=%s sslmode=disable", user, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Error(err)
	}

	return db
}
