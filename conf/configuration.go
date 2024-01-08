package conf

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Conf interface {
	GetString(key string, defaultValue ...string) string
	Unmarshal(key string, config any, defaultValue ...any) error
}

type conf struct {
	*viper.Viper
	val *validator.Validate
}

func New(cfgFile string) (Conf, error) {
	v := viper.New()

	v.SetConfigName(".changes")
	v.AddConfigPath(".")
	v.SetConfigFile(cfgFile)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return &conf{
		Viper: v,
		val:   validator.New(),
	}, nil
}

func (c *conf) GetString(key string, defaultValue ...string) string {
	if len(defaultValue) > 0 {
		c.SetDefault(key, defaultValue[0])
	}

	return c.Viper.GetString(key)
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
