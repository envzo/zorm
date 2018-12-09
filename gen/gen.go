// this is where magic happens
package gen

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"

	"github.com/envzo/zorm/parse"
)

type gen struct {
	T string
	D *parse.Def
	B *B
}

func Gen(file, folder, pkg string) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	m, err := parse.Parse(b)
	if err != nil {
		return err
	}

	for _, def := range m {
		g := &gen{
			T: ToCamel(def.TB),
			D: def,
			B: NewB(),
		}

		// gen_<tb>_<db>.go
		base := "gen_" + g.T + "_" + def.Engine

		fn := base + ".go"
		fn = filepath.Join(folder, fn)

		fset := token.NewFileSet()
		fileAST, err := parser.ParseFile(fset, "", g.genORM(pkg), parser.ParseComments)
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

		// gen_<tb>_<db>.sql

		fn = base + ".sql"
		fn = filepath.Join(folder, fn)
		if err = ioutil.WriteFile(fn, []byte(genSql(def)), 0644); err != nil {
			return err
		}
	}

	return nil
}
