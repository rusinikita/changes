package script

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/stoewer/go-strcase"

	"github.com/rusinikita/changes/commit"
	"github.com/rusinikita/changes/commit/value"
)

func customTypeProvider(props value.Properties) cel.EnvOption {
	return func(e *cel.Env) (*cel.Env, error) {
		return cel.CustomTypeProvider(&provider{
			Provider: e.CELTypeProvider(),
			props:    props,
		})(e)
	}
}

type provider struct {
	types.Provider
	props value.Properties
}

func (p *provider) FindStructFieldType(structType, fieldName string) (*types.FieldType, bool) {
	name := value.Name(fieldName)

	if _, ok := p.props[name]; !ok {
		return p.Provider.FindStructFieldType(structType, strcase.UpperCamelCase(fieldName))
	}

	fieldType := &types.FieldType{
		Type: types.StringType,
		IsSet: func(target any) bool {
			return len(target.(commit.Commit).Values[value.Name(fieldName)]) != 0
		},
		GetFrom: func(target any) (any, error) {
			return string(target.(commit.Commit).Values[value.Name(fieldName)]), nil
		},
	}

	return fieldType, true
}
