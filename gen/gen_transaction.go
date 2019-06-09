package gen

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
)

func genTransaction(file, folder, pkg string) error {
	// gen_<tb>_<db>.go
	base := "transaction.go"

	fn := filepath.Join(folder, base)



	fset := token.NewFileSet()
	fileAST, err := parser.ParseFile(fset, "", genTransactionFileString(pkg), parser.ParseComments)
	if err != nil {
		return err
	}
	ast.SortImports(fset, fileAST)

	var b bytes.Buffer

	if err = format.Node(&b, fset, fileAST); err != nil {
		return err
	}

	if err = ioutil.WriteFile(fn, b.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}

func genTransactionFileString(pkg string) []byte {
	b := NewBuf()
	b.WL("package orm")
	b.WL("import (")
	b.WL("	\"context\"")
	b.WL("	\"database/sql\"")
	b.WL("	\"github.com/envzo/zorm/db\"")
	b.WL(")")
	b.WL("")
	b.WL("type Ztx struct {")
	b.WL("	tx *sql.Tx")
	b.WL("}")
	b.WL("")
	b.WL("var Zotx = &Ztx{}")
	b.WL("func (ztx * Ztx) Begin() {")
	b.WL("	txTest, err := db.DB().Begin()")
	b.WL("	if err != nil {")
	b.WL("")
	b.WL("	}")
	b.WL("	ztx.tx = txTest")
	b.WL("}")
	b.WL("")
	b.WL("func (ztx *Ztx) Stmt(stmt *sql.Stmt) *sql.Stmt {")
	b.WL("	return ztx.tx.Stmt(stmt)")
	b.WL("}")
	b.WL("")
	b.WL("func (ztx *Ztx) StmtContext(ctx context.Context, stmt *sql.Stmt) *sql.Stmt {")
	b.WL("	return ztx.tx.StmtContext(ctx, stmt)")
	b.WL("}")
	b.WL("")
	b.WL("func (ztx *Ztx) Prepare(query string) (*sql.Stmt, error) {")
	b.WL("	return ztx.tx.Prepare(query)")
	b.WL("}")
	b.WL("")
	b.WL("func (ztx *Ztx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {")
	b.WL("	return ztx.tx.PrepareContext(ctx, query)")
	b.WL("}")
	b.WL("")
	b.WL("func (ztx *Ztx) Query(query string, args ...interface{}) (*sql.Rows, error) {")
	b.WL("	return ztx.tx.Query(query, args ...)")
	b.WL("}")
	b.WL("")
	b.WL("func (ztx *Ztx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {")
	b.WL("	return ztx.tx.QueryContext(ctx, query, args ...)")
	b.WL("}")
	b.WL("")
	b.WL("func (ztx *Ztx) QueryRow(query string, args ...interface{}) *sql.Row {")
	b.WL("	return ztx.tx.QueryRow(query, args ...)")
	b.WL("}")
	b.WL("")
	b.WL("func (ztx *Ztx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {")
	b.WL("	return ztx.tx.QueryRowContext(ctx, query, args ...)")
	b.WL("}")
	b.WL("")
	b.WL("func (ztx * Ztx) Exec(query string, args ...interface{}) (sql.Result, error) {")
	b.WL("	return ztx.tx.Exec(query, args ...)")
	b.WL("}")
	b.WL("")
	b.WL("func (ztx *Ztx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {")
	b.WL("	return ztx.tx.ExecContext(ctx, query, args ...)")
	b.WL("}")
	b.WL("")
	b.WL("func (ztx *Ztx) Commit() error {")
	b.WL("	return ztx.tx.Commit()")
	b.WL("}")
	b.WL("")
	b.WL("func (ztx *Ztx) Rollback() error {")
	b.WL("	return ztx.tx.Rollback()")
	b.WL("}")
	b.WL("")


	return b.Bytes()
}