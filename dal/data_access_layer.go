package dal

import (
	"database/sql"
	"fmt"

	"github.com/jmataya/gizmo/common"
)

type DataAccessLayer struct {
	db  common.DB
	err error
}

func NewDataAccessLayer(db common.DB) *DataAccessLayer {
	return &DataAccessLayer{db: db}
}

func (d *DataAccessLayer) ValidateInsert(model Model) {
	if err := model.Validate(); err != n
	d.err = model.Validate()
}

func (d *DataAccessLayer) Query(query string, args ...interface{}) *sql.Row {
	stmt := d.prepare(query)
	return d.stmtQueryRow(stmt, args...)
}

func (d *DataAccessLayer) prepare(query string) *sql.Stmt {
	if d.err != nil {
		return nil
	}

	stmt, err := d.db.Prepare(query)
	if err != nil {
		d.err = fmt.Errorf("Unable to prepare statement %s with error %s", query, err.Error())
		return nil
	}

	return stmt
}

func (d *DataAccessLayer) stmtQueryRow(stmt *sql.Stmt, args ...interface{}) *sql.Row {
	if d.err != nil {
		return nil
	}

	return stmt.QueryRow(args...)
}

func (d *DataAccessLayer) Scan(row *sql.Row, dest ...interface{}) {
	if d.err != nil {
		return
	}

	d.err = row.Scan(dest...)
}

func (d *DataAccessLayer) Result() error {
	return d.err
}
