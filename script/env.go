package script

import (
	"reflect"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"

	"github.com/rusinikita/changes/commit"
	"github.com/rusinikita/changes/commit/value"
	"github.com/rusinikita/changes/git"
)

func newData(commits []commit.Commit, changes []git.FileChange) Data {
	return map[string]any{
		"commits": commits,
		"changes": changes,
		"now":     time.Now(),
	}
}

func createEnv(props value.Properties) (*cel.Env, error) {
	return cel.NewEnv(
		ext.NativeTypes(
			reflect.ValueOf(commit.Commit{}),
			reflect.ValueOf(git.Signature{}),
			reflect.ValueOf(git.FileChange{}),
		),
		cel.Variable("commits", cel.ListType(cel.ObjectType("commit.Commit"))),
		cel.Variable("changes", cel.ListType(cel.ObjectType("git.FileChange"))),
		cel.Variable("now", cel.TimestampType),

		customTypeProvider(props),
	)
}

// cel.Function("method",
// cel.MemberOverload(
// "test_biba_func",
// []*cel.Type{cel.ObjectType("script.SomeStruct")},
// cel.StringType,
// cel.FunctionBinding(func(values ...ref.Val) ref.Val {
// 	return types.String(values[0].Value().(SomeStruct).Method())
// }),
// // Provide the implementation using cel.FunctionBinding()
// ),
// ),
