package internal

import (
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type Change struct {
	Repository *git.Repository
	// PR source branch ref
	SourceBranch plumbing.ReferenceName
	// PR target branch ref
	TargetBranch plumbing.ReferenceName
	// PR source branch last commit
	SourceLastCommit *object.Commit
	// PR target branch last commit
	TargetLastCommit *object.Commit
	// Source and target latest common commit
	CommonCommit *object.Commit
}

// InitPR - PR
func InitPR(repository *git.Repository, target plumbing.ReferenceName) (history Change, err error) {
	// target
	targetBranch, err := repository.Reference(target, true)
	if err != nil {
		return history, fmt.Errorf("%w '%s'", err, target)
	}

	targetCommit, err := repository.CommitObject(targetBranch.Hash())
	if err != nil {
		return history, err
	}

	// source
	sourceHead, err := repository.Head()
	if err != nil {
		return history, err
	}

	sourceCommit, err := repository.CommitObject(sourceHead.Hash())
	if err != nil {
		return history, err
	}

	// merge changes start
	base, err := sourceCommit.MergeBase(targetCommit)
	if err != nil {
		return history, err
	}

	if len(base) == 0 {
		return history, errors.New("no merge base commits")
	}

	return Change{
		Repository:       repository,
		SourceBranch:     sourceHead.Name(),
		TargetBranch:     targetBranch.Name(),
		SourceLastCommit: sourceCommit,
		TargetLastCommit: targetCommit,
		CommonCommit:     base[0],
	}, nil
}

// InitLast - last N commits
func InitLast(rep *git.Repository, lastCommitsCount int) (history Change, err error) {
	history.Repository = rep

	head, err := rep.Head()
	if err != nil {
		return history, err
	}

	history.SourceBranch = head.Name()

	commit, err := rep.CommitObject(head.Hash())
	if err != nil {
		return history, err
	}

	history.SourceLastCommit = commit
	history.CommonCommit, err = findNCommit(rep, lastCommitsCount)

	return history, err
}

func findNCommit(rep *git.Repository, lastCommitsCount int) (result *object.Commit, err error) {
	log, err := rep.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	err = log.ForEach(func(commit *object.Commit) error {
		result = commit

		if lastCommitsCount == 0 {
			return storer.ErrStop
		}

		lastCommitsCount--

		return nil
	})

	return result, err
}
