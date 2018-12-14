package util

import (
	"strings"
	"unicode"
)

const (
	I64 = "int64"
	I32 = "int32"
	Str = "string"

	// filed types
	YamlI32       = "i32"
	YamlI64       = "i64"
	YamlStr       = "str"
	YamlTimestamp = "timestamp"
)

func GoT(in string) string {
	switch in {
	case YamlI64, YamlTimestamp:
		return I64
	case YamlI32:
		return I32
	case YamlStr:
		return Str
	}
	return in
}

func NilSqlType(in string) string {
	switch in {
	case YamlI64, YamlTimestamp:
		return "sql.NullInt64"
	case YamlI32:
		return I32
	case YamlStr:
		return "sql.NullString"
	}
	return in
}

func SqlTypeName(in string) string {
	switch in {
	case YamlI64, YamlTimestamp:
		return "bigint"
	case YamlI32:
		return "int"
	case YamlStr:
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
