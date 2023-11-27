package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var long = strings.Repeat("long ", 20)

func Test_check(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		err := check(commit("normal"), nil)

		assert.NoError(t, err)
	})

	t.Run("long title", func(t *testing.T) {
		err := check(commit(long), nil)

		assert.Error(t, err)
	})

	t.Run("long body", func(t *testing.T) {
		err := check(commit("title\n"+long), nil)

		assert.NoError(t, err)
	})

	t.Run("skip merge", func(t *testing.T) {
		longM := "Merge" + long

		err := check(commit(longM), nil)

		assert.NoError(t, err)
	})

	t.Run("many changes", func(t *testing.T) {
		stats := object.FileStats{object.FileStat{
			Name:     "test.go",
			Addition: 401,
		}}

		err := check(commit("normal"), stats)

		assert.Error(t, err)
	})

	t.Run("skip no go changes", func(t *testing.T) {
		stats := object.FileStats{object.FileStat{
			Name:     "test",
			Addition: 401,
		}}

		err := check(commit("normal"), stats)

		assert.NoError(t, err)
	})
}

func Test_repositoryChecks(t *testing.T) {
	t.Run("ok, more 10", func(t *testing.T) {
		rep, err := git.Init(memory.NewStorage(), memfs.New())
		require.NoError(t, err)

		commits := append([]string{long}, stringRange(0, 10)...)
		fillCommits(commits, t, rep)

		err = repositoryChecks(rep)
		assert.NoError(t, err)
	})

	t.Run("long found", func(t *testing.T) {
		rep, err := git.Init(memory.NewStorage(), memfs.New())
		require.NoError(t, err)

		commits := append([]string{long}, stringRange(0, 5)...)
		fillCommits(commits, t, rep)

		err = repositoryChecks(rep)

		assert.Error(t, err)
	})

	t.Run("log error", func(t *testing.T) {
		rep := repMock{
			err: errors.New("test"),
		}

		err := repositoryChecks(rep)

		assert.Error(t, err)
	})

	t.Run("log error", func(t *testing.T) {
		rep := repMock{t: t}

		err := repositoryChecks(rep)

		assert.Error(t, err)
	})
}

func commit(message string) *object.Commit {
	return &object.Commit{Message: message}
}

func fillCommits(cc []string, t *testing.T, r *git.Repository) {
	wt, err := r.Worktree()
	require.NoError(t, err)

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

func stringRange(f, l int) []string {
	result := make([]string, 0, l-f)

	for i := f; i <= l; i++ {
		result = append(result, strconv.Itoa(i))
	}

	return result
}

func Test_main(_ *testing.T) {
	_ = run()
}

type repMock struct {
	t   *testing.T
	err error
}

func (repMock) Next() (*object.Commit, error) {
	return nil, errors.New("implement me")
}

func (r repMock) ForEach(f func(*object.Commit) error) error {
	rep, err := git.Init(memory.NewStorage(), memfs.New())
	assert.NoError(r.t, err)

	fillCommits([]string{"test"}, r.t, rep)

	iter, err := rep.Log(&git.LogOptions{})
	assert.NoError(r.t, err)

	return iter.ForEach(func(c *object.Commit) error {
		c.TreeHash = plumbing.NewHash("")

		return f(c)
	})
}

func (repMock) Close() {}

func (r repMock) Log(_ *git.LogOptions) (object.CommitIter, error) {
	return r, r.err
}
