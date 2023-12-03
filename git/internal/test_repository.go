package internal

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/require"
)

// TestRepository creates
func TestRepository(t *testing.T, steps ...Step) *git.Repository {
	rep, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err)

	// create test branches
	worktree, err := rep.Worktree()
	require.NoError(t, err)

	for _, o := range steps {
		o(t, rep, worktree)
	}

	return rep
}

type Step func(t *testing.T, r *git.Repository, wt *git.Worktree)

// Commits creates list new file commits with file name, content and message for list
func Commits(cc []string) Step {
	return func(t *testing.T, r *git.Repository, wt *git.Worktree) {
		for _, commit := range cc {
			fileName := fmt.Sprintf("%s.txt", commit)
			file, err := wt.Filesystem.Create(fileName)
			require.NoError(t, err)

			_, err = file.Write([]byte(commit))
			require.NoError(t, err)
			require.NoError(t, file.Close())

			_, err = wt.Add(fileName)
			require.NoError(t, err)
			_, err = wt.Commit(commit, &git.CommitOptions{
				Author: &object.Signature{
					Name:  commit,
					Email: commit,
				},
				Committer: &object.Signature{
					Name:  commit,
					Email: commit,
				},
			})
			require.NoError(t, err)
		}
	}
}

// NewBranch creates new branch with parent commit = current - prevCommits
func NewBranch(name string, prevCommits int) Step {
	return func(t *testing.T, r *git.Repository, wt *git.Worktree) {
		head, err := r.Head()
		require.NoError(t, err)

		commit, err := r.CommitObject(head.Hash())
		require.NoError(t, err)

		for i := 0; i < prevCommits; i++ {
			commit, err = commit.Parent(0)
			require.NoError(t, err)
		}

		err = wt.Checkout(&git.CheckoutOptions{
			Hash:   commit.Hash,
			Branch: plumbing.NewBranchReferenceName(name),
			Create: true,
		})
		require.NoError(t, err)
	}
}

func Range(first, last int) (result []string) {
	for i := first; i <= last; i++ {
		result = append(result, strconv.Itoa(i))
	}

	return result
}
