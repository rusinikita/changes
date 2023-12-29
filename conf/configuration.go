package conf

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Conf interface {
	Unmarshal(key string, config any, defaultValue ...any) error
}

type conf struct {
	*viper.Viper
	val *validator.Validate
}

func New() Conf {
	return &conf{
		Viper: viper.New(),
		val:   validator.New(),
	}
}

func (c *conf) Unmarshal(key string, config any, defaultValue ...any) (err error) {
	if len(defaultValue) > 0 {
		c.SetDefault(key, defaultValue[0])
	}

	err = c.UnmarshalKey(key, config)
	if err != nil {
		return fmt.Errorf("failed to unmarshall `%s`: %w", key, err)
	}

	return nil
}
