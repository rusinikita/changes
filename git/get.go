package git

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/rusinikita/changes/git/internal"
)

const defaultLastCommits = 10

type getSettings struct{}

type ChangeOption func(settings *getSettings)

func GetChange(opts ...ChangeOption) (Change, error) {
	deps := getDeps{
		PlainOpenWithOptions: git.PlainOpenWithOptions,
		InitLast:             internal.InitLast,
		InitPR:               internal.InitPR,
	}

	return get(deps, opts)
}

func get(deps getDeps, _ []ChangeOption) (Change, error) {
	rep, err := deps.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})
	if err != nil {
		return nil, err
	}

	var (
		result           = change{}
		isCI             = os.Getenv("CI") == "true"
		targetBranchName = os.Getenv("CI_MERGE_REQUEST_TARGET_BRANCH_NAME")
	)

	if targetBranchName == "" {
		targetBranchName = os.Getenv("GITHUB_BASE_REF")
	}

	if !isCI || targetBranchName == "" {
		result.repository, err = deps.InitLast(rep, defaultLastCommits)
		if err != nil {
			return nil, err
		}

		return &result, nil
	}

	result.repository, err = deps.InitPR(rep, plumbing.NewRemoteReferenceName("origin", targetBranchName))
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type getDeps struct {
	PlainOpenWithOptions func(path string, o *git.PlainOpenOptions) (*git.Repository, error)
	InitLast             func(rep *git.Repository, lastCommitsCount int) (history internal.Change, err error)
	InitPR               func(repository *git.Repository, target plumbing.ReferenceName) (internal.Change, error)
}
