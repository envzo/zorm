package gen

import (
	"bytes"

	"github.com/envzo/zorm/cls"
	"github.com/envzo/zorm/parse"
	"github.com/envzo/zorm/util"
)

/* todo
usage:
1 orm.Shop().Del().Where().XXX(xxx).Do()
2 orm.Shop().Sel().Where().XXX(xxx).YYY(yyy).Order().ZZZ(-1).AAA(1).Limit(a, b).Do()
*/

func (g *gen) genORM(pkg string) []byte {
	g.B.WL("// usage: ")
	g.B.WL2("// FindByXXX will not return sql.ErrNoRows, so it's caller's ability to check error")

	g.B.WL2("package ", pkg)

	g.B.WL(`import (`)
	g.B.WL(`"errors"`)
	g.B.WL(`"fmt"`)
	g.B.WL(`"database/sql"`)
	g.B.WL(`"strings"`)
	g.B.WL2(`"time"`)
	g.B.WL(`"github.com/ascode/zorm/db"`)
	g.B.WL(`"github.com/ascode/zorm/util"`)
	g.B.WL2(`)`)

	g.B.WL(`var _ = errors.New`)
	g.B.WL(`var _ = fmt.Printf`)
	g.B.WL(`var _ = strings.Trim`)
	g.B.WL(`var _ = sql.ErrNoRows`)
	g.B.WL(`var _ = util.I64`)
	g.B.WL(`var _ = time.Nanosecond`)

	g.B.WL("type ", g.T, " struct {")
	for _, f := range g.x.Fs {
		g.B.Tab().W(f.Camel).Spc().W(f.GoT).Ln()
	}
	g.B.Ln().WL("baby bool")
	g.B.WL("}")

	g.B.WL("func New", g.T, "() *", g.T, " {")
	g.B.WL("return &", g.T, "{baby: true}")
	g.B.WL("}")

	g.B.WL("type _", g.T, "Mgr struct {}")

	g.B.WL("var ", g.T, "Mgr = &_", g.T, "Mgr{}")

	for _, fs := range g.x.Uniques {
		g.genIsExists(fs)
		g.genTxIsExists(fs)
		g.genUniFind(fs)
		g.genTxUniFind(fs)
		g.genUniUpdate(fs)
		g.genTxUniUpdate(fs)
		g.genUniRm(fs)
		g.genTxUniRm(fs)
	}

	for _, fs := range g.x.Indexes {
		g.genFindByIndex(fs)
		g.genTxFindByIndex(fs)
		g.genCountByIndex(fs)
		g.genTxCountByIndex(fs)
	}

	g.genFindByMultiJoin()
	g.genCountByMultiJoin()
	g.genFindByJoin()
	g.genFindByCond()
	g.genFindAllByCond()
	g.genCreate().Ln()
	g.genTxCreate().Ln()
	g.genUpsert().Ln()
	g.genCountByRule()
	g.genRmByRule()
	g.genTxRmByRule()

	if g.x.PK != nil {
		g.genUniFindByPk()
		g.genTxUniFindByPk()
		g.genUpdateByPK().Ln()
		g.genTxUpdateByPK().Ln()
		g.genRmByPK()
		g.genTxRmByPK()
		g.genIsExistsByPK()
	}

	return g.B.Bytes()
}

func (g *gen) genIsExists(args []*parse.F) {
	var m bytes.Buffer
	m.WriteString("Is")
	for _, f := range args {
		m.WriteString(f.Camel)
	}
	m.WriteString("Exists")

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(")

	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(util.LowerFirstLetter(f.Camel)).Spc()

		if f.T == cls.YamlTimestamp { // it is convenient to use integer when querying
			g.B.W(util.I64)
		} else {
			g.B.W(f.GoT)
		}
	}
	g.B.W(")")
	g.B.Spc().W("(bool, error)").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")

	g.B.W("row := db.DB().QueryRow(`select count(1) from ", g.x.DB, ".", g.x.TB, " where ")

	for i, f := range args {
		if i > 0 {
			g.B.W(" and ")
		}
		g.B.W(f.Name, " = ?")
	}
	g.B.WL("`, ")

	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(util.LowerFirstLetter(f.Camel)).Spc()
	}
	g.B.WL2(")")

	g.B.WL("var c sql.NullInt64")
	g.B.Ln().W("if err := row.Scan(&c); err!= nil {")
	g.B.W("return false, err")
	g.B.WL("}")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")
	g.B.W("return c.Int64 > 0, nil")

	g.B.WL2("}")
}

func (g *gen) genTxIsExists(args []*parse.F) {
	var m bytes.Buffer
	m.WriteString("TxIs")
	for _, f := range args {
		m.WriteString(f.Camel)
	}
	m.WriteString("Exists")

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(")

	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(util.LowerFirstLetter(f.Camel)).Spc()

		if f.T == cls.YamlTimestamp { // it is convenient to use integer when querying
			g.B.W(util.I64)
		} else {
			g.B.W(f.GoT)
		}
	}
	g.B.W(")")
	g.B.Spc().W("(bool, error)").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")

	g.B.W("row := Zotx.QueryRow(`select count(1) from ", g.x.DB, ".", g.x.TB, " where ")

	for i, f := range args {
		if i > 0 {
			g.B.W(" and ")
		}
		g.B.W(f.Name, " = ?")
	}
	g.B.WL("`, ")

	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(util.LowerFirstLetter(f.Camel)).Spc()
	}
	g.B.WL2(")")

	g.B.WL("var c sql.NullInt64")
	g.B.Ln().W("if err := row.Scan(&c); err!= nil {")
	g.B.W("return false, err")
	g.B.WL("}")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")
	g.B.W("return c.Int64 > 0, nil")

	g.B.WL2("}")
}

func (g *gen) genIsExistsByPK() {
	var m bytes.Buffer
	m.WriteString("IsExistsByPK")

	g.B.WL("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(pk ", g.x.PK.GoT, ") (bool, error) {")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")
	g.B.WL("row := db.DB().QueryRow(`select count(1) from ", g.x.DB, ".", g.x.TB, " where ", g.x.PK.Name, " = ?`, pk)")
	g.B.WL("var c sql.NullInt64")
	g.B.Ln().WL("if err := row.Scan(&c); err!= nil {")
	g.B.WL("return false, err")
	g.B.WL("}")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")
	g.B.WL("return c.Int64 > 0, nil")
	g.B.WL2("}")
}

