package subject

import (
	"fmt"
	"regexp"

	"github.com/rusinikita/changes/commit/value"
	"github.com/rusinikita/changes/errors"
)

type Parser struct {
	format     string
	regexp     *regexp.Regexp
	properties value.Properties
}

func NewParser(format string, values value.Properties) (*Parser, error) {
	regexpString, foundValues := formatRegexp(format)

	err := checkValues(foundValues, values)

	compiledRegexp, compileErr := regexp.Compile(regexpString)
	err = errors.Add(err, compileErr)

	return &Parser{
		format:     format,
		regexp:     compiledRegexp,
		properties: values,
	}, err
}

func (p Parser) Parse(subject string) (values value.Values, err error) {
	values = value.Values{}

	matches := p.regexp.FindStringSubmatch(subject)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no match to format: %s", p.format)
	}

	expNames := p.regexp.SubexpNames()

	for i, n := range expNames {
		if n == "" {
			continue
		}

		v := matches[i]
		if v == "" {
			continue
		}

		name := value.Name(n)

		valueErr := p.properties[name].Validate(v)
		if valueErr != "" {
			err = errors.Add(err, errors.New(valueErr), n)
		}

		values[name] = value.Value(v)
	}

	return values, err
}
