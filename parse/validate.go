package parse

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v2"
)

func check(d *Def) error {
	if !dbs[d.Engine] {
		return errors.New("invalid db: " + d.Engine)
	}
	if len(d.Fields) == 0 {
		return errors.New("must have at least one field")
	}
	for _, f := range d.Fields {
		if err := checkField(f, d.PK); err != nil {
			return err
		}
	}

	if d.PK != "" {
		if err := checkPK(d.PK, d.Fields); err != nil {
			return err
		}
	}

	if err := checkIndex(d.Fields, d.Indexes); err != nil {
		return err
	}
	if err := checkIndex(d.Fields, d.Uniques); err != nil {
		return err
	}
	return nil
}

func checkPK(pk string, fields []yaml.MapSlice) error {
	// check if exists
	var f *yaml.MapItem
	for _, idx := range fields {
		if idx[0].Key.(string) == pk {
			f = &idx[0]
			break
		}
	}
	if f == nil {
		return errors.New("primary key not found: " + pk)
	}

	if t := f.Value.(string); t != I32 && t != I64 {
		return errors.New("primary key must be integer")
	}
	return nil
}

func checkField(m yaml.MapSlice, pk string) error {
	fn := m[0].Key.(string)
	if strings.HasPrefix(fn, "_") || strings.HasSuffix(fn, "_") || strings.Contains(fn, "-") {
		return errors.New("invalid field name: " + fn)
	}

	t := m[0].Value.(string)
	if !ts[t] {
		return errors.New("invalid field type: " + t)
	}

	// check attributes
	hasSize := false
	for i, v := range m {
		if i == 0 {
			continue
		}

		switch n := v.Key.(string); n {
		case AutoIncr:
			if t != I32 && t != I64 {
				return errors.New("auto incremented field must be integer")
			}
			// must be primary key too
			if fn != pk {
				return errors.New("auto incremented field must be primary key")
			}
		case Size:
			if t != Str {
				return errors.New("field has size must be string")
			}
			if _, ok := v.Value.(int); !ok {
				return errors.New("size attribute must be integer")
			}
			hasSize = true
		case Comment:
		case Nullable:
		default:
			return errors.New("unknown field: " + n)
		}
	}

	if t == Str && !hasSize {
		return errors.New("field with str type must have size attribute")
	}

	return nil
}

func checkIndex(fields []yaml.MapSlice, indexes [][]string) error {
	for _, idx := range indexes {
	Field:
		for _, field := range idx {
			for _, f := range fields {
				if f[0].Key.(string) == field {
					continue Field
				}
			}
			return errors.New("index field not found: " + field)
		}
	}
	return nil
}
