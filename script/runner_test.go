package script

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rusinikita/changes/change"
	"github.com/rusinikita/changes/commit"
	"github.com/rusinikita/changes/commit/value"
	"github.com/rusinikita/changes/conf/mock"
	"github.com/rusinikita/changes/errors"
)

func TestGetScriptsRunner(t *testing.T) {
	t.Run("unmarshall err", func(t *testing.T) {
		config := &mock.ConfMock{
			UnmarshalFunc: func(key string, config any, defaultValue ...any) error {
				return errors.New("test")
			},
		}

		scripts, err := GetScriptsRunner(config, value.DefaultProperties)

		assert.Nil(t, scripts)
		assert.Error(t, err)
	})

	t.Run("partial success", func(t *testing.T) {
		configs := []Config{
			{Message: "good1", Func: "true"},
			{Message: "err1"},
			{Message: "err2"},
			{Message: "good2", Func: "true"},
		}
		config := &mock.ConfMock{
			UnmarshalFunc: func(key string, config any, defaultValue ...any) error {
				reflect.ValueOf(config).Elem().Set(reflect.ValueOf(configs))

				return nil
			},
		}

		scripts, err := GetScriptsRunner(config, value.DefaultProperties)

		assert.NotEmpty(t, scripts)
		assert.Equal(t, 2, errors.Len(err))
		assert.Len(t, scripts.(*runner).scripts, 2)
	})
}

func Test_runner_Run(t *testing.T) {
	configs := []Config{
		{Message: "err", Func: "true"},
		{Message: "err", Func: "commits"},
		{Message: "err", Func: "changes"},
		{Message: "err", Func: "1/0 == 1"},
		{Message: "no err", Func: "false"},
	}
	config := &mock.ConfMock{
		UnmarshalFunc: func(key string, config any, defaultValue ...any) error {
			reflect.ValueOf(config).Elem().Set(reflect.ValueOf(configs))

			return nil
		},
	}

	scripts, err := GetScriptsRunner(config, value.DefaultProperties)
	require.NoError(t, err)

	err = scripts.Run(make([]commit.Commit, 3), make([]change.Change, 3))

	assert.Equal(t, 8, errors.Len(err))
}
