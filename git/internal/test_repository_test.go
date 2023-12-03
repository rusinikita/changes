package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange(t *testing.T) {
	expected := []string{"4", "5", "6", "7", "8", "9", "10", "11"}

	assert.Equal(t, expected, Range(4, 11))
}
