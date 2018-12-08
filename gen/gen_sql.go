package gen

import (
	"strconv"

	"github.com/envzo/zorm/parse"
)

func genSql(d *parse.Def) string {
	b := NewB()

	b.WL("use ", d.DB, ";")
	b.W("create table if not exists `", d.TB, "` (").Ln()
	for i, v := range d.Fields {
		fn := v[0].Key.(string)
		t := v[0].Value.(string)

		b.W("  `", fn, "` ")
		b.W(SqlTypeName(t))

		switch t {
		case parse.Str:
			b.W("(")
			for _, f := range v {
				if f.Key.(string) == parse.Size {
					b.W(strconv.FormatInt(int64(f.Value.(int)), 10))
				}
			}
			b.W(")")
		}

		for _, f := range v {
			if f.Key.(string) == parse.Nullable {
				if !f.Value.(bool) {
					b.W(" not null")
				}
			} else if f.Key.(string) == parse.AutoIncr {
				b.W(" auto_increment")
			} else if f.Key.(string) == parse.Comment {
				b.W(" comment '" + f.Value.(string) + "'")
			}
		}

		if i != len(d.Fields)-1 {
			b.W(",").Ln()
		}
	}

	if d.PK != "" {
		b.WL(",").W("  primary key (`", d.PK, "`)")
	}

	appendIndex(b, d, true)
	appendIndex(b, d, false)

	b.Ln().W(") engine=InnoDB default charset=utf8mb4")
	b.WL(" comment '", d.Comment, "';")

	return b.String()
}

func appendIndex(b *B, d *parse.Def, uniq bool) {
	indexes := d.Uniques
	if !uniq {
		indexes = d.Indexes
	}

	for i, v := range indexes {
		if i == 0 {
			b.W(",").Ln()
		}

		if uniq {
			b.W("  unique key `uni")
		} else {
			b.W("  index `idx")
		}

		for _, index := range v {
			b.W("_" + index)
		}
		b.W("` (")

		for j, index := range v {
			if j != 0 {
				b.W(", ")
			}
			b.W("`" + index + "`")
		}
		b.W(")")
		if i != len(indexes)-1 {
			b.W(",").Ln()
		}
	}
}
