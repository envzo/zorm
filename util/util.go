package util

import (
	"strings"
	"unicode"

	"github.com/envzo/zorm/cls"
)

const (
	I64     = "int64"
	I32     = "int32"
	Str     = "string"
	Bool    = "bool"
	Float32 = "float32"
	Time    = "*time.Time"
)

func GoT(in string) string {
	switch in {
	case cls.YamlI64, cls.YamlTimestamp:
		return I64
	case cls.YamlI32:
		return I32
	case cls.YamlStr:
		return Str
	case cls.YamlBool:
		return Bool
	case cls.YamlFloat:
		return Float32
	case cls.YamlDate:
		return Time
	default:
		panic("unknown type: " + in)
	}
}

func NilSqlType(in string) string {
	switch in {
	case cls.YamlI32, cls.YamlI64, cls.YamlTimestamp:
		return "sql.NullInt64"
	case cls.YamlStr, cls.YamlDate:
		return "sql.NullString"
	case cls.YamlBool:
		return "sql.NullBool"
	case cls.YamlFloat:
		return "sql.NullFloat64"
	default:
		panic("unknown type: " + in)
	}
}

func DerefNilSqlType(n, t string) string {
	switch t {
	case cls.YamlI64, cls.YamlTimestamp:
		return n + ".Int64"
	case cls.YamlI32:
		return "int32(" + n + ".Int64)"
	case cls.YamlFloat:
		return "float32(" + n + ".Float64)"
	case cls.YamlStr:
		return n + ".String"
	case cls.YamlDate:
		if n == "" {
			return "nil"
		}
		return `util.MustParseDateStr(` + n + `.String)`
	case cls.YamlBool:
		return n + ".Bool"
	default:
		panic("unknown type: " + t)
	}
}

func SqlTypeName(in string) string {
	switch in {
	case cls.YamlI64, cls.YamlTimestamp:
		return "bigint"
	case cls.YamlI32:
		return "int"
	case cls.YamlStr:
		return "varchar"
	case cls.YamlBool:
		return "tinyint(1)"
	case cls.YamlFloat:
		return "float"
	case cls.YamlDate:
		return "date"
	default:
		panic("unknown type: " + in)
	}
}

func ToCamel(in string) string {
	runes := []rune(in)
	var out []rune

	for i, r := range runes {
		if r == '_' {
			continue
		}
		if i == 0 || runes[i-1] == '_' {
			out = append(out, unicode.ToUpper(r))
			continue
		}
		out = append(out, r)
	}

	return string(out)
}

func LowerFirstLetter(in string) string {
	return strings.ToLower(in[:1]) + in[1:]
}
