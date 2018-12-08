package gen

import (
	"strings"
	"unicode"

	"github.com/envzo/zorm/parse"
)

const (
	I64 = "int64"
)

func TypeName(in string) string {
	switch in {
	case parse.I64:
		return I64
	case parse.I32:
		return "int32"
	case parse.Str:
		return "string"
	case parse.Timestamp:
		return "*time.Time"
	}
	return in
}

func TmpSqlType(in string) string {
	switch in {
	case parse.I64, parse.Timestamp:
		return "sql.NullInt64"
	case parse.I32:
		return "int32"
	case parse.Str:
		return "sql.NullString"
	}
	return in
}

func SqlTypeName(in string) string {
	switch in {
	case parse.I64, parse.Timestamp:
		return "bigint"
	case parse.I32:
		return "int"
	case parse.Str:
		return "varchar"
	}
	return in
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

func UpperFirstLetter(in string) string {
	return strings.ToUpper(in[:1]) + in[1:]
}

func LowerFirstLetter(in string) string {
	return strings.ToLower(in[:1]) + in[1:]
}
