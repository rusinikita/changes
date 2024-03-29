package change

import (
	"github.com/rusinikita/changes/errors"
	"github.com/rusinikita/changes/git"
)

type Change git.FileChange

func Changes(gitChanges []git.FileChange) []Change {
	result := make([]Change, 0, len(gitChanges))

	for _, change := range gitChanges {
		result = append(result, Change(change))
	}

	return result
}

func (c Change) ErrPrefix() []string {
	p := c.Path
	if p == "" {
		p = c.PrevPath
	}

	return []string{errors.FilesDiffGroup, p}
}
