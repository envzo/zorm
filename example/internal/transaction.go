package orm

import (
	"database/sql"
	"github.com/envzo/zorm/db"
)

type Ztx struct {
	tx *sql.Tx
}

var Zotx = &Ztx{}

func (ztx * Ztx) Begin() {
	txTest, err := db.DB().Begin()
	if err != nil {

	}
	ztx.tx = txTest
}

func (ztx * Ztx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return ztx.tx.Exec(query, args ...)
}

func (ztx *Ztx) Commit() error {
	return ztx.tx.Commit()
}

func (ztx *Ztx) Rollback() error {
	return ztx.tx.Rollback()
}