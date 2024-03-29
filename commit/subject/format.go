package subject

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rusinikita/changes/commit/value"
	"github.com/rusinikita/changes/errors"
)

const (
	tileValueRegexp = "?P<%s>.+"
	valueRegexp     = `?P<%s>[^\(\)\[\] ]+`
	titleErr        = "must contain a title"
)

var (
	valueNameRegexp = regexp.MustCompile("[a-z-_]+")
	replacer        = strings.NewReplacer(`((`, `(\((`, `))`, `)\))`, `[`, `\[`, `]`, `\]`)
)

func formatRegexp(format string) (string, []string) {
	found := valueNameRegexp.FindAllString(format, -1)

	format = replacer.Replace(format)

	for _, s := range found {
		regexpReplacing := valueRegexp
		if s == value.TitleValue {
			regexpReplacing = tileValueRegexp
		}

		format = strings.ReplaceAll(format, s, fmt.Sprintf(regexpReplacing, s))
	}

	format = "^" + format + "$"

	return format, found
}

func checkValues(foundValues []string, values value.Properties) error {
	var (
		errs           error
		notFoundValues []string
		hasTitle       bool
	)

	for _, v := range foundValues {
		name := value.Name(v)
		if _, ok := values[name]; !ok {
			notFoundValues = append(notFoundValues, v)
			continue
		}

		if v == value.TitleValue {
			hasTitle = true
		}
	}

	if len(notFoundValues) > 0 {
		errs = errors.Add(errs, fmt.Errorf("unknown keys: %s", notFoundValues))
	}

	if !hasTitle {
		errs = errors.Add(errs, errors.New(titleErr))
	}

	return errs
}
