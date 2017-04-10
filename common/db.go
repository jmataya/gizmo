package common

import "database/sql"

type DB interface {
	Prepare(query string) (*sql.Stmt, error)
}
