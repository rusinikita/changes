package commit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rusinikita/changes/commit/subject"
	"github.com/rusinikita/changes/commit/value"
	"github.com/rusinikita/changes/conf/mock"
	"github.com/rusinikita/changes/errors"
	"github.com/rusinikita/changes/git"
)

func TestGetParser(t *testing.T) {
	t.Run("get error", func(t *testing.T) {
		config := &mock.ConfMock{
			GetStringFunc: func(key string, defaultValue ...string) string {
				return "???"
			},
			UnmarshalFunc: func(key string, config any, defaultValue ...any) error {
				return errors.New("test")
			},
		}

		parser, err := GetParser(config)
		assert.Nil(t, parser)
		assert.Equal(t, 3, errors.Len(err))
	})

	t.Run("parse", func(t *testing.T) {
		parser, err := subject.NewParser("(type): (title)", value.DefaultProperties)
		require.NoError(t, err)

		p := Parser{
			Parser: parser,
		}

		commits := []git.Commit{
			{
				Message: "bad",
			},
			{
				Message: "bad2",
			},
			{
				Message: "feat: good",
			},
		}

		enrichedCommits, err := p.Parse(commits)

		assert.Equal(t, 2, errors.Len(err))
		assert.Len(t, enrichedCommits, 3)
		assert.EqualValues(t, "feat", enrichedCommits[2].Values[value.TypeValue])
	})
}
