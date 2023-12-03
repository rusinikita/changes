package git

import (
	"slices"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rusinikita/changes/git/internal"
)

func Test_GetCommits(t *testing.T) {
	testBranch := "new_branch"
	rep := internal.TestRepository(t,
		internal.Commits(internal.Range(1, 10)),
		internal.NewBranch(testBranch, 4),
		internal.Commits(internal.Range(11, 20)),
	)

	t.Run("last", func(t *testing.T) {
		history, err := internal.InitLast(rep, 15)
		require.NoError(t, err)

		g := change{history}
		commits, _ := g.Commits()

		expected := append(internal.Range(2, 6), internal.Range(11, 20)...)
		slices.Reverse(expected)

		assert.Len(t, commits, 15)
		for i, msg := range expected {
			c := commits[i]
			t.Run(msg, func(t *testing.T) {
				assert.Equal(t, msg, c.Message)
				assert.EqualValues(t, 1, c.ParentsCount)
				assert.Equal(t, msg, c.Author.Name)
				assert.Equal(t, msg, c.Committer.Name)
			})
		}
	})

	t.Run("pr", func(t *testing.T) {
		history, err := internal.InitPR(rep, plumbing.NewBranchReferenceName("master"))
		require.NoError(t, err)

		g := change{history}
		commits, _ := g.Commits()

		expected := internal.Range(11, 20)
		slices.Reverse(expected)

		assert.Len(t, commits, 10)
		for i, msg := range expected {
			c := commits[i]
			t.Run(msg, func(t *testing.T) {
				assert.Equal(t, msg, c.Message)
				assert.EqualValues(t, 1, c.ParentsCount)
				assert.Equal(t, msg, c.Author.Name)
				assert.Equal(t, msg, c.Committer.Name)
			})
		}
	})
}
