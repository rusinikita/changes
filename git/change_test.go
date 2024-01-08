package git

import (
	"strings"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rusinikita/changes/git/internal"
)

func Test_Changes(t *testing.T) {
	testBranch := "new_branch"
	rep := internal.TestRepository(t,
		internal.Commits(internal.Range(1, 10)),
		internal.NewBranch(testBranch, 4),
		internal.Commits(internal.Range(11, 20)),
	)

	history, err := internal.InitPR(rep, plumbing.NewBranchReferenceName("master"))
	require.NoError(t, err)

	g := change{history}
	changes, _ := g.FilesDiff()

	assert.Len(t, changes, 10)

	for _, change := range changes {
		assert.Empty(t, change.PrevPath)
		assert.NotEmpty(t, change.Path)
		assert.Len(t, change.Chunks, 1)

		for _, chunk := range change.Chunks {
			assert.Equal(t, Add, chunk.Type)
			assert.NotEmpty(t, chunk.Content)
		}
	}
}

func TestFileChange_Stats(t *testing.T) {
	c := FileChange{
		Path:     "path",
		PrevPath: "prevPath",
		Chunks: []Chunk{
			{
				Type:    Add,
				Content: strings.Repeat("line\n", 134),
			},
			{
				Type:    Equal,
				Content: strings.Repeat("line\n", 134),
			},
			{
				Type:    Delete,
				Content: strings.Repeat("line\n", 130),
			},
			{
				Type:    Delete,
				Content: strings.Repeat("line\n", 4),
			},
			{
				Type:    Add,
				Content: strings.Repeat("line\n", 13),
			},
			{
				Type:    Add,
				Content: "line\nlast",
			},
		},
	}

	stats := c.Stats()

	assert.Equal(t, 149, stats.Additions)
	assert.Equal(t, 134, stats.Deletions)
}
