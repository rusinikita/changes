package script

import (
	"reflect"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"

	"github.com/rusinikita/changes/change"
	"github.com/rusinikita/changes/commit"
	"github.com/rusinikita/changes/commit/value"
	"github.com/rusinikita/changes/git"
)

var (
	commitListType  = cel.ListType(cel.ObjectType("commit.Commit"))
	changesListType = cel.ListType(cel.ObjectType("change.Change"))
)

func newData(commits []commit.Commit, changes []change.Change) Data {
	return map[string]any{
		"commits": commits,
		"changes": changes,
		"now":     time.Now(),
	}
}

func createEnv(props value.Properties) (*cel.Env, error) {
	env, err := cel.NewEnv(
		ext.NativeTypes(
			reflect.ValueOf(commit.Commit{}),
			reflect.ValueOf(git.Signature{}),
			reflect.ValueOf(git.Chunk{}),
			reflect.ValueOf(change.Change{}),
			reflect.ValueOf(git.Stat{}),
		),
		cel.Variable("commits", commitListType),
		cel.Variable("changes", changesListType),
		cel.Variable("now", cel.TimestampType),

		customTypeProvider(props),
	)
	if err != nil {
		return nil, err
	}

	return env.Extend(
		cel.Function("stats",
			cel.MemberOverload(
				"file_change_stats",
				[]*cel.Type{cel.ObjectType("change.Change")},
				cel.ObjectType("git.Stat"),
				cel.UnaryBinding(func(value ref.Val) ref.Val {
					return env.CELTypeAdapter().NativeToValue(git.FileChange(value.Value().(change.Change)).Stats())
				}),
			),
		),
	)
}
