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

func Open(t, dsn string, openConn, idleConn int) error {
	var err error
	openOnce.Do(func() {
		db, err = sql.Open(t, dsn)
		db.SetMaxOpenConns(openConn)
		db.SetMaxIdleConns(idleConn)
	})
	if err != nil {
		return err
	}
	return nil
}

func DB() *sql.DB { return db }

type Rule struct {
	S string
	P interface{}
}

type Join struct {
	T    string
	Rule []Rule
}
