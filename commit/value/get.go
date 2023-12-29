package value

import (
	"regexp"

	"github.com/rusinikita/changes/conf"
	"github.com/rusinikita/changes/errors"
)

const (
	ErrorNameLen   = "value name must be 3 - 20 characters long"
	ErrorNameChars = "value name must be only latin lowercase characters or -_"
	NameMinLen     = 3
	NameMaxLen     = 20

	ErrorValueEmpty = "regexp or allowed must be filled"
	ErrorValueBoth  = "regexp and allowed can't be filled together"

	Regexp          = "regexp"
	ValuesConfigKey = "commit.values"

	TitleValue   = "title"
	TypeValue    = "type"
	ContextValue = "context"
	IssueValue   = "issue"
)

var (
	NameRegex         = regexp.MustCompile("^[a-z-_]+$")
	DefaultProperties = Properties{
		TypeValue: {
			Allowed: map[string]struct{}{
				"fix":  {},
				"feat": {},
			},
		},
		ContextValue: {
			Regexp: regexp.MustCompile("^[a-z-_]+$"),
			MaxLen: 10,
		},
		TitleValue: {
			Regexp: regexp.MustCompile(`^[\w ]+$`),
			MaxLen: 50,
		},
		IssueValue: {
			Regexp: regexp.MustCompile(`^#\d+$`),
			MaxLen: 5,
		},
	}
)

func Get(c conf.Conf) (Properties, error) {
	configValues := map[string]propertyConf{}

	err := c.Unmarshal(ValuesConfigKey, &configValues)
	if err != nil {
		return nil, errors.Prefix(err, ValuesConfigKey)
	}

	var (
		values  = DefaultProperties
		value   Property
		vErr    error
		vErrors error
	)

	for name, config := range configValues {
		vErr = validateName(name)

		vErrors = errors.Add(vErrors, vErr, name)

		value, vErr = parseValue(config)

		vErrors = errors.Add(vErrors, vErr, name)

		values[Name(name)] = value
	}

	return values, vErrors
}

func validateName(s string) error {
	if len(s) < NameMinLen || len(s) > NameMaxLen {
		return errors.New(ErrorNameLen)
	}

	if !NameRegex.MatchString(s) {
		return errors.New(ErrorNameChars)
	}

	return nil
}

func parseValue(c propertyConf) (Property, error) {
	if len(c.Allowed) == 0 && len(c.Regexp) == 0 {
		return Property{}, errors.New(ErrorValueEmpty)
	}

	if len(c.Allowed) > 0 && len(c.Regexp) > 0 {
		return Property{}, errors.New(ErrorValueBoth)
	}

	exp, err := mapRegex(c.Regexp)
	if err != nil {
		return Property{}, errors.Prefix(err, Regexp)
	}

	v := Property{
		MaxLen:  c.MaxLen,
		Allowed: mapAllowed(c.Allowed),
		Regexp:  exp,
	}

	return v, nil
}

func mapAllowed(allowed []string) map[string]struct{} {
	if len(allowed) == 0 {
		return nil
	}

	result := map[string]struct{}{}

	for _, s := range allowed {
		result[s] = struct{}{}
	}

	return result
}

func mapRegex(r string) (*regexp.Regexp, error) {
	if len(r) == 0 {
		return nil, nil
	}

	return regexp.Compile(r)
}
