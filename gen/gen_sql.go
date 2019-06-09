package gen

import (
	"fmt"
	"github.com/envzo/zorm/cls"
	"github.com/envzo/zorm/parse"
	"github.com/envzo/zorm/util"
	"strconv"
	"strings"
)

var reservedWord = []string {
	"CURRENT_TIMESTAMP",
	"now()",
}

func isReserveWord(w string) bool {
	for _, s := range reservedWord {
		if strings.ToLower(w) == strings.ToLower(s) {
			return true
		}
	}
	return false
}

func genSql(x *parse.X) string {
	b := NewBuf()

	b.WL2("use ", x.DB, ";")
	b.W("create table if not exists `", x.TB, "` (").Ln()
	for i, v := range x.Fs {
		b.W("  `", v.Name, "` ")
		b.W(util.SqlTypeName(v.T))

		switch v.T {
		case cls.YamlStr:
			b.W("(", strconv.FormatInt(v.Size, 10), ")")
		}

		if !v.Nullable {
			b.W(" not null")
		}

		if v.AutoIncr {
			b.W(" auto_increment")
		}

		if v.Default != "" && v.Default != nil {
			if isReserveWord(fmt.Sprintf("%-v", v.Default)) {
				b.W(fmt.Sprintf(" default %-v", v.Default))
			} else {
				b.W(fmt.Sprintf(" default '%-v'", v.Default))
			}
		}

		if v.Comment != "" {
			b.W(" comment '" + v.Comment + "'")
		}

		if i != len(x.Fs)-1 {
			b.W(",").Ln()
		}
	}

	if x.PK != nil {
		b.WL(",").W("  primary key (`", x.PK.Name, "`)")
	}

	genIndex(b, x, true)
	genIndex(b, x, false)

	b.Ln().W(") engine=InnoDB default charset=utf8mb4")
	b.WL(" comment '", x.Comment, "';")

	return b.String()
}

func genIndex(b *Buf, x *parse.X, uniq bool) {
	indexes := x.Uniques
	if !uniq {
		indexes = x.Indexes
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
			b.W("_" + index.Name)
		}
		b.W("` (")

		for j, index := range v {
			if j != 0 {
				b.W(", ")
			}
			b.W("`" + index.Name + "`")
		}
		b.W(")")
		if i != len(indexes)-1 {
			b.W(",").Ln()
		}
	}
}
