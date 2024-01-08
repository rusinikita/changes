package change

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChange_ErrPrefix(t *testing.T) {
	c := Change{
		Path:     "test",
		PrevPath: "prev_test",
	}

	assert.Equal(t, "file/test", c.ErrPrefix())

	c.Path = ""
	assert.Equal(t, "file/prev_test", c.ErrPrefix())
}
