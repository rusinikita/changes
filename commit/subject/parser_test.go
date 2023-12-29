package subject

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rusinikita/changes/commit/value"
)

func TestNewParser(t *testing.T) {
	t.Run("regexp err", func(t *testing.T) {
		values := value.Properties{value.TitleValue: {}}
		format := "( ( ( (title"
		parser, multiError := NewParser(format, values)

		require.NotNil(t, parser)
		assert.Error(t, multiError)

		assert.Nil(t, parser.regexp)
		assert.Equal(t, parser.properties, values)
		assert.Equal(t, parser.format, format)
	})
}

func TestParser_Parse(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		parser, err := NewParser("(type)((context))?: (title)", value.DefaultProperties)
		require.NoError(t, err)

		values, e := parser.Parse("fix: something fixed")

		expected := value.Values{
			"type":           "fix",
			value.TitleValue: "something fixed",
		}

		assert.NoError(t, e)
		assert.Equal(t, expected, values)
	})

	t.Run("err: wrong", func(t *testing.T) {
		parser, err := NewParser("(type): (title)", value.DefaultProperties)
		require.NoError(t, err)

		values, e := parser.Parse("123: something fixed")

		expected := value.Values{
			"type":           "123",
			value.TitleValue: "something fixed",
		}

		assert.Error(t, e)
		assert.Equal(t, expected, values)
	})

	t.Run("err: no value", func(t *testing.T) {
		parser, err := NewParser("(type)((context)): (title)", value.DefaultProperties)
		require.NoError(t, err)

		values, e := parser.Parse("fix: something fixed")

		assert.Error(t, e)
		assert.Empty(t, values)
	})
}
