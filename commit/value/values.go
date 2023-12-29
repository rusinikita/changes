package value

import (
	"fmt"
	"regexp"

	"golang.org/x/exp/maps"
)

type Properties map[Name]Property

type Values map[Name]Value

type Name string

type Value string

type propertyConf struct {
	Regexp  string
	Allowed []string
	MaxLen  int `mapstructure:"max-len"`
}

type Property struct {
	Regexp  *regexp.Regexp
	Allowed map[string]struct{}
	MaxLen  int
}

func (p Property) Validate(v string) string {
	if p.MaxLen > 0 && len(v) > p.MaxLen {
		return fmt.Sprintf("longer than %d (%d)", p.MaxLen, len(v))
	}

	if len(p.Allowed) > 0 {
		if _, ok := p.Allowed[v]; ok {
			return ""
		}

		return fmt.Sprintf("value must be one of %v", maps.Keys(p.Allowed))
	}

	if p.Regexp.MatchString(v) {
		return ""
	}

	return fmt.Sprintf("does not match regexp `%s`", p.Regexp.String())
}
