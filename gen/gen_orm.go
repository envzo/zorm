package gen

import (
	"github.com/envzo/zorm/parse"
)

func (g *gen) genORM(pkg string) []byte {
	g.B.WL("// usage: ")
	g.B.WL2("// FindByXXX will not return sql.ErrNoRows, so it's caller's ability to check error")

	g.B.WL2("package ", pkg)

	g.B.WL(`import (`)
	g.B.WL(`"fmt"`)
	g.B.WL(`"database/sql"`)
	g.B.WL(`"github.com/envzo/zorm/db"`)
	g.B.WL2(`)`)

	g.B.WL2(`var _ = fmt.Printf`)
	g.B.WL2(`var _ = sql.ErrNoRows`)

	g.B.WL("type ", g.T, " struct {")
	for _, f := range g.D.Fields {
		g.B.Tab().W(ToCamel(f[0].Key.(string))).Spc().W(TypeName(f[0].Value.(string))).Ln()
	}
	g.B.WL("}")

	g.B.WL("func New", g.T, "() *", g.T, " {")
	g.B.WL("return &", g.T, "{}")
	g.B.WL("}")

	g.B.WL("type _", g.T, "Mgr struct {}")

	g.B.WL("var ", g.T, "Mgr = &_", g.T, "Mgr{}")

	fields := make([]*Field, len(g.D.Fields))
	for i, f := range g.D.Fields {
		n := f[0].Key.(string)
		t := f[0].Value.(string)
		fields[i] = &Field{
			Origin:  n,
			OriginT: t,
			Camel:   ToCamel(n),
			GoT:     TypeName(t),
		}
	}

	for _, fs := range g.D.Uniques {
		g.genIsExistsOne(fields, fs)
		g.genUniFindOne(fields, fs)
	}

	for _, fs := range g.D.Indexes {
		g.genFindByIndex(fields, fs)
		g.genCountByIndex(fields, fs)
	}

	g.genFindByJoin(fields)
	g.genCreate(fields).Ln()

	if g.D.PK != "" {
		g.genUpdateByPK(fields)
	}

	return g.B.Bytes()
}

type Field struct {
	Origin  string
	OriginT string
	Camel   string
	GoT     string
}

func (g *gen) genIsExistsOne(fields []*Field, args []string) {
	g.B.W("func (mgr", " *_", g.T, "Mgr) Is")
	for _, f := range args {
		g.B.W(ToCamel(f))
	}
	g.B.W("Exists(")

	for i, arg := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(LowerFirstLetter(ToCamel(arg))).Spc()

		// todo need refine: gather info in parse phase
		for _, f := range fields {
			if f.Origin == arg {
				if f.OriginT == parse.Timestamp { // it is convenient to use integer when querying
					g.B.W(I64)
				} else {
					g.B.W(f.GoT)
				}
				break
			}
		}
	}
	g.B.W(")")
	g.B.Spc().W("(bool, error)").WL("{")

	g.B.W("row := db.DB().QueryRow(`select count(1) from ", g.D.DB, ".", g.D.TB, " where ")

	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(f, "=?")
	}
	g.B.WL("`, ")

	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(LowerFirstLetter(ToCamel(f))).Spc()
	}
	g.B.WL2(")")

	g.B.WL("var c sql.NullInt64")
	g.B.Ln().W("if err := row.Scan(&c); err!= nil {")
	g.B.W("return false, err")
	g.B.WL("}")
	g.B.W("return c.Int64 > 0, nil")

	g.B.WL2("}")
}

func (g *gen) genUniFindOne(fields []*Field, args []string) {
	g.B.W("func (mgr", " *_", g.T, "Mgr) UniFindOneBy")
	for _, f := range args {
		g.B.W(ToCamel(f))
	}
	g.B.W("(")

	for i, arg := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(LowerFirstLetter(ToCamel(arg))).Spc()

		// todo need refine: gather info in parse phase
		for _, f := range fields {
			if f.Origin == arg {
				if f.OriginT == parse.Timestamp { // it is convenient to use integer when querying
					g.B.W(I64)
				} else {
					g.B.W(f.GoT)
				}
				break
			}
		}
	}
	g.B.W(")")
	g.B.Spc().W("(*" + g.T + ", error)").WL("{")

	g.B.W("row := db.DB().QueryRow(`select ")
	for i, f := range fields {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(f.Origin)
	}
	g.B.W(" from ", g.D.DB, ".", g.D.TB, " where ")

	for i, f := range args {
		if i > 0 {
			g.B.W("and ")
		}
		g.B.W(f, "=?")
	}
	g.B.WL("`, ")

	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(LowerFirstLetter(ToCamel(f))).Spc()
	}
	g.B.WL2(")")

	// temp variables
	vm := map[string]string{}

	for _, f := range fields {
		n := LowerFirstLetter(f.Camel)

		// todo need refine
		for _, arg := range args {
			if arg == f.Origin {
				n += "_1"
				break
			}
		}

		g.B.W("var ", n)
		vm[f.Camel] = n

		g.B.Spc().WL(TmpSqlType(f.OriginT))
	}

	g.B.Ln().W("if err := row.Scan(")
	for i, f := range fields {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W("&", vm[f.Camel])
	}
	g.B.W("); err!= nil {")
	g.B.W("return nil, err")
	g.B.WL2("}")

	g.B.WL("d := ", g.T, "{")
	for _, f := range fields {
		g.B.W(f.Camel, ":", vm[f.Camel])
		switch f.OriginT {
		case parse.I64, parse.Timestamp:
			g.B.W(".Int64")
		case parse.Str:
			g.B.W(".String")
		}

		g.B.WL(",")
	}
	g.B.WL2("}")

	g.B.W("return &d, nil")

	g.B.WL2("}")
}

