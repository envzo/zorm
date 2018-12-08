package parse

import "gopkg.in/yaml.v2"

// filed types
const (
	I32       = "i32"
	I64       = "i64"
	Str       = "str"
	Timestamp = "timestamp"
)

// field attributes
const (
	AutoIncr = "__auto_incr"
	Size     = "__size"
	Comment  = "__comment"
	Nullable = "__nullable"
)

var (
	dbs = map[string]bool{"mysql": true}
	ts  = map[string]bool{
		I32:       true,
		I64:       true,
		Str:       true,
		Timestamp: true,
	}
)

type Def struct {
	Engine  string          `yaml:"engine"`
	DB      string          `yaml:"db"`
	TB      string          `yaml:"tb"`
	Comment string          `yaml:"comment"`
	Fields  []yaml.MapSlice `yaml:"fields"`
	PK      string          `yaml:"pk"`
	Indexes [][]string      `yaml:"indexes"`
	Uniques [][]string      `yaml:"uniques"`
}
