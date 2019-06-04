package db

import "database/sql"

type T struct {
	tx *sql.Tx
}

var t = T{}

func Tx() *T {
	return &t
}

func (t *T) Begin(query string, args ...interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	t.tx = tx

	return nil
}

func (t *T) Exec(query string, args ...interface{}) (sql.Result, error) {
	return t.tx.Exec(query, args...)
}

func (t *T) Rollback() error {
	return t.tx.Rollback()
}

func (t *T) Commit() error {
	return t.tx.Commit()
}
