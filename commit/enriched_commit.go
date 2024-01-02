package commit

import (
	"github.com/rusinikita/changes/commit/value"
	"github.com/rusinikita/changes/git"
)

type Commit struct {
	git.Commit
	value.Values
}
