package errors

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_multiToTree(t *testing.T) {
	testErrStr := "test"
	testErr := errors.New("test")

	err := multiErr{
		errs: []pathErr{
			{
				path: gen(3),
				err:  testErr,
			},
			{
				path: gen(2),
				err:  testErr,
			},
			{
				path: gen(1),
				err:  testErr,
			},
			{
				path: gen(2),
				err:  testErr,
			},
			{
				path: gen(3),
				err:  testErr,
			},
			{
				path: []string{"1", "3"},
				err:  testErr,
			},
		},
	}

	tree := multiToTree(err)

	expected := OutputTree{
		{
			Name:     "1",
			Messages: []string{testErrStr},
			Groups: []Node{
				{
					Name: "2",
					Messages: []string{
						"3: test",
						"3: test",
						testErrStr,
						testErrStr,
					},
				},
				{
					Name: "3",
					Messages: []string{
						testErrStr,
					},
				},
			},
		},
	}

	assert.Equal(t, expected, tree)
}

func gen(n int) (ss []string) {
	for i := 1; i <= n; i++ {
		ss = append(ss, strconv.Itoa(i))
	}

	return ss
}
