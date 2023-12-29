package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	m := NewUnmarshal("test", nil)

	v := "wrong"

	assert.NoError(t, m.Unmarshal("key", &v))
	assert.Equal(t, "test", v)
}
