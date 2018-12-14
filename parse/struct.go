package parse

import (
	"gopkg.in/yaml.v2"

	"github.com/envzo/zorm/util"
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
	Ts  = map[string]bool{
		util.YamlI32:       true,
		util.YamlI64:       true,
		util.YamlStr:       true,
		util.YamlTimestamp: true,
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
