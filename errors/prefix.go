package errors

import (
	"net/url"
	"strings"
	"unicode"
)

func StrToPathPrefix(s string) string {
	const maxLen = 20

	if len(s) > maxLen {
		s = s[:maxLen]
	}

	s = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return '_'
		}

		return r
	}, s)

	return url.PathEscape(s)
}
