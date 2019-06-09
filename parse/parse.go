package parse

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/envzo/zorm/cls"
	"github.com/envzo/zorm/util"
)

func Parse(b []byte) ([]*X, error) {
	m := map[string]*Def{}
	if err := yaml.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	r := make([]*X, 0, len(m))

	for _, def := range m {
		if !DBs[def.Engine] {
			return nil, errors.New("invalid engine: " + def.Engine)
		}

		if def.TB == "" {
			return nil, errors.New("must have tb")
		}

		if len(def.Fields) == 0 {
			return nil, errors.New("must have at least one field")
		}

		x, err := newParser(def).parse()
		if err != nil {
			return nil, err
		}

		r = append(r, x)
	}

	return r, nil
}

type parser struct {
	d *Def
}

func newParser(d *Def) *parser {
	return &parser{d: d}
}

func (p *parser) parse() (*X, error) {
	x := &X{
		Engine:  p.d.Engine,
		DB:      p.d.DB,
		TB:      p.d.TB,
		Fs:      []*F{},
		Comment: p.d.Comment,
	}

	pk, err := p.parsePK()
	if err != nil {
		return nil, err
	}
	x.PK = pk
	for _, s := range p.d.Fields {
		f, err := p.parseField(s)
		if err != nil {
			return nil, err
		}
		x.Fs = append(x.Fs, f)
	}

	idx, err := p.parseIndex(p.d.Indexes)
	if err != nil {
		return nil, err
	}
	x.Indexes = idx

	if idx, err = p.parseIndex(p.d.Uniques); err != nil {
		return nil, err
	}
	x.Uniques = idx

	return x, nil
}

func (p *parser) parsePK() (*F, error) {
	if p.d.PK == "" {
		return nil, nil
	}

	// check if exists

	var f yaml.MapSlice
	for _, s := range p.d.Fields {
		if s[0].Key.(string) == p.d.PK {
			f = s
			break
		}
	}
	if f == nil {
		return nil, errors.New("primary key not found: " + p.d.PK)
	}

	pk, err := p.parseField(f)
	if err != nil {
		return nil, err
	}

	if pk.T != cls.YamlI32 && pk.T != cls.YamlI64 {
		return nil, errors.New("primary key must be integer")
	}

	return pk, nil
}

func (p *parser) parseField(s yaml.MapSlice) (*F, error) {
	n := s[0].Key.(string)
	if strings.HasPrefix(n, "_") || strings.HasSuffix(n, "_") || strings.Contains(n, "-") {
		return nil, errors.New("invalid field name: " + n)
	}

	t := s[0].Value.(string)
	if !TS[t] {
		return nil, errors.New("invalid field type: " + t)
	}

	f := F{
		Name:  n,
		T:     t,
		Camel: util.ToCamel(n),
		GoT:   util.GoT(t),
	}

	// check attributes
	hasAttr := map[string]bool{}
	for i, attr := range s {
		if i == 0 {
			continue
		}

		attrName := attr.Key.(string)
		switch attrName {
		case AutoIncr:
			if t != cls.YamlI32 && t != cls.YamlI64 {
				return nil, errors.New("field has auto increment attr must be integer: " + f.Name)
			}
			if n != p.d.PK {
				return nil, errors.New("auto incremented field must be primary key")
			}
			f.AutoIncr = attr.Value.(bool)
		case Size:
			if t != cls.YamlStr {
				return nil, errors.New("field has size attr must be string: " + f.Name)
			}

			size, ok := attr.Value.(int)
			if !ok {
				return nil, errors.New("value of size attr must be integer: " + f.Name)
			}
			f.Size = int64(size)
		case Comment:
			f.Comment = attr.Value.(string)
		case Nullable:
			f.Nullable = attr.Value.(bool)
		case Default:
			f.Default = attr.Value.(string)
		default:
			return nil, errors.New("unknown field: " + attrName)
		}
		hasAttr[attrName] = true
	}

	if t == cls.YamlStr && !hasAttr[Size] {
		return nil, errors.New("field with string type must have size attr")
	}

	return &f, nil
}

func (p *parser) parseIndex(indexes [][]string) ([][]*F, error) {
	var a [][]*F

	for _, idx := range indexes {
		var b []*F

	Field:
		for _, field := range idx {
			for _, f := range p.d.Fields {
				if f[0].Key.(string) == field {
					t := f[0].Value.(string)
					b = append(b, &F{
						Name:  field,
						T:     t,
						Camel: util.ToCamel(field),
						GoT:   util.GoT(t),
					})
					continue Field
				}
			}
			return nil, errors.New("index field not found: " + field)
		}

		a = append(a, b)
	}
	return a, nil
}
