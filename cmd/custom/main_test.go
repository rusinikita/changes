package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rusinikita/changes/git"
)

var (
	longTitle    = strings.Repeat("long ", 20)
	largeContent = strings.Repeat("row\n", maxFileChangesLen+1)
)

func Test_checkDiff(t *testing.T) {
	t.Run("many changes", func(t *testing.T) {
		diff := []git.FileChange{{
			Path: "test.go",
			Chunks: []git.Chunk{{
				Type:    git.Add,
				Content: largeContent,
			}},
		}}

		err := checkDiff(diff)

		assert.Error(t, err)
	})

	t.Run("skip no go changes", func(t *testing.T) {
		diff := []git.FileChange{{
			Path: "test",
			Chunks: []git.Chunk{{
				Type:    git.Add,
				Content: largeContent,
			}},
		}}

		err := checkDiff(diff)

		assert.NoError(t, err)
	})
}

func Test_checkCommit(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		err := checkCommit(commit("normal"))

		assert.NoError(t, err)
	})

	t.Run("long title", func(t *testing.T) {
		err := checkCommit(commit(longTitle))

		assert.Error(t, err)
	})

	t.Run("long title", func(t *testing.T) {
		err := checkCommit(commit(longTitle))

		assert.Error(t, err)
	})

	t.Run("long body", func(t *testing.T) {
		err := checkCommit(commit("title\n" + longTitle))

		assert.NoError(t, err)
	})

	t.Run("skip merge", func(t *testing.T) {
		longM := "Merge" + longTitle

		err := checkCommit(commit(longM))

		assert.NoError(t, err)
	})

	t.Run("email check", func(t *testing.T) {
		c := commit("normal")
		c.Author.Email = "test@test.in"

		err := checkCommit(c)

		assert.Error(t, err)
	})
}

func Test_checkChange(t *testing.T) {
	testErr := errors.New("test")

	t.Run("empty ok", func(t *testing.T) {
		m := changeMock{}

		assert.NoError(t, changesCheck(m))
	})

	t.Run("commits err", func(t *testing.T) {
		m := changeMock{diffErr: testErr}

		assert.Error(t, changesCheck(m))
	})

	t.Run("diff err", func(t *testing.T) {
		m := changeMock{commitsErr: testErr}

		assert.Error(t, changesCheck(m))
	})

	t.Run("commits check", func(t *testing.T) {
		m := changeMock{commits: []git.Commit{{
			Message: longTitle,
		}}}

		assert.Error(t, changesCheck(m))
	})

	t.Run("diff check", func(t *testing.T) {
		m := changeMock{diff: []git.FileChange{{
			Path: "test.go",
			Chunks: []git.Chunk{{
				Type:    git.Add,
				Content: largeContent,
			}},
		}}}

		assert.Error(t, changesCheck(m))
	})
}

func commit(message string) git.Commit {
	return git.Commit{
		Message: message,
		Author: git.Signature{
			Name:  "Test",
			Email: "test@gmail.com",
		},
	}
}

func Test_main(_ *testing.T) {
	_ = run()
}

type changeMock struct {
	commits    []git.Commit
	commitsErr error
	diff       []git.FileChange
	diffErr    error
}

func (changeMock) Branches() (source, target string) {
	return "", ""
}

func (c changeMock) Commits() ([]git.Commit, error) {
	return c.commits, c.commitsErr
}

func (c changeMock) FilesDiff() ([]git.FileChange, error) {
	return c.diff, c.diffErr
}
