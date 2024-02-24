package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownOutput(t *testing.T) {
	expected := `# Changes report
## Commits
- no more than 4 commits

**9999999: bla bla**
- some error
- some error

**2141234: blabla bla**
- some error
- some error

## Changes
**commit/value/get.go**
- some error
- some error`

	assert.Equal(t, expected, MarkdownOutput(testTree))
}

func TestMarkdownEmptyOutput(t *testing.T) {
	expected := `# Changes report
OK`

	assert.Equal(t, expected, MarkdownOutput(nil))
}