func (g *gen) genUniFind(args []*parse.F) {
	var m bytes.Buffer
	m.WriteString("UniFindBy")
	for _, f := range args {
		m.WriteString(f.Camel)
	}

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(")

	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(util.LowerFirstLetter(f.Camel)).Spc()

		if f.T == cls.YamlTimestamp { // it is convenient to use integer when querying
			g.B.W(util.I64)
		} else {
			g.B.W(f.GoT)
		}
	}
	g.B.W(")")
	g.B.Spc().W("(*" + g.T + ", error)").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")

	g.B.W("row := db.DB().QueryRow(`select ")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(f.Name)
	}
	g.B.W(" from ", g.x.DB, ".", g.x.TB, " where ")

	for i, f := range args {
		if i > 0 {
			g.B.W(" and ")
		}
		g.B.W(f.Name, " = ?")
	}
	g.B.WL("`, ")

	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		if f.T == cls.YamlDate {
			g.B.W(util.LowerFirstLetter(f.Camel), `.Format("2006-01-02")`).Spc()
		} else if f.T == cls.YamlDateTime {
			g.B.W(util.LowerFirstLetter(f.Camel), `.Format("2006-01-02 15:04:05")`).Spc()
		} else {
			g.B.W(util.LowerFirstLetter(f.Camel)).Spc()
		}
	}
	g.B.WL2(")")

	// temp variables
	vm := map[string]string{}

	for _, f := range g.x.Fs {
		n := util.LowerFirstLetter(f.Camel)

		if n == "type" {
			n += "_"
		}

		for _, arg := range args {
			if arg.Name == f.Name {
				n += "_1"
				break
			}
		}

		g.B.W("var ", n)
		vm[f.Camel] = n

		g.B.Spc().WL(util.NilSqlType(f.T))
	}

	g.B.Ln().W("if err := row.Scan(")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W("&", vm[f.Camel])
	}
	g.B.W("); err!= nil {")
	g.B.W("return nil, err")
	g.B.WL2("}")

	g.B.WL("d := ", g.T, "{}")
	for _, f := range g.x.Fs {
		g.B.W("d.", f.Camel, "=", util.DerefNilSqlType(vm[f.Camel], f.T)).Ln()
	}

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.W("return &d, nil")

	g.B.WL2("}")
}

func (g *gen) genTxUniFind(args []*parse.F) {
	var m bytes.Buffer
	m.WriteString("TxUniFindBy")
	for _, f := range args {
		m.WriteString(f.Camel)
	}

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(")

	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(util.LowerFirstLetter(f.Camel)).Spc()

		if f.T == cls.YamlTimestamp { // it is convenient to use integer when querying
			g.B.W(util.I64)
		} else {
			g.B.W(f.GoT)
		}
	}
	g.B.W(")")
	g.B.Spc().W("(*" + g.T + ", error)").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")

	g.B.W("row := Zotx.QueryRow(`select ")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(f.Name)
	}
	g.B.W(" from ", g.x.DB, ".", g.x.TB, " where ")

	for i, f := range args {
		if i > 0 {
			g.B.W(" and ")
		}
		g.B.W(f.Name, " = ?")
	}
	g.B.WL("`, ")

	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		if f.T == cls.YamlDate {
			g.B.W(util.LowerFirstLetter(f.Camel), `.Format("2006-01-02")`).Spc()
		} else if f.T == cls.YamlDateTime {
			g.B.W(util.LowerFirstLetter(f.Camel), `.Format("2006-01-02 15:04:05")`).Spc()
		} else {
			g.B.W(util.LowerFirstLetter(f.Camel)).Spc()
		}
	}
	g.B.WL2(")")

	// temp variables
	vm := map[string]string{}

	for _, f := range g.x.Fs {
		n := util.LowerFirstLetter(f.Camel)

		if n == "type" {
			n += "_"
		}

		for _, arg := range args {
			if arg.Name == f.Name {
				n += "_1"
				break
			}
		}

		g.B.W("var ", n)
		vm[f.Camel] = n

		g.B.Spc().WL(util.NilSqlType(f.T))
	}

	g.B.Ln().W("if err := row.Scan(")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W("&", vm[f.Camel])
	}
	g.B.W("); err!= nil {")
	g.B.W("return nil, err")
	g.B.WL2("}")

	g.B.WL("d := ", g.T, "{}")
	for _, f := range g.x.Fs {
		g.B.W("d.", f.Camel, "=", util.DerefNilSqlType(vm[f.Camel], f.T)).Ln()
	}

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.W("return &d, nil")

	g.B.WL2("}")
}

func (g *gen) genUniUpdate(args []*parse.F) {
	var m bytes.Buffer
	m.WriteString("UpdateBy")
	for _, f := range args {
		m.WriteString(f.Camel)
	}

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(d *", g.T, ")")
	g.B.Spc().W("(int64, error) ").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")
	g.B.W("r,err := db.DB().Exec(`update ", g.x.DB, ".", g.x.TB, " set ")
	flag := false
SetField:
	for _, f := range g.x.Fs {
		if g.x.PK != nil && f.Name == g.x.PK.Name && g.x.PK.AutoIncr {
			continue
		}

		for _, fields := range g.x.Uniques {
			for _, field := range fields {
				if f.Name == field.Name {
					continue SetField
				}
			}
		}

		if flag {
			g.B.W(", ")
		}
		g.B.W(f.Name, " = ?")
		flag = true
	}
	g.B.W(" where ")

	for i, f := range args {
		if i > 0 {
			g.B.W(" and ")
		}
		g.B.W(f.Name, " = ?")
	}
	g.B.W("`, ")

	// params
SetParam:
	for _, f := range g.x.Fs {
		if g.x.PK != nil && f.Name == g.x.PK.Name && g.x.PK.AutoIncr {
			continue
		}

		for _, fields := range g.x.Uniques {
			for _, field := range fields {
				if f.Name == field.Name {
					continue SetParam
				}
			}
		}

		g.B.W("d.", f.Camel, ", ")
	}

	for i, f := range args {
		g.B.W("d.", f.Camel)
		if i != len(args)-1 {
			g.B.W(", ")
		}
	}
	g.B.WL(")")
	g.B.WL("if err != nil {")
	g.B.WL("	return 0, err")
	g.B.WL("}")
	g.B.WL("n,err := r.RowsAffected()")
	g.B.WL("if err != nil {")
	g.B.WL("	return 0, err")
	g.B.WL("}")

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.WL("return n, nil")
	g.B.WL("}").Ln()
}

