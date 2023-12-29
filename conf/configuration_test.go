package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_conf_Unmarshal(t *testing.T) {
	c := New()

	t.Run("error", func(t *testing.T) {
		err := c.Unmarshal("some.key", "", "")

		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		value := map[string]string{
			"test":      "test",
			"more_test": "more_test",
		}

		str := struct {
			Test     string
			MoreTest string `mapstructure:"more_test"`
		}{}

		err := c.Unmarshal("some.key", &str, value)

		assert.NoError(t, err)
		assert.Equal(t, "test", str.Test)
		assert.Equal(t, "more_test", str.MoreTest)
	})
}
