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

type JoinType string

const (
	// 内联
	JoinTypeNone JoinType = ""
	// 内联
	JoinTypeInner JoinType = "inner"
	// 左联
	JoinTypeLeft JoinType = "left"
	// 右联
	JoinTypeRight JoinType = "right"
)

func Open(t, dsn string, openConn, idleConn int) error {
	var err error
	openOnce.Do(func() {
		db, err = sql.Open(t, dsn)
		// db.SetMaxOpenConns(openConn)
		// db.SetMaxIdleConns(idleConn)
		// db.SetConnMaxLifetime(10 * time.Minute)
	})
	return err
}

func DB() *sql.DB { return db }

type Rule struct {
	S string
	P interface{}
}

type Join struct {
	T     string
	Rule  []Rule
	JType JoinType
}
