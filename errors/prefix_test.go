package errors

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrToPathPrefix(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{
			s:    strings.Repeat("abc .", 10),
			want: strings.Repeat("abc_.", 4),
		},
		{
			s:    strings.Repeat("abc .", 3),
			want: strings.Repeat("abc_.", 3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, StrToPathPrefix(tt.s))
		})
	}
}