func (g *gen) genCreate(fields []*Field) *B {
	g.B.WL("func (mgr *_", g.T, "Mgr) Create(d *", g.T, ") error {")
	if g.D.PK != "" {
		g.B.W("r, err")
	} else {
		g.B.W("_, err")
	}
	g.B.W(" := db.DB().Exec(`insert into ", g.D.DB, ".", g.D.TB, " (")
	cnt := 0
	for i, f := range fields {
		if f.Origin == g.D.PK {
			continue
		}
		g.B.W(f.Origin)
		if i != len(fields)-1 {
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

	for i, f := range fields {
		if f.Origin == g.D.PK {
			continue
		}
		g.B.W("d.", f.Camel)
		if i != len(fields)-1 {
			g.B.W(",")
		}
	}

	g.B.WL(")")
	g.B.WL("if err != nil {")
	g.B.WL("return err")
	g.B.WL("}")

	if g.D.PK != "" {
		g.B.WL("id,err := r.LastInsertId()")
		g.B.WL("if err!=nil {")
		g.B.WL("return err")
		g.B.WL("}")
		g.B.W("d.", ToCamel(g.D.PK), "=")

		for _, f := range fields { // check pk type
			if f.Origin == g.D.PK {
				if f.OriginT == parse.I64 {
					g.B.WL("id")
				} else if f.OriginT == parse.I32 {
					g.B.WL("int32(id)")
				}
				break
			}
		}
	}
	g.B.WL("return nil")
	return g.B.WL("}")
}

func (g *gen) genUpdateByPK(fields []*Field) {
	g.B.WL("func (mgr *_", g.T, "Mgr) Update(d *", g.T, ") (int64, error) {")
	g.B.W("r,err:=db.DB().Exec(`update ", g.D.DB, ".", g.D.TB, " set ")
	for i, f := range fields {
		if f.Origin == g.D.PK {
			continue
		}
		g.B.W(f.Origin, "=?")
		if i != len(fields)-1 {
			g.B.W(", ")
		}
	}
	g.B.WL(" where ", g.D.PK, "=?`, d.", ToCamel(g.D.PK), ")")
	g.B.WL("if err!=nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")
	g.B.WL("n,err:=r.RowsAffected()")
	g.B.WL("if err!=nil {")
	g.B.WL("return 0, err")
	g.B.WL("}")
	g.B.WL("return n, nil")
	g.B.WL("}")
}

func (g *gen) genFindByIndex(fields []*Field, args []string) {
	g.B.W("func (mgr", " *_", g.T, "Mgr) FindBy")
	for _, f := range args {
		g.B.W(ToCamel(f))
	}
	g.B.W("(")

	for i, arg := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(LowerFirstLetter(ToCamel(arg))).Spc()

		// todo need refine: gather info in parse phase
		for _, f := range fields {
			if f.Origin == arg {
				if f.OriginT == parse.Timestamp { // it is convenient to use integer when querying
					g.B.W(I64)
				} else {
					g.B.W(f.GoT)
				}
				break
			}
		}
	}
	g.B.W(", order []string, offset, limit int64)")
	g.B.Spc().W("([]*" + g.T + ", error)").WL("{")

	// make query sel
	g.B.W("query := `select ")
	for i, f := range fields {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(f.Origin)
	}
	g.B.W(" from ", g.D.DB, ".", g.D.TB, " where ")
	for i, f := range args {
		if i > 0 {
			g.B.W(" and ")
		}
		g.B.W(f, "=?")
	}
	g.B.WL("`")
	g.B.WL("for i, o := range order {")
	g.B.WL("if i==0 {")
	g.B.WL(`query += " order by "`)
	g.B.WL("} else {")
	g.B.WL(`query += ", "`)
	g.B.WL("}")
	g.B.WL("query += o[1:]")
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
		g.B.W(LowerFirstLetter(ToCamel(f))).Spc()
	}
	g.B.WL(")")
	g.B.WL("if err!=nil {")
	g.B.WL("return nil,err")
	g.B.WL2("}")

	// temp variables
	vm := map[string]string{}

	for _, f := range fields {
		n := LowerFirstLetter(f.Camel)

		// todo need refine
		for _, arg := range args {
			if arg == f.Origin {
				n += "_1"
				break
			}
		}

		g.B.W("var ", n)
		vm[f.Camel] = n

		g.B.Spc().WL(TmpSqlType(f.OriginT))
	}

	g.B.Ln().WL2("var ret []*", g.T)

	g.B.WL("for rows.Next(){")

	g.B.W("if err = rows.Scan(")
	for i, f := range fields {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W("&", vm[f.Camel])
	}
	g.B.W("); err!= nil {")
	g.B.W("return nil, err")
	g.B.WL2("}")

	g.B.WL("d := ", g.T, "{")
	for _, f := range fields {
		g.B.W(f.Camel, ":", vm[f.Camel])
		switch f.OriginT {
		case parse.I64, parse.Timestamp:
			g.B.W(".Int64")
		case parse.Str:
			g.B.W(".String")
		}

		g.B.WL(",")
	}
	g.B.WL("}")

	g.B.WL("ret = append(ret, &d)")

	g.B.WL("}") // end rows loop

	g.B.W("return ret, nil")

	g.B.WL2("}")
}

