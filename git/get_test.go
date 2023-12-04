package git

import (
	"errors"
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rusinikita/changes/git/internal"
)

func Test_get(t *testing.T) {
	var (
		testErr = errors.New("test")
		testRep = &git.Repository{}
	)

	require.NoError(t, os.Setenv("CI", ""))
	require.NoError(t, os.Setenv("GITHUB_BASE_REF", ""))

	t.Run("open error", func(t *testing.T) {
		deps := getDeps{
			PlainOpenWithOptions: func(path string, o *git.PlainOpenOptions) (*git.Repository, error) {
				return nil, testErr
			},
		}

		result, err := get(deps, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("last 10", func(t *testing.T) {
		callArgument := 0
		deps := getDeps{
			PlainOpenWithOptions: func(path string, o *git.PlainOpenOptions) (*git.Repository, error) {
				return testRep, nil
			},
			InitLast: func(rep *git.Repository, lastCommitsCount int) (history internal.Change, err error) {
				callArgument = lastCommitsCount

				return internal.Change{}, nil
			},
		}

		require.NoError(t, os.Setenv("CI", "true"))

		result, err := get(deps, nil)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 10, callArgument)
	})

	t.Run("last err", func(t *testing.T) {
		deps := getDeps{
			PlainOpenWithOptions: func(path string, o *git.PlainOpenOptions) (*git.Repository, error) {
				return testRep, nil
			},
			InitLast: func(rep *git.Repository, lastCommitsCount int) (history internal.Change, err error) {
				return internal.Change{}, testErr
			},
		}

		result, err := get(deps, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("PR", func(t *testing.T) {
		var callArgument plumbing.ReferenceName
		deps := getDeps{
			PlainOpenWithOptions: func(path string, o *git.PlainOpenOptions) (*git.Repository, error) {
				return testRep, nil
			},
			InitPR: func(repository *git.Repository, target plumbing.ReferenceName) (history internal.Change, err error) {
				callArgument = target

				return internal.Change{}, nil
			},
		}

		require.NoError(t, os.Setenv("CI", "true"))
		require.NoError(t, os.Setenv("GITHUB_BASE_REF", "main"))

		result, err := get(deps, nil)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, plumbing.NewRemoteReferenceName("origin", "main"), callArgument)
	})

	t.Run("PR err", func(t *testing.T) {
		deps := getDeps{
			PlainOpenWithOptions: func(path string, o *git.PlainOpenOptions) (*git.Repository, error) {
				return testRep, nil
			},
			InitPR: func(repository *git.Repository, target plumbing.ReferenceName) (history internal.Change, err error) {
				return internal.Change{}, testErr
			},
		}

		result, err := get(deps, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
