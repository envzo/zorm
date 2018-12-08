package db

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	openOnce sync.Once
	db       *sql.DB
)

func Open(t, dsn string) error {
	var err error
	openOnce.Do(func() {
		db, err = sql.Open(t, dsn)
	})
	if err != nil {
		return err
	}
	return nil
}

func DB() *sql.DB { return db }
