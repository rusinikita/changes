package internal

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitLast(t *testing.T) {
	rep := TestRepository(t, Commits(Range(1, 10)))

	t.Run("ok", func(t *testing.T) {
		change, err := InitLast(rep, 5)
		require.NoError(t, err)

		assert.Equal(t, "master", change.SourceBranch.Short())
		assert.Empty(t, change.TargetBranch)
		assert.Equal(t, "10", change.SourceLastCommit.Message)
		assert.Nil(t, change.TargetLastCommit)
		assert.Equal(t, "5", change.CommonCommit.Message)
		assert.NotNil(t, change.Repository)
	})

	t.Run("not enough commits", func(t *testing.T) {
		change, err := InitLast(rep, 11)
		require.NoError(t, err)

		assert.Equal(t, "1", change.CommonCommit.Message)
	})
}

func TestInitPR(t *testing.T) {
	var (
		testBranch = "test_branch"
		master     = "master"
		rep        = TestRepository(t,
			Commits(Range(1, 10)),
			NewBranch(testBranch, 4),
			Commits(Range(11, 20)),
		)
	)

	t.Run("ok", func(t *testing.T) {
		change, err := InitPR(rep, plumbing.NewBranchReferenceName(master))
		require.NoError(t, err)

		assert.Equal(t, testBranch, change.SourceBranch.Short())
		assert.Equal(t, master, change.TargetBranch.Short())
		assert.Equal(t, "20", change.SourceLastCommit.Message)
		assert.Equal(t, "10", change.TargetLastCommit.Message)
		assert.Equal(t, "6", change.CommonCommit.Message)
		assert.NotNil(t, change.Repository)
	})

	t.Run("no merge base", func(t *testing.T) {
		t.Skip("TODO")
	})
}