func (g *gen) genTxUniUpdate(args []*parse.F) {
	var m bytes.Buffer
	m.WriteString("TxUpdateBy")
	for _, f := range args {
		m.WriteString(f.Camel)
	}

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(d *", g.T, ")")
	g.B.Spc().W("(int64, error) ").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")
	g.B.W("r,err := Zotx.Exec(`update ", g.x.DB, ".", g.x.TB, " set ")
	flag := false
SetField:
	for _, f := range g.x.Fs {
		if g.x.PK != nil && f.Name == g.x.PK.Name && g.x.PK.AutoIncr {
			continue
		}

		for _, fields := range g.x.Uniques {
			for _, field := range fields {
				if f.Name == field.Name {
					continue SetField
				}
			}
		}

		if flag {
			g.B.W(", ")
		}
		g.B.W(f.Name, " = ?")
		flag = true
	}
	g.B.W(" where ")

	for i, f := range args {
		if i > 0 {
			g.B.W(" and ")
		}
		g.B.W(f.Name, " = ?")
	}
	g.B.W("`, ")

	// params
SetParam:
	for _, f := range g.x.Fs {
		if g.x.PK != nil && f.Name == g.x.PK.Name && g.x.PK.AutoIncr {
			continue
		}

		for _, fields := range g.x.Uniques {
			for _, field := range fields {
				if f.Name == field.Name {
					continue SetParam
				}
			}
		}

		g.B.W("d.", f.Camel, ", ")
	}

	for i, f := range args {
		g.B.W("d.", f.Camel)
		if i != len(args)-1 {
			g.B.W(", ")
		}
	}
	g.B.WL(")")
	g.B.WL("if err != nil {")
	g.B.WL("	return 0, err")
	g.B.WL("}")
	g.B.WL("n,err := r.RowsAffected()")
	g.B.WL("if err != nil {")
	g.B.WL("	return 0, err")
	g.B.WL("}")

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.WL("return n, nil")
	g.B.WL("}").Ln()
}

func (g *gen) genUniRm(args []*parse.F) {
	var m bytes.Buffer
	m.WriteString("UniRmBy")
	for _, f := range args {
		m.WriteString(f.Camel)
	}

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(d *", g.T, ")")
	g.B.Spc().W("(int64, error) ").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")
	g.B.W("r,err := db.DB().Exec(`delete from ", g.x.DB, ".", g.x.TB, " where ")

	for i, f := range args {
		if i > 0 {
			g.B.W(" and ")
		}
		g.B.W(f.Name, " = ?")
	}
	g.B.W("`, ")

	// params
	for i, f := range args {
		g.B.W("d.", f.Camel)
		if i != len(args)-1 {
			g.B.W(", ")
		}
	}
	g.B.WL(")")
	g.B.WL("if err != nil {")
	g.B.WL("	return 0, err")
	g.B.WL("}")
	g.B.WL("n,err := r.RowsAffected()")
	g.B.WL("if err != nil {")
	g.B.WL("	return 0, err")
	g.B.WL("}")

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.WL("return n, nil")
	g.B.WL("}").Ln()
}

func (g *gen) genTxUniRm(args []*parse.F) {
	var m bytes.Buffer
	m.WriteString("TxUniRmBy")
	for _, f := range args {
		m.WriteString(f.Camel)
	}

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(d *", g.T, ")")
	g.B.Spc().W("(int64, error) ").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")
	g.B.W("r,err := Zotx.Exec(`delete from ", g.x.DB, ".", g.x.TB, " where ")

	for i, f := range args {
		if i > 0 {
			g.B.W(" and ")
		}
		g.B.W(f.Name, " = ?")
	}
	g.B.W("`, ")

	// params
	for i, f := range args {
		g.B.W("d.", f.Camel)
		if i != len(args)-1 {
			g.B.W(", ")
		}
	}
	g.B.WL(")")
	g.B.WL("if err != nil {")
	g.B.WL("	return 0, err")
	g.B.WL("}")
	g.B.WL("n,err := r.RowsAffected()")
	g.B.WL("if err != nil {")
	g.B.WL("	return 0, err")
	g.B.WL("}")

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.WL("return n, nil")
	g.B.WL("}").Ln()
}

func (g *gen) genCreate() *Buf {
	var m bytes.Buffer
	m.WriteString("Create")

	g.B.WL("func (mgr *_", g.T, "Mgr) ", m.String(), "(d *", g.T, ") error {")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")
	if g.x.PK != nil && g.x.PK.AutoIncr {
		g.B.W("r, err")
	} else {
		g.B.W("_, err")
	}
	g.B.W(" := db.DB().Exec(`insert into ", g.x.DB, ".", g.x.TB, " (")
	cnt := 0
	for i, f := range g.x.Fs {
		if g.x.PK != nil && g.x.PK.AutoIncr && f.Name == g.x.PK.Name {
			continue
		}
		g.B.W(f.Name)
		if i != len(g.x.Fs)-1 {
			g.B.W(", ")
		}
		cnt++
	}
	g.B.W(") value (")
	for i := 0; i < cnt; i++ {
		if i > 0 {
			g.B.W(",")
		}
		g.B.W("?")
	}
	g.B.W(")`,")

	for i, f := range g.x.Fs {
		if g.x.PK != nil && g.x.PK.AutoIncr && f.Name == g.x.PK.Name {
			continue
		}
		g.B.W("d.", f.Camel)
		if i != len(g.x.Fs)-1 {
			g.B.W(",")
		}
	}

	g.B.WL(")")
	g.B.WL("if err != nil {")
	g.B.WL("	return err")
	g.B.WL("}")

	if g.x.PK != nil && g.x.PK.AutoIncr {
		g.B.WL("id,err := r.LastInsertId()")
		g.B.WL("if err != nil {")
		g.B.WL("return err")
		g.B.WL("}")
		g.B.W("d.", g.x.PK.Camel, "=")

		if g.x.PK.T == cls.YamlI64 {
			g.B.WL("id")
		} else if g.x.PK.T == cls.YamlI32 {
			g.B.WL("int32(id)")
		}
	}

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.WL("return nil")
	return g.B.WL("}")
}

