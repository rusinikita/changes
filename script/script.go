package script

import (
	"fmt"
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"

	"github.com/rusinikita/changes/change"
	"github.com/rusinikita/changes/commit"
	"github.com/rusinikita/changes/errors"
)

type Data map[string]any

type Config struct {
	Message string `mapstructure:"message"`
	Func    string `mapstructure:"func"`
}

type Script struct {
	cel.Program
	message string
}

func New(c Config, env *cel.Env) (*Script, error) {
	if c.Func == "" || c.Message == "" {
		return nil, errors.New("message and func are required")
	}

	ast, issues := env.Compile(c.Func)
	if err := issues.Err(); err != nil {
		return nil, errors.Prefix(issues.Err(), "func")
	}

	if err := checkOutType(ast.OutputType()); err != nil {
		return nil, errors.Prefix(err, "func")
	}

	prg, err := env.Program(ast)
	if err != nil {
		return nil, errors.Prefix(err, "func")
	}

	script := &Script{
		Program: prg,
		message: c.Message,
	}

	return script, nil
}

func (s *Script) Run(data Data) (any, error) {
	out, _, err := s.Eval(map[string]any(data))
	if err != nil {
		return nil, err
	}

	switch v := out.Value().(type) {
	case bool:
		return v, nil
	case []ref.Val:
		return convertList(out, v)
	case []change.Change:
		return v, nil
	case []commit.Commit:
		return v, nil
	default:
		return nil, fmt.Errorf("unsuported return `%s` `%T`", v, v)
	}
}

func convertList(out ref.Val, list []ref.Val) (any, error) {
	if len(list) == 0 {
		return nil, nil
	}

	switch item := list[0].Value().(type) {
	case change.Change:
		return out.ConvertToNative(reflect.TypeOf([]change.Change{}))
	case commit.Commit:
		return out.ConvertToNative(reflect.TypeOf([]commit.Commit{}))
	default:
		return nil, fmt.Errorf("unsuported return `%s` `%T list`", out.Value(), item)
	}
}

func checkOutType(t *cel.Type) error {
	if t.Kind() == cel.BoolKind {
		return nil
	}

	if t.Kind() != cel.ListKind {
		return fmt.Errorf("only bool and list return supported, but `%s`", t.TypeName())
	}

	if t.IsExactType(commitListType) || t.IsExactType(changesListType) {
		return nil
	}

	return fmt.Errorf("only bool and list return supported, but `%s`", t.TypeName())
}
