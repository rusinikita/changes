package git

import (
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// Commit data
type Commit struct {
	// Message is the commit message, contains arbitrary text.
	Message string
	// Author is the original author of the commit.
	Author Signature
	// Committer is the one performing the commit, might be different from Author.
	Committer Signature
	// ParentsCount is count of the parent commits of the commit.
	ParentsCount int
}

func (c Commit) Subject() string {
	return strings.SplitN(c.Message, "\n", 2)[0] //nolint:revive
}

// Signature is used to identify who and when created a commit or tag.
type Signature struct {
	// Name represents a person name. It is an arbitrary string.
	Name string
	// Email is an email, but it cannot be assumed to be well-formed.
	Email string
	// When is the timestamp of the signature.
	When time.Time
}

func (g *change) Commits() (commits []Commit, err error) {
	logHistory, err := g.repository.Repository.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	err = logHistory.ForEach(func(commit *object.Commit) error {
		if commit.Hash == g.repository.CommonCommit.Hash {
			return storer.ErrStop
		}

		commits = append(commits, Commit{
			Message:      commit.Message,
			Author:       signature(commit.Author),
			Committer:    signature(commit.Committer),
			ParentsCount: commit.NumParents(),
		})

		return nil
	})

	return commits, err
}

func signature(s object.Signature) Signature {
	return Signature{
		Name:  s.Name,
		Email: s.Email,
		When:  s.When,
	}
}