func (g *gen) genTxCreate() *Buf {
	var m bytes.Buffer
	m.WriteString("TxCreate")

	g.B.WL("func (mgr *_", g.T, "Mgr) ", m.String(), "(d *", g.T, ") error {")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")
	if g.x.PK != nil && g.x.PK.AutoIncr {
		g.B.W("r, err")
	} else {
		g.B.W("_, err")
	}
	g.B.W(" := Zotx.Exec(`insert into ", g.x.DB, ".", g.x.TB, " (")
	cnt := 0
	for i, f := range g.x.Fs {
		if g.x.PK != nil && g.x.PK.AutoIncr && f.Name == g.x.PK.Name {
			continue
		}
		g.B.W(f.Name)
		if i != len(g.x.Fs)-1 {
			g.B.W(", ")
		}
		cnt++
	}
	g.B.W(") value (")
	for i := 0; i < cnt; i++ {
		if i > 0 {
			g.B.W(",")
		}
		g.B.W("?")
	}
	g.B.W(")`,")

	for i, f := range g.x.Fs {
		if g.x.PK != nil && g.x.PK.AutoIncr && f.Name == g.x.PK.Name {
			continue
		}
		g.B.W("d.", f.Camel)
		if i != len(g.x.Fs)-1 {
			g.B.W(",")
		}
	}

	g.B.WL(")")
	g.B.WL("if err != nil {")
	g.B.WL("	return err")
	g.B.WL("}")

	if g.x.PK != nil && g.x.PK.AutoIncr {
		g.B.WL("id,err := r.LastInsertId()")
		g.B.WL("if err != nil {")
		g.B.WL("return err")
		g.B.WL("}")
		g.B.W("d.", g.x.PK.Camel, "=")

		if g.x.PK.T == cls.YamlI64 {
			g.B.WL("id")
		} else if g.x.PK.T == cls.YamlI32 {
			g.B.WL("int32(id)")
		}
	}

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.WL("return nil")
	return g.B.WL("}")
}

func (g *gen) genUpsert() *Buf {
	g.B.WL("func (mgr *_", g.T, "Mgr) Upsert(d *", g.T, ") error {").
		WL("if d.baby {").
		WL("	return mgr.Create(d)").
		WL("}")
	if g.x.PK != nil {
		g.B.WL("_, err := mgr.Update(d)").
			WL("return err")
	} else {
		g.B.WL(`return errors.New("unimplemented upsert: maybe adding pk is a good option")`)
	}
	return g.B.WL("}")
}

func (g *gen) genUniFindByPk() {
	var m bytes.Buffer
	m.WriteString("UniFindByPK")

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(", util.LowerFirstLetter(g.x.PK.Camel))

	if g.x.PK.T == cls.YamlTimestamp { // it is convenient to use integer when querying
		g.B.Spc().W(util.I64)
	} else {
		g.B.Spc().W(g.x.PK.GoT)
	}
	g.B.W(")")
	g.B.Spc().W("(*" + g.T + ", error)").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")

	g.B.W("row := db.DB().QueryRow(`select ")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(f.Name)
	}
	g.B.WL2(" from ", g.x.DB, ".", g.x.TB, " where ", g.x.PK.Name, " = ?`, ", util.LowerFirstLetter(g.x.PK.Camel), ")")

	// temp variables
	vm := map[string]string{}

	for _, f := range g.x.Fs {
		n := util.LowerFirstLetter(f.Camel)

		if n == "type" {
			n += "_"
		}

		if f.Name == g.x.PK.Name {
			n += "_1"
		}

		g.B.W("var ", n)
		vm[f.Camel] = n

		g.B.Spc().WL(util.NilSqlType(f.T))
	}

	g.B.Ln().W("if err := row.Scan(")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W("&", vm[f.Camel])
	}
	g.B.W("); err != nil {")
	g.B.W("return nil, err")
	g.B.WL2("}")

	g.B.WL("d := ", g.T, "{}")
	for _, f := range g.x.Fs {
		g.B.W("d.", f.Camel, "=", util.DerefNilSqlType(vm[f.Camel], f.T)).Ln()
	}

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.W("return &d, nil")

	g.B.WL2("}")
}

func (g *gen) genTxUniFindByPk() {
	var m bytes.Buffer
	m.WriteString("TxUniFindByPK")

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(", util.LowerFirstLetter(g.x.PK.Camel))

	if g.x.PK.T == cls.YamlTimestamp { // it is convenient to use integer when querying
		g.B.Spc().W(util.I64)
	} else {
		g.B.Spc().W(g.x.PK.GoT)
	}
	g.B.W(")")
	g.B.Spc().W("(*" + g.T + ", error)").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")

	g.B.W("row := Zotx.QueryRow(`select ")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(f.Name)
	}
	g.B.WL2(" from ", g.x.DB, ".", g.x.TB, " where ", g.x.PK.Name, " = ?`, ", util.LowerFirstLetter(g.x.PK.Camel), ")")

	// temp variables
	vm := map[string]string{}

	for _, f := range g.x.Fs {
		n := util.LowerFirstLetter(f.Camel)

		if n == "type" {
			n += "_"
		}

		if f.Name == g.x.PK.Name {
			n += "_1"
		}

		g.B.W("var ", n)
		vm[f.Camel] = n

		g.B.Spc().WL(util.NilSqlType(f.T))
	}

	g.B.Ln().W("if err := row.Scan(")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W("&", vm[f.Camel])
	}
	g.B.W("); err != nil {")
	g.B.W("return nil, err")
	g.B.WL2("}")

	g.B.WL("d := ", g.T, "{}")
	for _, f := range g.x.Fs {
		g.B.W("d.", f.Camel, "=", util.DerefNilSqlType(vm[f.Camel], f.T)).Ln()
	}

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.W("return &d, nil")

	g.B.WL2("}")
}

func (g *gen) genUpdateByPK() *Buf {
	var m bytes.Buffer
	m.WriteString("Update")

	g.B.WL("func (mgr *_", g.T, "Mgr) ", m.String(), "(d *", g.T, ") (int64, error) {")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")
	g.B.W("r,err:=db.DB().Exec(`update ", g.x.DB, ".", g.x.TB, " set ")
	for i, f := range g.x.Fs {
		if f.Name == g.x.PK.Name {
			continue
		}
		g.B.W(f.Name, " = ?")
		if i != len(g.x.Fs)-1 {
			g.B.W(", ")
		}
	}
	g.B.W(" where ", g.x.PK.Name, " = ?`, ")

	// params
	for i, f := range g.x.Fs {
		if f.Name == g.x.PK.Name {
			continue
		}
		g.B.W("d.", f.Camel)
		if i != len(g.x.Fs)-1 {
			g.B.W(", ")
		}
	}

	g.B.WL(", d.", g.x.PK.Camel, ")")
	g.B.WL("if err != nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")
	g.B.WL("n,err := r.RowsAffected()")
	g.B.WL("if err != nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")
	g.B.WL("return n, nil")
	return g.B.WL("}")
}

