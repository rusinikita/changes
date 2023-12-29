package subject

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rusinikita/changes/commit/value"
)

func Test_formatRegexp(t *testing.T) {
	tests := []struct {
		format     string
		wantRegexp string
		wantValues []string
	}{
		{
			format:     `(type) (title)`,
			wantRegexp: `^(?P<type>[^\(\)\[\] ]+) (?P<title>.+)$`,
			wantValues: []string{"type", "title"},
		},
		{
			format:     `(type) (issue) (title)`,
			wantRegexp: `^(?P<type>[^\(\)\[\] ]+) (?P<issue>[^\(\)\[\] ]+) (?P<title>.+)$`,
			wantValues: []string{"type", "issue", "title"},
		},
		{
			format:     `(type)((context))?(!)?: (title)`,
			wantRegexp: `^(?P<type>[^\(\)\[\] ]+)(\((?P<context>[^\(\)\[\] ]+)\))?(!)?: (?P<title>.+)$`,
			wantValues: []string{"type", "context", "title"},
		},
		{
			format:     `(task) ([context])? (type) (title)`,
			wantRegexp: `^(?P<task>[^\(\)\[\] ]+) (\[?P<context>[^\(\)\[\] ]+\])? (?P<type>[^\(\)\[\] ]+) (?P<title>.+)$`,
			wantValues: []string{"task", "context", "type", "title"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			gotRegexp, gotValues := formatRegexp(tt.format)

			assert.Equal(t, tt.wantRegexp, gotRegexp)
			assert.Equal(t, tt.wantValues, gotValues)
		})
	}
}

func Test_checkValues(t *testing.T) {
	values := value.Properties{
		"test": {},
		"name": {},
	}

	t.Run("title err", func(t *testing.T) {
		err := checkValues([]string{"test"}, values)

		assert.Contains(t, err.Error(), titleErr)
	})

	t.Run("not found err", func(t *testing.T) {
		err := checkValues([]string{"test_test", "title"}, values)

		assert.Contains(t, err.Error(), "unknown keys")
	})

	t.Run("not found err", func(t *testing.T) {
		err := checkValues([]string{"test_test", "bla"}, values)

		assert.Contains(t, err.Error(), titleErr)
		assert.Contains(t, err.Error(), "unknown keys")
		assert.Contains(t, err.Error(), "test_test")
		assert.Contains(t, err.Error(), "bla")
	})
}
