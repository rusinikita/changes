package value

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rusinikita/changes/conf/mock"
	"github.com/rusinikita/changes/errors"
)

func Test_validateName(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		wantErr string
	}{
		{
			name: "success",
			args: "abc",
		},
		{
			name: "success",
			args: strings.Repeat("a", 20),
		},
		{
			name:    "fail: too short",
			args:    "ab",
			wantErr: ErrorNameLen,
		},
		{
			name:    "fail: too long",
			args:    strings.Repeat("a", 21),
			wantErr: ErrorNameLen,
		},
		{
			name:    "fail: banned char",
			args:    "qwerty123abc",
			wantErr: ErrorNameChars,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateName(tt.args)
			if tt.wantErr == "" {
				assert.Nil(t, err)
			}

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

func Test_parseValue(t *testing.T) {
	tests := []struct {
		name    string
		c       propertyConf
		want    Property
		wantErr string
	}{
		{
			name:    "fail:empty",
			c:       propertyConf{},
			wantErr: ErrorValueEmpty,
		},
		{
			name: "fail:both",
			c: propertyConf{
				Regexp:  "123xyz",
				Allowed: []string{"test"},
			},
			wantErr: ErrorValueBoth,
		},
		{
			name: "fail:regex",
			c: propertyConf{
				Regexp: "123xyz[",
				MaxLen: 10,
			},
			wantErr: "missing closing ]",
		},
		{
			name: "ok",
			c: propertyConf{
				Allowed: []string{"test"},
				MaxLen:  1,
			},
			want: Property{
				Allowed: map[string]struct{}{
					"test": {},
				},
				MaxLen: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseValue(tt.c)
			if tt.wantErr != "" {
				require.NotNil(t, err)
				assert.ErrorContains(t, err, tt.wantErr)
				assert.Zero(t, got)

				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRegexp(t *testing.T) {
	exp := regexp.MustCompile("^[a-z-_]+$")

	tests := []string{
		"test",
		"qwerty123test",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			t.Log(exp.FindAllString(test, -1))
		})
	}
}

func TestParse(t *testing.T) {
	t.Run("unmarshall err", func(t *testing.T) {
		m := newUnmarshal(nil, errors.New("test"))

		v, err := Get(m)

		assert.Nil(t, v)
		assert.Error(t, err)
	})

	t.Run("default values", func(t *testing.T) {
		m := newUnmarshal(map[string]propertyConf{}, nil)

		v, err := Get(m)

		assert.NoError(t, err)
		assert.Equal(t, DefaultProperties, v)
	})
}

func newUnmarshal(value any, err error) *mock.ConfMock {
	return &mock.ConfMock{
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