func (g *gen) genTxUpdateByPK() *Buf {
	var m bytes.Buffer
	m.WriteString("TxUpdate")

	g.B.WL("func (mgr *_", g.T, "Mgr) ", m.String(), "(d *", g.T, ") (int64, error) {")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")
	g.B.W("r,err:=Zotx.Exec(`update ", g.x.DB, ".", g.x.TB, " set ")
	for i, f := range g.x.Fs {
		if f.Name == g.x.PK.Name {
			continue
		}
		g.B.W(f.Name, " = ?")
		if i != len(g.x.Fs)-1 {
			g.B.W(", ")
		}
	}
	g.B.W(" where ", g.x.PK.Name, " = ?`, ")

	// params
	for i, f := range g.x.Fs {
		if f.Name == g.x.PK.Name {
			continue
		}
		g.B.W("d.", f.Camel)
		if i != len(g.x.Fs)-1 {
			g.B.W(", ")
		}
	}

	g.B.WL(", d.", g.x.PK.Camel, ")")
	g.B.WL("if err != nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")
	g.B.WL("n,err := r.RowsAffected()")
	g.B.WL("if err != nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")
	g.B.WL("return n, nil")
	return g.B.WL("}")
}

func (g *gen) genTxRmByPK() {
	var m bytes.Buffer
	m.WriteString("TxRmByPK")

	g.B.WL("func (mgr *_", g.T, "Mgr) ", m.String(), "(pk ", g.x.PK.GoT, ") (int64, error) {")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")
	g.B.WL(`query := "delete from `, g.x.DB, ".", g.x.TB, " where ", g.x.PK.Name, ` = ?"`)

	g.B.WL("r,err := Zotx.Exec(query, pk)")
	g.B.WL("if err != nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")
	g.B.WL("n,err := r.RowsAffected()")
	g.B.WL("if err != nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.WL("return n, nil")
	g.B.WL("}")
}

func (g *gen) genRmByPK() {
	var m bytes.Buffer
	m.WriteString("RmByPK")

	g.B.WL("func (mgr *_", g.T, "Mgr) ", m.String(), "(pk ", g.x.PK.GoT, ") (int64, error) {")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")
	g.B.WL(`query := "delete from `, g.x.DB, ".", g.x.TB, " where ", g.x.PK.Name, ` = ?"`)

	g.B.WL("r,err := db.DB().Exec(query, pk)")
	g.B.WL("if err != nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")
	g.B.WL("n,err := r.RowsAffected()")
	g.B.WL("if err != nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.WL("return n, nil")
	g.B.WL("}")
}

func (g *gen) genRmByRule() {
	var m bytes.Buffer
	m.WriteString("RmByRule")

	g.B.WL("func (mgr *_", g.T, "Mgr) ", m.String(), "(rules ...db.Rule) (int64, error) {")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")
	g.B.WL(`query := "delete from `, g.x.DB, ".", g.x.TB, ` where "`)

	g.B.WL("var p []interface{}")

	// params
	g.B.WL(`for i, r := range rules {`)
	g.B.WL(`if i > 0 {`)
	g.B.WL(`	query += " and "`)
	g.B.WL(`}`)
	g.B.WL(`	query += r.S`)
	g.B.WL(`	p = append(p, r.P)`)
	g.B.WL(`}`)

	g.B.WL("r,err := db.DB().Exec(query, p...)")
	g.B.WL("if err != nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")
	g.B.WL("n,err := r.RowsAffected()")
	g.B.WL("if err!=nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.WL("return n, nil")
	g.B.WL("}")
}

func (g *gen) genTxRmByRule() {
	var m bytes.Buffer
	m.WriteString("TxRmByRule")

	g.B.WL("func (mgr *_", g.T, "Mgr) ", m.String(), "(rules ...db.Rule) (int64, error) {")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")
	g.B.WL(`query := "delete from `, g.x.DB, ".", g.x.TB, ` where "`)

	g.B.WL("var p []interface{}")

	// params
	g.B.WL(`for i, r := range rules {`)
	g.B.WL(`if i > 0 {`)
	g.B.WL(`	query += " and "`)
	g.B.WL(`}`)
	g.B.WL(`	query += r.S`)
	g.B.WL(`	p = append(p, r.P)`)
	g.B.WL(`}`)

	g.B.WL("r,err := Zotx.Exec(query, p...)")
	g.B.WL("if err != nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")
	g.B.WL("n,err := r.RowsAffected()")
	g.B.WL("if err!=nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.WL("return n, nil")
	g.B.WL("}")
}

func (g *gen) genFindByIndex(args []*parse.F) {
	var m bytes.Buffer
	m.WriteString("FindBy")
	for _, f := range args {
		m.WriteString(f.Camel)
	}

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(")

	for i, arg := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(util.LowerFirstLetter(arg.Camel)).Spc()

		if arg.T == cls.YamlTimestamp { // it is convenient to use integer when querying
			g.B.W(util.I64)
		} else {
			g.B.W(arg.GoT)
		}
	}
	g.B.W(", order []string, offset, limit int64)")
	g.B.Spc().W("([]*" + g.T + ", error)").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")

	// make query sel
	g.B.W("query := `select ")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(f.Name)
	}
	g.B.W(" from ", g.x.DB, ".", g.x.TB, " where ")
	for i, f := range args {
		if i > 0 {
			g.B.W(" and ")
		}
		g.B.W(f.Name, " = ?")
	}
	g.B.WL("`")
	g.B.WL("for i, o := range order {")
	g.B.WL("if i==0 {")
	g.B.WL(`query += " order by "`)
	g.B.WL("} else {")
	g.B.WL(`query += ", "`)
	g.B.WL("}")
	g.B.WL(`if strings.HasPrefix(o, "-") {`)
	g.B.WL("	query += o[1:]")
	g.B.WL("} else {")
	g.B.WL("	query += o")
	g.B.WL("}")
	g.B.WL("if o[0] == '-' {")
	g.B.WL(`query += " desc"`)
	g.B.WL("}")
	g.B.WL("}")
	g.B.WL("if offset != -1 && limit != -1 {")
	g.B.WL(`query += fmt.Sprintf(" limit %d, %d", offset, limit)`)
	g.B.WL2("}")
	// end make query sql

	g.B.W("rows, err := db.DB().Query(query, ")
	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(util.LowerFirstLetter(f.Camel)).Spc()
	}
	g.B.WL(")")
	g.B.WL("if err != nil {")
	g.B.WL("return nil, err")
	g.B.WL("}")
	g.B.WL2("defer rows.Close()")

	// temp variables
	vm := map[string]string{}

	for _, f := range g.x.Fs {
		n := util.LowerFirstLetter(f.Camel)

		if n == "type" {
			n += "_"
		}

		// todo need refine
		for _, arg := range args {
			if arg.Name == f.Name {
				n += "_1"
				break
			}
		}

		g.B.W("var ", n)
		vm[f.Camel] = n

		g.B.Spc().WL(util.NilSqlType(f.T))
	}

	g.B.Ln().WL2("var ret []*", g.T)

	g.B.WL("for rows.Next(){")

	g.B.W("if err = rows.Scan(")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W("&", vm[f.Camel])
	}
	g.B.W("); err!= nil {")
	g.B.W("return nil, err")
	g.B.WL2("}")

	g.B.WL("d := ", g.T, "{}")
	for _, f := range g.x.Fs {
		g.B.W("d.", f.Camel, "=", util.DerefNilSqlType(vm[f.Camel], f.T)).Ln()
	}

	g.B.WL("ret = append(ret, &d)")

	g.B.WL("}") // end rows loop

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.W("return ret, nil")

	g.B.WL2("}")
}

func (g *gen) genTxFindByIndex(args []*parse.F) {
	var m bytes.Buffer
	m.WriteString("TxFindBy")
	for _, f := range args {
		m.WriteString(f.Camel)
	}

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(")

	for i, arg := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(util.LowerFirstLetter(arg.Camel)).Spc()

		if arg.T == cls.YamlTimestamp { // it is convenient to use integer when querying
			g.B.W(util.I64)
		} else {
			g.B.W(arg.GoT)
		}
	}
	g.B.W(", order []string, offset, limit int64)")
	g.B.Spc().W("([]*" + g.T + ", error)").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")

	// make query sel
	g.B.W("query := `select ")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(f.Name)
	}
	g.B.W(" from ", g.x.DB, ".", g.x.TB, " where ")
	for i, f := range args {
		if i > 0 {
			g.B.W(" and ")
		}
		g.B.W(f.Name, " = ?")
	}
	g.B.WL("`")
	g.B.WL("for i, o := range order {")
	g.B.WL("if i==0 {")
	g.B.WL(`query += " order by "`)
	g.B.WL("} else {")
	g.B.WL(`query += ", "`)
	g.B.WL("}")
	g.B.WL(`if strings.HasPrefix(o, "-") {`)
	g.B.WL("	query += o[1:]")
	g.B.WL("} else {")
	g.B.WL("	query += o")
	g.B.WL("}")
	g.B.WL("if o[0] == '-' {")
	g.B.WL(`query += " desc"`)
	g.B.WL("}")
	g.B.WL("}")
	g.B.WL("if offset != -1 && limit != -1 {")
	g.B.WL(`query += fmt.Sprintf(" limit %d, %d", offset, limit)`)
	g.B.WL2("}")
	// end make query sql

	g.B.W("rows, err := Zotx.Query(query, ")
	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(util.LowerFirstLetter(f.Camel)).Spc()
	}
	g.B.WL(")")
	g.B.WL("if err != nil {")
	g.B.WL("return nil, err")
	g.B.WL("}")
	g.B.WL2("defer rows.Close()")

	// temp variables
	vm := map[string]string{}

	for _, f := range g.x.Fs {
		n := util.LowerFirstLetter(f.Camel)

		if n == "type" {
			n += "_"
		}

		// todo need refine
		for _, arg := range args {
			if arg.Name == f.Name {
				n += "_1"
				break
			}
		}

		g.B.W("var ", n)
		vm[f.Camel] = n

		g.B.Spc().WL(util.NilSqlType(f.T))
	}

	g.B.Ln().WL2("var ret []*", g.T)

	g.B.WL("for rows.Next(){")

	g.B.W("if err = rows.Scan(")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W("&", vm[f.Camel])
	}
	g.B.W("); err!= nil {")
	g.B.W("return nil, err")
	g.B.WL2("}")

	g.B.WL("d := ", g.T, "{}")
	for _, f := range g.x.Fs {
		g.B.W("d.", f.Camel, "=", util.DerefNilSqlType(vm[f.Camel], f.T)).Ln()
	}

	g.B.WL("ret = append(ret, &d)")

	g.B.WL("}") // end rows loop

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.W("return ret, nil")

	g.B.WL2("}")
}

