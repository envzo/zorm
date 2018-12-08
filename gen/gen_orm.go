package gen

import (
	"github.com/envzo/zorm/parse"
)

var bt string // bean type

func genORM(pkg string, d *parse.Def) []byte {
	b := NewB()
	b.W("package ", pkg).Ln2()

	b.WL(`import (`)
	b.WL(`"database/sql"`)
	b.WL2(`"time"`)
	b.WL(`"github.com/envzo/zorm/db"`)
	b.WL2(`)`)

	b.W("var _ = time.Time{}").Ln2()

	bt = ToCamel(d.TB)

	b.WL("type ", bt, " struct {")
	for _, f := range d.Fields {
		b.Tab().W(ToCamel(f[0].Key.(string))).Spc().W(TypeName(f[0].Value.(string))).Ln()
	}
	b.WL("}")

	b.WL("func New", bt, "() *", bt, " {")
	b.WL("return &", bt, "{}")
	b.WL("}")

	tmgr := bt + "Mgr"

	b.WL("type _", tmgr, " struct {}")

	b.WL("var ", tmgr, " = &_", tmgr, "{}")

	fields := make([]*Field, len(d.Fields))
	for i, f := range d.Fields {
		n := f[0].Key.(string)
		t := f[0].Value.(string)
		fields[i] = &Field{
			Origin:  n,
			OriginT: t,
			Camel:   ToCamel(n),
			GoT:     TypeName(t),
		}
	}

	for _, fs := range d.Uniques {
		genQueryOne(fields, d.DB, d.TB, b, fs)
	}

	genCreate(fields, d.DB, d.TB, d.PK, b).Ln()
	genUpdate(fields, d.DB, d.TB, d.PK, b)

	return b.Bytes()
}

type Field struct {
	Origin  string
	OriginT string
	Camel   string
	GoT     string
}

func genQueryOne(fields []*Field, db, tb string, b *B, args []string) {
	b.W("func (mgr", " *_", bt, "Mgr) FindOneBy")
	for _, f := range args {
		b.W(ToCamel(f))
	}
	b.W("(")

	for i, arg := range args {
		if i > 0 {
			b.W(", ")
		}
		b.W(LowerFirstLetter(ToCamel(arg))).Spc()

		// todo need refine: gather info in parse phase
		for _, f := range fields {
			if f.Origin == arg {
				if f.OriginT == parse.Timestamp { // it is convenient to use integer when querying
					b.W(I64)
				} else {
					b.W(f.GoT)
				}
				break
			}
		}
	}
	b.W(")")
	b.Spc().W("(*" + bt + ", error)").WL("{")

	b.W("row := db.DB().QueryRow(`select ")
	for i, f := range fields {
		if i > 0 {
			b.W(", ")
		}
		b.W(f.Origin)
	}
	b.W(" from ", db, ".", tb, " where ")

	for i, f := range args {
		if i > 0 {
			b.W(", ")
		}
		b.W(f, "=?")
	}
	b.WL("`, ")

	for i, f := range args {
		if i > 0 {
			b.W(", ")
		}
		b.W(LowerFirstLetter(ToCamel(f))).Spc()
	}
	b.WL2(")")

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

		b.W("var ", n)
		vm[f.Camel] = n

		b.Spc().WL(TmpSqlType(f.OriginT))
	}

	b.Ln().W("if err := row.Scan(")
	for i, f := range fields {
		if i > 0 {
			b.W(", ")
		}
		b.W("&", vm[f.Camel])
	}
	b.W("); err!= nil {")
	b.W("return nil, err")
	b.WL2("}")

	b.WL("data := ", bt, "{")
	for _, f := range fields {
		if f.OriginT == parse.Timestamp {
			continue
		}

		b.W(f.Camel, ":", vm[f.Camel])
		switch f.OriginT {
		case parse.I64:
			b.W(".Int64")
		case parse.Str:
			b.W(".String")
		}

		b.WL(",")
	}
	b.WL2("}")

	// convert timestamp to time
	ok := true
	for _, f := range fields {
		if f.OriginT != parse.Timestamp {
			continue
		}
		b.W("t")
		if ok {
			b.W(":")
			ok = false
		}
		b.WL("=time.Unix(", vm[f.Camel], ".Int64, 0)")
		b.WL("data.", f.Camel, "=&t")
	}

	b.W("return &data, nil")

	b.WL2("}")
}

func genCreate(fields []*Field, db, tb, pk string, b *B) *B {
	b.WL("func (mgr *_", bt, "Mgr) Create(d *", bt, ") error {")
	b.W("r,err:=db.DB().Exec(`insert into ", db, ".", tb, " (")
	cnt := 0
	for i, f := range fields {
		if f.Origin == pk {
			continue
		}
		b.W(f.Origin)
		if i != len(fields)-1 {
			b.W(", ")
		}
		cnt++
	}
	b.W(") value (")
	for i := 0; i < cnt; i++ {
		if i > 0 {
			b.W(",")
		}
		b.W("?")
	}
	b.W(")`,")

	for i, f := range fields {
		if f.Origin == pk {
			continue
		}
		b.W("d.", f.Camel)
		if i != len(fields)-1 {
			b.W(",")
		}
	}

	b.WL(")")
	b.WL("if err!=nil {")
	b.WL("return err")
	b.WL("}")
	if pk != "" {
		b.WL("id,err:=r.LastInsertId()")
		b.WL("if err!=nil {")
		b.WL("return err")
		b.WL("}")
		b.W("d.", ToCamel(pk), "=")
		// check pk type
		for _, f := range fields {
			if f.Origin == pk {
				if f.OriginT == parse.I64 {
					b.WL("id")
				} else if f.OriginT == parse.I32 {
					b.WL("int32(id)")
				}
				break
			}
		}
	}
	b.WL("return nil")
	return b.WL("}")
}

func genUpdate(fields []*Field, db, tb, pk string, b *B) {
	b.WL("func (mgr *_", bt, "Mgr) Update(d *", bt, ") (int64, error) {")
	b.W("r,err:=db.DB().Exec(`update ", db, ".", tb, " set ")
	for i, f := range fields {
		if f.Origin == pk {
			continue
		}
		b.W(f.Origin, "=?")
		if i != len(fields)-1 {
			b.W(", ")
		}
	}
	b.WL(" where ", pk, "=?`, d.", ToCamel(pk), ")")
	b.WL("if err!=nil {")
	b.WL("return 0, err")
	b.WL("}")
	b.WL("n,err:=r.RowsAffected()")
	b.WL("if err!=nil {")
	b.WL("return 0, err")
	b.WL("}")
	b.WL("return n, nil")
	b.WL("}")
}
