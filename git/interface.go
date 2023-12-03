package git

import "github.com/rusinikita/changes/git/internal"

// Change is abstraction upon git using go-git library
// Encapsulates work with pull request history diff
type Change interface {
	// Branches returns PR source and target branch
	Branches() (source, target string)
	// Commits returns PR commits
	Commits() ([]Commit, error)
	// FilesDiff returns PR repository files diff
	FilesDiff() ([]FileChange, error)
}

type change struct {
	repository internal.Change
}

func (g *change) Branches() (source string, target string) {
	return g.repository.SourceBranch.Short(), g.repository.TargetBranch.Short()
}
