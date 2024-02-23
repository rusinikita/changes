package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:dupword
func TestMarkdownOutput(t *testing.T) {
	expected := `## Commits
- no more than 4 commits

9999999: bla bla
- some error
- some error

2141234: blabla bla
- some error
- some error

## Changes
commit/value/get.go
- some error
- some error`

	assert.Equal(t, expected, MarkdownOutput(testTree))
}
