package output

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rusinikita/changes/errors"
)

var (
	testErr = "some error"

	//nolint:dupword
	testTree = errors.OutputTree{
		{
			Name:     errors.CommitGroup,
			Messages: []string{"no more than 4 commits"},
			Groups: []errors.Node{
				{
					Name: "9999999: bla bla",
					Messages: []string{
						testErr,
						testErr,
					},
				},
				{
					Name: "2141234: blabla bla",
					Messages: []string{
						testErr,
						testErr,
					},
				},
			},
		},
		{
			Name: errors.FilesDiffGroup,
			Groups: []errors.Node{{
				Name: "commit/value/get.go",
				Messages: []string{
					testErr,
					testErr,
				},
			}},
		},
	}
)

//nolint:dupword
func TestTerminalOutput(t *testing.T) {
	expected := `
Commits
-------
- no more than 4 commits

9999999: bla bla
- some error
- some error

2141234: blabla bla
- some error
- some error

Changes
-------
commit/value/get.go
- some error
- some error`

	assert.Equal(t, strings.TrimSpace(expected), TerminalOutput(testTree))
}
