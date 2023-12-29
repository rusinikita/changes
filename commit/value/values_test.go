package value

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProperty_Validate(t *testing.T) {
	t.Run("len", func(t *testing.T) {
		p := Property{
			MaxLen: 10,
		}

		err := p.Validate(strings.Repeat("a", 11))

		assert.Contains(t, err, "10")
	})

	t.Run("not allowed", func(t *testing.T) {
		p := DefaultProperties[TypeValue]

		err := p.Validate("test")

		assert.Contains(t, err, "must be one of")
	})

	t.Run("not match", func(t *testing.T) {
		p := DefaultProperties[IssueValue]

		err := p.Validate("test")

		assert.Contains(t, err, "not match")
	})

	t.Run("not match", func(t *testing.T) {
		p := DefaultProperties[IssueValue]

		err := p.Validate("#123")

		assert.Empty(t, err)
	})
}