func (g *gen) genFindByJoin(fields []*Field) {
	g.B.W("func (mgr", " *_", g.T, "Mgr) FindByJoin(t string, on, where []db.Rule, order []string, offset, limit int64)")
	g.B.Spc().W("([]*" + g.T + ", error)").WL("{")

	g.B.WL2("var params []interface{}")

	// make query sel
	g.B.W("query := `select ")
	for i, f := range fields {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(f.Origin)
	}
	g.B.WL(" from ", g.D.DB, ".", g.D.TB, " join t on `")
	g.B.WL(`for i, v := range on {`)
	g.B.WL(`	if i > 0 {`)
	g.B.WL(`		query += " and "`)
	g.B.WL(`	}`)
	g.B.WL(`	query += v.S`)
	g.B.WL(`	if v.P != nil {`)
	g.B.WL(`		params = append(params, v.P)`)
	g.B.WL(`	}`)
	g.B.WL(`}`)
	g.B.WL(`for i, v := range where {`)
	g.B.WL(`	if i == 0 {`)
	g.B.WL(`		query += " where "`)
	g.B.WL(`	} else if i != len(where)-1 {`)
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
	g.B.WL("query += o[1:]")
	g.B.WL("if o[0] == '-' {")
	g.B.WL(`query += " desc"`)
	g.B.WL("}")
	g.B.WL("}")
	g.B.WL("if offset != -1 && limit != -1 {")
	g.B.WL(`query += fmt.Sprintf(" limit %d, %d", offset, limit)`)
	g.B.WL2("}")
	// end make query sql

	g.B.WL("rows, err := db.DB().Query(query, params...)")
	g.B.WL("if err!=nil {")
	g.B.WL("return nil,err")
	g.B.WL2("}")

	for _, f := range fields {
		n := LowerFirstLetter(f.Camel)

		g.B.W("var ", n)

		g.B.Spc().WL(TmpSqlType(f.OriginT))
	}

	g.B.Ln().WL2("var ret []*", g.T)

	g.B.WL("for rows.Next(){")

	g.B.W("if err = rows.Scan(")
	for i, f := range fields {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W("&", LowerFirstLetter(f.Camel))
	}
	g.B.W("); err!= nil {")
	g.B.W("return nil, err")
	g.B.WL2("}")

	g.B.WL("d := ", g.T, "{")
	for _, f := range fields {
		g.B.W(f.Camel, ":", LowerFirstLetter(f.Camel))
		switch f.OriginT {
		case parse.I64, parse.Timestamp:
			g.B.W(".Int64")
		case parse.Str:
			g.B.W(".String")
		}

		g.B.WL(",")
	}
	g.B.WL("}")

	g.B.WL("ret = append(ret, &d)")

	g.B.WL("}") // end rows loop

	g.B.W("return ret, nil")

	g.B.WL2("}")
}

func (g *gen) genCountByIndex(fields []*Field, args []string) {
	g.B.W("func (mgr", " *_", g.T, "Mgr) CountBy")
	for _, f := range args {
		g.B.W(ToCamel(f))
	}
	g.B.W("(")

	for i, arg := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(LowerFirstLetter(ToCamel(arg))).Spc()

		// todo need refine: gather info in parse phase
		for _, f := range fields {
			if f.Origin == arg {
				if f.OriginT == parse.Timestamp { // it is convenient to use integer when querying
					g.B.W(I64)
				} else {
					g.B.W(f.GoT)
				}
				break
			}
		}
	}
	g.B.W(")")
	g.B.Spc().W("(int64, error)").WL("{")

	// make query sel
	g.B.W("query := `select count(1) from ", g.D.DB, ".", g.D.TB, " where ")
	for i, f := range args {
		if i > 0 {
			g.B.W(" and ")
		}
		g.B.W(f, "=?")
	}
	g.B.WL("`")
	// end make query sql

	g.B.W("row := db.DB().QueryRow(query, ")
	for i, f := range args {
		if i > 0 {
			g.B.W(", ")
		}
		g.B.W(LowerFirstLetter(ToCamel(f))).Spc()
	}
	g.B.WL(")")

	g.B.Ln().WL2("var c sql.NullInt64")

	g.B.W("if err := row.Scan(&c); err != nil {")
	g.B.W("return 0, err")
	g.B.WL2("}")

	g.B.W("return c.Int64, nil")

	g.B.WL2("}")
}
