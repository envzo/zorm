package gen

import (
	"fmt"
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
	g.B.WL(`"github.com/envzo/zorm/db"`)
	g.B.WL(`"github.com/envzo/zorm/util"`)
	g.B.WL2(`)`)

	g.B.WL(`var _ = errors.New`)
	g.B.WL(`var _ = fmt.Printf`)
	g.B.WL(`var _ = strings.Trim`)
	g.B.WL(`var _ = sql.ErrNoRows`)
	g.B.WL(`var _ = util.I64`)
	g.B.WL(`var _ = time.Nanosecond`)

	g.B.WL("type ", g.T, " struct {")
	for _, f := range g.x.Fs {
		g.B.Tab().WL("// ", f.Comment)
		g.B.Tab().W(f.Camel).Spc().W(f.GoT).Ln()
	}
	g.B.Ln().WL("// 调用Upsert方法时，baby为true则insert，反之update")
	g.B.WL("baby bool")
	g.B.WL("}")

	g.B.WL("func New", g.T, "() *", g.T, " {")
	g.B.WL("return &", g.T, "{baby: true}")
	g.B.WL("}")

	g.B.WL("type _", g.T, "Mgr struct {}")

	g.B.Ln()

	for _, fs := range g.x.Indexes {
		g.B.W("type ", g.getIxEntityTypeName(fs))
		g.B.WL(" struct {")
		for _, f := range fs {
			g.B.Tab().WL(f.Camel, " ", f.GoT)
		}
		g.B.WL("}")

		g.B.Ln()
	}

	g.B.WL("var ", g.T, "Mgr = &_", g.T, "Mgr{}")

	for _, fs := range g.x.Uniques {
		g.genIsExists(fs)
		g.genUniFind(fs)
		g.genUniUpdate(fs)
		g.genUniRm(fs)
	}

	for _, fs := range g.x.Indexes {
		g.genFindByIndex(fs)
		g.genFindByIndexArray(fs)
		g.genCountByIndex(fs)
	}

	g.genFindByMultiJoin()
	g.genCountByMultiJoin()
	g.genFindByJoin()
	g.genFindByCond()
	g.genFindAllByCond()
	g.genCreate().Ln()
	g.genUpsert().Ln()
	g.genCountByRule()
	g.genRmByRule()

	if g.x.PK != nil {
		g.genUniFindByPk()
		g.genUpdateByPK().Ln()
		g.genRmByPK()
		g.genIsExistsByPK()
	}

	return g.B.Bytes()
}

func (g *gen) genIsExists(args []*parse.F) {
	g.B.W("func (mgr", " *_", g.T, "Mgr) Is")
	for _, f := range args {
		g.B.W(f.Camel)
	}
	g.B.W("Exists(")

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
	g.B.W("return c.Int64 > 0, nil")

	g.B.WL2("}")
}

func (g *gen) genIsExistsByPK() {
	g.B.WL("func (mgr", " *_", g.T, "Mgr) IsExistsByPK(pk ", g.x.PK.GoT, ") (bool, error) {")
	g.B.WL("row := db.DB().QueryRow(`select count(1) from ", g.x.DB, ".", g.x.TB, " where ", g.x.PK.Name, " = ?`, pk)")
	g.B.WL("var c sql.NullInt64")
	g.B.Ln().WL("if err := row.Scan(&c); err!= nil {")
	g.B.WL("return false, err")
	g.B.WL("}")
	g.B.WL("return c.Int64 > 0, nil")
	g.B.WL2("}")
}

func (g *gen) genUniFind(args []*parse.F) {
	g.B.W("func (mgr", " *_", g.T, "Mgr) UniFindBy")
	for _, f := range args {
		g.B.W(f.Camel)
	}
	g.B.W("(")

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

	g.B.W("return &d, nil")

	g.B.WL2("}")
}

func (g *gen) genUniUpdate(args []*parse.F) {
	g.B.W("func (mgr", " *_", g.T, "Mgr) UpdateBy")
	for _, f := range args {
		g.B.W(f.Camel)
	}
	g.B.W("(d *", g.T, ")")
	g.B.Spc().W("(int64, error) ").WL("{")

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
	g.B.WL("return n, nil")
	g.B.WL("}").Ln()
}

func (g *gen) genUniRm(args []*parse.F) {
	g.B.W("func (mgr", " *_", g.T, "Mgr) UniRmBy")
	for _, f := range args {
		g.B.W(f.Camel)
	}
	g.B.W("(d *", g.T, ")")
	g.B.Spc().W("(int64, error) ").WL("{")

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
	g.B.WL("return n, nil")
	g.B.WL("}").Ln()
}

func (g *gen) genCreate() *Buf {
	g.B.WL("func (mgr *_", g.T, "Mgr) Create(d *", g.T, ") error {")
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
	g.B.W("func (mgr", " *_", g.T, "Mgr) UniFindByPK(", util.LowerFirstLetter(g.x.PK.Camel))

	if g.x.PK.T == cls.YamlTimestamp { // it is convenient to use integer when querying
		g.B.Spc().W(util.I64)
	} else {
		g.B.Spc().W(g.x.PK.GoT)
	}
	g.B.W(")")
	g.B.Spc().W("(*" + g.T + ", error)").WL("{")

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

	g.B.W("return &d, nil")

	g.B.WL2("}")
}

func (g *gen) genUpdateByPK() *Buf {
	g.B.WL("func (mgr *_", g.T, "Mgr) Update(d *", g.T, ") (int64, error) {")
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
	g.B.WL("return n, nil")
	return g.B.WL("}")
}