func (g *gen) genFindByJoin() {
	g.B.W("func (mgr", " *_", g.T, "Mgr) FindByJoin(t string, on, where []db.Rule, order []string, offset, limit int64) ")
	g.B.W("([]*" + g.T + ", error)").WL("{")
	g.B.WL("return mgr.FindByMultiJoin([]db.Join{")
	g.B.WL("	{T: t, Rule: on},")
	g.B.WL("}, where, order, offset, limit)")
	g.B.WL("}")
}

func (g *gen) genFindByMultiJoin() {
	var m bytes.Buffer
	m.WriteString("FindByMultiJoin")

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(joins []db.Join, where []db.Rule, order []string, offset, limit int64)")
	g.B.Spc().W("([]*" + g.T + ", error)").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")

	g.B.WL2("var params []interface{}")

	// make query sel
	g.B.W("query := `select ")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(g.x.TB, ".", f.Name)
	}
	g.B.WL(" from ", g.x.DB, ".", g.x.TB, "`")

	g.B.WL("for _, join := range joins {")
	g.B.WL("query += ` join ", g.x.DB, ".` + ", "join.T + ` on `")
	g.B.WL(`for i, v := range join.Rule {`)
	g.B.WL(`	if i > 0 {`)
	g.B.WL(`		query += " and "`)
	g.B.WL(`	}`)
	g.B.WL(`	query += v.S`)
	g.B.WL(`	if v.P != nil {`)
	g.B.WL(`		params = append(params, v.P)`)
	g.B.WL(`	}`)
	g.B.WL(`}`)
	g.B.WL("}")

	g.B.WL(`for i, v := range where {`)
	g.B.WL(`	if i == 0 {`)
	g.B.WL(`		query += " where "`)
	g.B.WL(`	} else {`)
	g.B.WL(`		query += " and "`)
	g.B.WL(`	}`)
	g.B.WL(`	query += v.S`)
	g.B.WL(`	if v.P != nil {`)
	g.B.WL(`		params = append(params, v.P)`)
	g.B.WL(`	}`)
	g.B.WL(`}`)

	g.B.WL("for i, o := range order {")
	g.B.WL("if i == 0 {")
	g.B.WL(`query += " order by "`)
	g.B.WL("} else {")
	g.B.WL(`query += ", "`)
	g.B.WL("}")
	g.B.WL(`if strings.HasPrefix(o, "-") {`)
	g.B.WL("	query += o[1:]")
	g.B.WL("} else {")
	g.B.WL("	query += o")
	g.B.WL("}")
	g.B.WL("if o[0] == '-' {")
	g.B.WL(`query += " desc"`)
	g.B.WL("}")
	g.B.WL("}")
	g.B.WL("if offset != -1 && limit != -1 {")
	g.B.WL(`query += fmt.Sprintf(" limit %d, %d", offset, limit)`)
	g.B.WL2("}")
	// end make query sql

	g.B.WL("rows, err := db.DB().Query(query, params...)")
	g.B.WL("if err != nil {")
	g.B.WL("return nil, err")
	g.B.WL("}")
	g.B.WL2("defer rows.Close()")

	// temp variables
	vm := map[string]string{}

	args := map[string]bool{
		"t":      true,
		"on":     true,
		"where":  true,
		"order":  true,
		"offset": true,
		"limit":  true,
	}

	for _, f := range g.x.Fs {
		n := util.LowerFirstLetter(f.Camel)

		if n == "type" {
			n += "_"
		}

		if _, ok := args[f.Name]; ok {
			n += "_1"
		}

		g.B.W("var ", n)
		vm[f.Camel] = n

		g.B.Spc().WL(util.NilSqlType(f.T))
	}

	g.B.Ln().WL2("var ret []*", g.T)

	g.B.WL("for rows.Next() {")

	g.B.W("if err = rows.Scan(")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W("&", vm[f.Camel])
	}
	g.B.W("); err != nil {")
	g.B.W("return nil, err")
	g.B.WL2("}")

	g.B.WL("d := ", g.T, "{}")
	for _, f := range g.x.Fs {
		g.B.W("d.", f.Camel, "=")
		g.B.W(util.DerefNilSqlType(vm[f.Camel], f.T)).Ln()
	}

	g.B.WL("ret = append(ret, &d)")

	g.B.WL("}") // end rows loop

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.W("return ret, nil")

	g.B.WL2("}")
}

