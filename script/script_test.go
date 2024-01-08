package script

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/rusinikita/changes/change"
	"github.com/rusinikita/changes/commit"
	"github.com/rusinikita/changes/commit/value"
	"github.com/rusinikita/changes/git"
)

func TestNew(t *testing.T) {
	props := maps.Clone(value.DefaultProperties)
	props[value.Name("custom")] = value.Property{}

	env, err := createEnv(props)
	require.NoError(t, err)

	tests := []struct {
		name        string
		script      string
		err         bool
		errContains string
	}{
		{
			name:   "ok default value, custom value, diff usage",
			script: `commits.exists(c, c.title == c.custom) && changes.exists(script, script.path == "test")`,
			err:    false,
		},
		{
			name:   "ok commits return",
			script: `commits.filter(c, c.title == c.custom)`,
			err:    false,
		},
		{
			name:   "ok changes return",
			script: `changes.filter(c, c.path == "test")`,
			err:    false,
		},
		{
			name:   "unknown return type",
			script: `["test"]`,
			err:    true,
		},
		{
			name:   "unknown return type",
			script: `"test"`,
			err:    true,
		},
		{
			name:        "message or func is empty",
			script:      "",
			err:         true,
			errContains: "message and func are required",
		},
		{
			name:   "error: unknown variable usage",
			script: "commit",
			err:    true,
		},
		{
			name:   "error: unknown commit field usage",
			script: "commits.exists(c, c.test == `test`)",
			err:    true,
		},
		{
			name:   "error: unknown diff method usage",
			script: `changes.exists(script, script.test == "test")`,
			err:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := Config{
				Message: "test",
				Func:    test.script,
			}

			s, err := New(c, env)

			if test.err {
				assert.Error(t, err)
				if test.errContains != "" {
					assert.ErrorContains(t, err, test.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, s)
				assert.Equal(t, "test", s.message)
			}
		})
	}
}

func TestScript_Run(t *testing.T) {
	props := maps.Clone(value.DefaultProperties)
	commits := []commit.Commit{
		{
			Commit: git.Commit{
				Author: git.Signature{
					When: time.Now(),
				},
			},
			Values: value.Values{
				value.TitleValue: "test",
			},
		},
		{
			Commit: git.Commit{},
			Values: value.Values{
				value.TitleValue: "test",
			},
		},
	}
	changes := []change.Change{
		{
			Path: "test",
		},
		{
			Path: "test.go",
		},
		{
			Path: "tes",
		},
	}

	env, err := createEnv(props)
	require.NoError(t, err)

	data := newData(commits, changes)

	tests := []struct {
		name     string
		script   string
		expected any
		err      bool
	}{
		{
			name:     "bool",
			script:   `commits.exists(script, script.title == "test")`,
			expected: true,
		},
		{
			name:     "has macros",
			script:   `has(commits[0].title)`,
			expected: true,
		},
		{
			name:     "timestamp usage",
			script:   `commits.exists(script, script.author.when > (now - duration("1h")))`,
			expected: true,
		},
		{
			name:     "single commit return",
			script:   `commits.filter(script, script.author.when > (now - duration("1h")))`,
			expected: []commit.Commit{commits[0]},
		},
		{
			name:     "commits list return",
			script:   `commits.filter(script, script.title == "test")`,
			expected: []commit.Commit{commits[0], commits[1]},
		},
		{
			name:     "changes single return",
			script:   `changes.filter(script, script.path == "test.go")`,
			expected: []change.Change{changes[1]},
		},
		{
			name:     "changes list return",
			script:   `changes.filter(script, script.path.size() > 3)`,
			expected: []change.Change{changes[0], changes[1]},
		},
		{
			name:     "changes empty list return",
			script:   `changes.filter(script, script.path.size() < 0)`,
			expected: nil,
		},
		{
			name:   "exec error",
			script: `10/0 == 0`,
			err:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := Config{
				Message: "test",
				Func:    test.script,
			}

			srt, err := New(c, env)
			require.NoError(t, err)

			result, err := srt.Run(data)

			if test.err {
				require.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.IsType(t, test.expected, result)
			assert.Equal(t, test.expected, result)
		})
	}
}
