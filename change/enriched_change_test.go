package change

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rusinikita/changes/errors"
)

func TestChange_ErrPrefix(t *testing.T) {
	c := Change{
		Path:     "test",
		PrevPath: "prev_test",
	}

	assert.Equal(t, errors.FilesDiffGroup+"/test", c.ErrPrefix())

	c.Path = ""
	assert.Equal(t, errors.FilesDiffGroup+"/prev_test", c.ErrPrefix())
}
