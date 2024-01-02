package script

import (
	"fmt"
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"

	"github.com/rusinikita/changes/commit"
	"github.com/rusinikita/changes/errors"
	"github.com/rusinikita/changes/git"
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

	prg, err := env.Program(ast)

	script := &Script{
		Program: prg,
		message: c.Message,
	}

	return script, errors.Prefix(err, "func")
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
	default:
		return nil, fmt.Errorf("only bool and list return supported, but `%T`", v)
	}
}

func convertList(out ref.Val, list []ref.Val) (any, error) {
	if len(list) == 0 {
		return nil, nil
	}

	switch item := list[0].Value().(type) {
	case git.FileChange:
		return out.ConvertToNative(reflect.TypeOf([]git.FileChange{}))
	case commit.Commit:
		return out.ConvertToNative(reflect.TypeOf([]commit.Commit{}))
	default:
		return nil, fmt.Errorf("only commits and changes return list supported, but `%T`", item)
	}
}