func (g *gen) genCountByMultiJoin() {
	var m bytes.Buffer
	m.WriteString("CountByMultiJoin")

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(joins []db.Join, where []db.Rule) ")
	g.B.W("(int64, error)").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")

	g.B.WL2("var params []interface{}")

	// make query sel
	g.B.W("query := `select count(1) from (select ")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(g.x.TB, ".", f.Name)
	}
	g.B.WL(" from ", g.x.DB, ".", g.x.TB, "`")

	g.B.WL("for _, join := range joins {")
	g.B.WL("query += ` join ", g.x.DB, ".` + ", "join.T + ` on `")
	g.B.WL(`for i, v := range join.Rule {`)
	g.B.WL(`	if i > 0 {`)
	g.B.WL(`		query += " and "`)
	g.B.WL(`	}`)
	g.B.WL(`	query += v.S`)
	g.B.WL(`	if v.P != nil {`)
	g.B.WL(`		params = append(params, v.P)`)
	g.B.WL(`	}`)
	g.B.WL(`}`)
	g.B.WL("}")

	g.B.WL(`for i, v := range where {`)
	g.B.WL(`	if i == 0 {`)
	g.B.WL(`		query += " where "`)
	g.B.WL(`	} else {`)
	g.B.WL(`		query += " and "`)
	g.B.WL(`	}`)
	g.B.WL(`	query += v.S`)
	g.B.WL(`	if v.P != nil {`)
	g.B.WL(`		params = append(params, v.P)`)
	g.B.WL(`	}`)
	g.B.WL(`}`)
	g.B.WL2(`query += ") t"`)

	// end make query sql

	g.B.WL("row := db.DB().QueryRow(query, params...)")

	g.B.WL("var c sql.NullInt64")
	g.B.WL("if err := row.Scan(&c); err != nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.W("return c.Int64, nil")

	g.B.WL2("}")
}

func (g *gen) genFindByCond() {
	var m bytes.Buffer
	m.WriteString("FindByCond")

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(where []db.Rule, order []string, offset, limit int64)")
	g.B.Spc().W("([]*" + g.T + ", error)").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")

	g.B.WL2("var params []interface{}")

	// make query sel
	g.B.W("query := `select ")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(f.Name)
	}
	g.B.WL(" from ", g.x.DB, ".", g.x.TB, " where `")
	g.B.WL(`for i, v := range where {`)
	g.B.WL(`	if i > 0 {`)
	g.B.WL(`		query += " and "`)
	g.B.WL(`	}`)
	g.B.WL(`	query += v.S`)
	g.B.WL(`	if v.P != nil {`)
	g.B.WL(`		params = append(params, v.P)`)
	g.B.WL(`	}`)
	g.B.WL(`}`)

	g.B.WL("for i, o := range order {")
	g.B.WL("if i == 0 {")
	g.B.WL(`query += " order by "`)
	g.B.WL("} else {")
	g.B.WL(`query += ", "`)
	g.B.WL("}")
	g.B.WL(`if strings.HasPrefix(o, "-") {`)
	g.B.WL("	query += o[1:]")
	g.B.WL("} else {")
	g.B.WL("	query += o")
	g.B.WL("}")
	g.B.WL("if o[0] == '-' {")
	g.B.WL(`query += " desc"`)
	g.B.WL("}")
	g.B.WL("}")
	g.B.WL("if offset != -1 && limit != -1 {")
	g.B.WL(`query += fmt.Sprintf(" limit %d, %d", offset, limit)`)
	g.B.WL2("}")
	// end make query sql

	g.B.WL("rows, err := db.DB().Query(query, params...)")
	g.B.WL("if err != nil {")
	g.B.WL("return nil, err")
	g.B.WL("}")
	g.B.WL2("defer rows.Close()")

	// temp variables
	vm := map[string]string{}

	args := map[string]bool{
		"t":      true,
		"on":     true,
		"where":  true,
		"order":  true,
		"offset": true,
		"limit":  true,
	}

	for _, f := range g.x.Fs {
		n := util.LowerFirstLetter(f.Camel)

		if n == "type" {
			n += "_"
		}

		if _, ok := args[f.Name]; ok {
			n += "_1"
		}

		g.B.W("var ", n)
		vm[f.Camel] = n

		g.B.Spc().WL(util.NilSqlType(f.T))
	}

	g.B.Ln().WL2("var ret []*", g.T)

	g.B.WL("for rows.Next() {")

	g.B.W("if err = rows.Scan(")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W("&", vm[f.Camel])
	}
	g.B.W("); err != nil {")
	g.B.W("return nil, err")
	g.B.WL2("}")

	g.B.WL("d := ", g.T, "{}")
	for _, f := range g.x.Fs {
		g.B.W("d.", f.Camel, "=", util.DerefNilSqlType(vm[f.Camel], f.T)).Ln()
	}

	g.B.WL("ret = append(ret, &d)")

	g.B.WL("}") // end rows loop

	g.B.W("return ret, nil")

	g.B.WL2("}")
}

