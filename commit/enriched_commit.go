package commit

import (
	"path"

	"github.com/rusinikita/changes/commit/value"
	"github.com/rusinikita/changes/errors"
	"github.com/rusinikita/changes/git"
)

type Commit struct {
	git.Commit
	value.Values
}

func (c Commit) ErrPrefix() string {
	return path.Join(errors.CommitGroup, c.Subject())
}
