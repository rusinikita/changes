package mock

import (
	"reflect"

	"github.com/rusinikita/changes/conf"
)

//go:generate moq -rm -stub -out conf_mock.go . Conf
type Conf = conf.Conf

func NewUnmarshal(value any, err error) *ConfMock {
	return &ConfMock{
		UnmarshalFunc: func(key string, config any, defaultValue ...any) error {
			if err != nil {
				return err
			}

			rValue := reflect.ValueOf(config)
			rValue.Elem().Set(reflect.ValueOf(value))

			return nil
		},
	}
}