func (g *gen) genRmByPK() {
	g.B.WL("func (mgr *_", g.T, "Mgr) RmByPK(pk ", g.x.PK.GoT, ") (int64, error) {")
	g.B.WL(`query := "delete from `, g.x.DB, ".", g.x.TB, " where ", g.x.PK.Name, ` = ?"`)

	g.B.WL("r,err := db.DB().Exec(query, pk)")
	g.B.WL("if err != nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")
	g.B.WL("n,err := r.RowsAffected()")
	g.B.WL("if err != nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")
	g.B.WL("return n, nil")
	g.B.WL("}")
}

func (g *gen) genRmByRule() {
	g.B.WL("func (mgr *_", g.T, "Mgr) RmByRule(rules ...db.Rule) (int64, error) {")
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
	g.B.WL("return n, nil")
	g.B.WL("}")
}

func (g *gen) genFindByIndex(args []*parse.F) {
	g.B.W("func (mgr", " *_", g.T, "Mgr) FindBy")
	for _, f := range args {
		g.B.W(f.Camel)
	}
	g.B.W("(")

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
	g.B.WL("if err!=nil {")
	g.B.WL("return nil,err")
	g.B.WL2("}")

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

	g.B.W("return ret, nil")

	g.B.WL2("}")
}

func (g *gen) genFindByIndexArray(args []*parse.F) {
	g.B.WL("// 通过索引数组查询")
	g.B.W("func (mgr", " *_", g.T, "Mgr) FindBy")
	for _, f := range args {
		g.B.W(f.Camel)
	}
	g.B.W("Array(", "entities []*", g.getIxEntityTypeName(args))

	g.B.W(", order []string, offset, limit int64)")
	g.B.Spc().W("([]*" + g.T + ", string, error)").WL("{")

	g.B.WL("if len(entities) == 0 {")
	g.B.WL("return nil, \"\", errors.New(\"input entities empty. \")")
	g.B.WL2("}")

	// 构造问号占位符
	a := "("
	for i := range args {
		a += "?"
		if i != len(args)-1 {
			a += ","
		} else {
			a += ")"
		}
	}

	if len(args) == 1 {
		a = "?"
	}

	b := "," + a
	g.B.WL("str := \"", a, "\" + strings.Repeat(\"", b, "\", len(entities) - 1)")

	// make query sel
	g.B.W("query := fmt.Sprintf(\"select ")
	for i, f := range g.x.Fs {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W("`", f.Name, "`")
	}
	g.B.W(" from ", g.x.DB, ".", g.x.TB, " where (")
	for i, f := range args {
		g.B.W("`", f.Name, "`")
		if i != len(args)-1 {
			g.B.W(", ")
		}
	}
	g.B.WL(") in (%s)\", str)")

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

	g.B.WL("params := make([]interface{}, 0, ", fmt.Sprintf("%d", len(args)), " * len(entities))")
	g.B.WL("for _, entity := range entities {")
	for _, f := range args {
		g.B.WL("params = append(params, entity.", f.Camel, ")")
	}
	g.B.WL("}")

	g.B.WL("rows, err := db.DB().Query(query, params...)")
	g.B.WL("if err!=nil {")
	g.B.WL("return nil, query, err")
	g.B.WL2("}")

	// temp variables
	vm := map[string]string{}

	for _, f := range g.x.Fs {
		n := util.LowerFirstLetter(f.Camel)

		if n == "type" {
			n += "_"
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
	g.B.W("return nil, query, err")
	g.B.WL2("}")

	g.B.WL("d := ", g.T, "{}")
	for _, f := range g.x.Fs {
		g.B.W("d.", f.Camel, "=", util.DerefNilSqlType(vm[f.Camel], f.T)).Ln()
	}

	g.B.WL("ret = append(ret, &d)")

	g.B.WL("}") // end rows loop

	g.B.W("return ret, query, nil")

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
	g.B.W("func (mgr", " *_", g.T, "Mgr) FindByMultiJoin(joins []db.Join, where []db.Rule, order []string, offset, limit int64)")
	g.B.Spc().W("([]*" + g.T + ", error)").WL("{")

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
	g.B.WL("return nil,err")
	g.B.WL2("}")

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

	g.B.W("return ret, nil")

	g.B.WL2("}")
}

func (g *gen) genCountByMultiJoin() {
	g.B.W("func (mgr", " *_", g.T, "Mgr) CountByMultiJoin(joins []db.Join, where []db.Rule) ")
	g.B.W("(int64, error)").WL("{")

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

	g.B.W("return c.Int64, nil")

	g.B.WL2("}")
}

func (g *gen) genFindByCond() {
	g.B.W("func (mgr", " *_", g.T, "Mgr) FindByCond(where []db.Rule, order []string, offset, limit int64)")
	g.B.Spc().W("([]*" + g.T + ", error)").WL("{")

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
	g.B.WL("return nil,err")
	g.B.WL2("}")

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
	g.B.W("func (mgr", " *_", g.T, "Mgr) FindAllByCond(where []db.Rule, order []string)")
	g.B.Spc().W("([]*" + g.T + ", error)").WL("{")

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
	g.B.WL("return nil,err")
	g.B.WL2("}")

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

func (g *gen) genCountByIndex(args []*parse.F) {
	g.B.W("func (mgr", " *_", g.T, "Mgr) CountBy")
	for _, f := range args {
		g.B.W(f.Camel)
	}
	g.B.W("(")

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

	g.B.W("return c.Int64, nil")

	g.B.WL2("}")
}

func (g *gen) genCountByRule() {
	g.B.WL("func (mgr", " *_", g.T, "Mgr) CountByRule(rules ...db.Rule) (int64, error) {")

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
	g.B.WL("return c.Int64, nil")
	g.B.WL2("}")
}

func (g *gen) getIxEntityTypeName(fs []*parse.F) string {
	s := "IxEntity" + g.T

	for _, f := range fs {
		s += f.Camel
	}

	return s
}
