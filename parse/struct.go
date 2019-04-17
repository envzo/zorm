package parse

import (
	"github.com/envzo/zorm/cls"
	"gopkg.in/yaml.v2"
)

// field attributes
const (
	AutoIncr = "__auto_incr"
	Size     = "__size"
	Comment  = "__comment"
	Nullable = "__nullable"
)

var (
	DBs = map[string]bool{"mysql": true}
	TS  = map[string]bool{
		cls.YamlI32:       true,
		cls.YamlI64:       true,
		cls.YamlStr:       true,
		cls.YamlTimestamp: true,
		cls.YamlFloat:     true,
		cls.YamlDouble:    true,
		cls.YamlBool:      true,
		cls.YamlDate:      true,
		cls.YamlDateTime:  true,
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

type F struct {
	Name     string
	T        string
	Size     int64
	Nullable bool
	AutoIncr bool
	Comment  string

	Camel string
	GoT   string
}

type X struct {
	Engine  string
	DB      string
	TB      string
	PK      *F
	Fs      []*F
	Comment string
	Uniques [][]*F
	Indexes [][]*F
}
