package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testErr = New("test")

func TestAdd(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.Equal(t, Add(nil, testErr), testErr)
		assert.Equal(t, Add(testErr, nil), testErr)
		assert.NoError(t, Add(nil, nil))
	})

	t.Run("multierr", func(t *testing.T) {
		err1 := Add(Prefix(Prefix(testErr, "test"), "test"), testErr)
		err2 := Add(testErr, testErr, "test")

		err := Add(err1, err2, "test")

		expected := ": test\ntest/test: test\ntest/test: test\ntest: test"
		assert.EqualError(t, err, expected)
	})
}

func TestPrefix(t *testing.T) {
	t.Run("nil or empty", func(t *testing.T) {
		assert.NoError(t, Prefix(nil, "1"))
		assert.Equal(t, Prefix(testErr), testErr)
	})

	t.Run("add path", func(t *testing.T) {
		err := Prefix(testErr, "test")
		err = Prefix(err, "test")

		assert.EqualError(t, err, "test/test: test")
	})
}
