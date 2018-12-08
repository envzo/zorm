package parse

import (
	"strings"

	"gopkg.in/yaml.v2"
)

func Parse(b []byte) (map[string]*Def, error) {
	m := map[string]*Def{}
	if err := yaml.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	for name, v := range m {
		// pre repair data

		m[name].Engine = strings.ToLower(m[name].Engine)

		if v.TB == "" {
			m[name].TB = name
		}

		if err := check(v); err != nil {
			return nil, err
		}
	}
	return m, nil
}