func (g *gen) genFindAllByCond() {
	var m bytes.Buffer
	m.WriteString("FindAllByCond")

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(where []db.Rule, order []string)")
	g.B.Spc().W("([]*" + g.T + ", error)").WL("{")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")

	g.B.WL2("var params []interface{}")

	// make query sel
	g.B.W("query := `select ")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(f.Name)
	}
	g.B.WL(" from ", g.x.DB, ".", g.x.TB, " `")
	g.B.WL(`for i, v := range where {`)
	g.B.WL(`	if i == 0 {`)
	g.B.WL(`		query += " where "`)
	g.B.WL(`	} else if i > 0 {`)
	g.B.WL(`		query += " and "`)
	g.B.WL(`	}`)
	g.B.WL(`	query += v.S`)
	g.B.WL(`	if v.P != nil {`)
	g.B.WL(`		params = append(params, v.P)`)
	g.B.WL(`	}`)
	g.B.WL(`}`)

	g.B.WL("for i, o := range order {")
	g.B.WL("if i == 0 {")
	g.B.WL(`query += " order by "`)
	g.B.WL("} else {")
	g.B.WL(`query += ", "`)
	g.B.WL("}")
	g.B.WL(`if strings.HasPrefix(o, "-") {`)
	g.B.WL("	query += o[1:]")
	g.B.WL("} else {")
	g.B.WL("	query += o")
	g.B.WL("}")
	g.B.WL("if o[0] == '-' {")
	g.B.WL(`query += " desc"`)
	g.B.WL("}")
	g.B.WL("}")
	// end make query sql

	g.B.WL("rows, err := db.DB().Query(query, params...)")
	g.B.WL("if err != nil {")
	g.B.WL("return nil, err")
	g.B.WL("}")
	g.B.WL2("defer rows.Close()")

	// temp variables
	vm := map[string]string{}

	args := map[string]bool{
		"t":      true,
		"on":     true,
		"where":  true,
		"order":  true,
		"offset": true,
		"limit":  true,
	}

	for _, f := range g.x.Fs {
		n := util.LowerFirstLetter(f.Camel)

		if n == "type" {
			n += "_"
		}

		if _, ok := args[f.Name]; ok {
			n += "_1"
		}

		g.B.W("var ", n)
		vm[f.Camel] = n

		g.B.Spc().WL(util.NilSqlType(f.T))
	}

	g.B.Ln().WL2("var ret []*", g.T)

	g.B.WL("for rows.Next() {")

	g.B.W("if err = rows.Scan(")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W("&", vm[f.Camel])
	}
	g.B.W("); err != nil {")
	g.B.W("return nil, err")
	g.B.WL2("}")

	g.B.WL("d := ", g.T, "{}")
	for _, f := range g.x.Fs {
		g.B.W("d.", f.Camel, "=", util.DerefNilSqlType(vm[f.Camel], f.T)).Ln()
	}

	g.B.WL("ret = append(ret, &d)")

	g.B.WL("}") // end rows loop

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.W("return ret, nil")

	g.B.WL2("}")
}

func (g *gen) genCountByIndex(args []*parse.F) {
	var m bytes.Buffer
	m.WriteString("CountBy")
	for _, f := range args {
		m.WriteString(f.Camel)
	}

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(")

	for i, arg := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(util.LowerFirstLetter(arg.Camel)).Spc()

		if arg.T == cls.YamlTimestamp { // it is convenient to use integer when querying
			g.B.W(util.I64)
		} else {
			g.B.W(arg.GoT)
		}
	}
	g.B.W(")")
	g.B.Spc().W("(int64, error)").WL("{")

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")

	// make query sel
	g.B.W("query := `select count(1) from ", g.x.DB, ".", g.x.TB, " where ")
	for i, f := range args {
		if i > 0 {
			g.B.W(" and ")
		}
		g.B.W(f.Name, " = ?")
	}
	g.B.WL("`")
	// end make query sql

	g.B.W("row := db.DB().QueryRow(query, ")
	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(util.LowerFirstLetter(f.Camel)).Spc()
	}
	g.B.WL(")")

	g.B.Ln().WL2("var c sql.NullInt64")

	g.B.W("if err := row.Scan(&c); err != nil {")
	g.B.W("return 0, err")
	g.B.WL2("}")

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.W("return c.Int64, nil")

	g.B.WL2("}")
}

func (g *gen) genTxCountByIndex(args []*parse.F) {
	var m bytes.Buffer
	m.WriteString("TxCountBy")
	for _, f := range args {
		m.WriteString(f.Camel)
	}

	g.B.W("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(")

	for i, arg := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(util.LowerFirstLetter(arg.Camel)).Spc()

		if arg.T == cls.YamlTimestamp { // it is convenient to use integer when querying
			g.B.W(util.I64)
		} else {
			g.B.W(arg.GoT)
		}
	}
	g.B.W(")")
	g.B.Spc().W("(int64, error)").WL("{")

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")

	// make query sel
	g.B.W("query := `select count(1) from ", g.x.DB, ".", g.x.TB, " where ")
	for i, f := range args {
		if i > 0 {
			g.B.W(" and ")
		}
		g.B.W(f.Name, " = ?")
	}
	g.B.WL("`")
	// end make query sql

	g.B.W("row := Zotx.QueryRow(query, ")
	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(util.LowerFirstLetter(f.Camel)).Spc()
	}
	g.B.WL(")")

	g.B.Ln().WL2("var c sql.NullInt64")

	g.B.W("if err := row.Scan(&c); err != nil {")
	g.B.W("return 0, err")
	g.B.WL2("}")

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.W("return c.Int64, nil")

	g.B.WL2("}")
}

func (g *gen) genCountByRule() {
	var m bytes.Buffer
	m.WriteString("CountByRule")

	g.B.WL("func (mgr", " *_", g.T, "Mgr) ", m.String(), "(rules ...db.Rule) (int64, error) {")
	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), "`)")

	// make query sel
	g.B.WL(`var p []interface{}`)
	g.B.WL("query := `select count(1) from ", g.x.DB, ".", g.x.TB, " where `")
	g.B.WL(`for i, rule := range rules {`)
	g.B.WL(`	if i > 0 {`)
	g.B.WL(`		query += " and "`)
	g.B.WL(`	}`)
	g.B.WL(`	query += rule.S`)
	g.B.WL(`	if rule.P != nil {`)
	g.B.WL(`		p = append(p, rule.P)`)
	g.B.WL2(`	}`)
	g.B.WL2(`}`)

	// end make query sql

	g.B.WL2("row := db.DB().QueryRow(query, p...)")

	g.B.WL2("var c sql.NullInt64")

	g.B.WL("if err := row.Scan(&c); err != nil {")
	g.B.WL("return 0, err")
	g.B.WL2("}")

	g.B.WL("util.Log(`", g.x.DB, ".", g.x.TB, "`, `", m.String(), " ... done`)")

	g.B.WL("return c.Int64, nil")
	g.B.WL2("}")
}
