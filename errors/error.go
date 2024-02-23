package errors

import (
	"errors"
	"fmt"
	"strings"
)

func New(e string) error {
	return errors.New(e)
}

func Add(current, newErr error, dirs ...string) error {
	if newErr == nil {
		return current
	}

	newErr = Prefix(newErr, dirs...)

	if current == nil {
		return newErr
	}

	switch err := current.(type) {
	case *pathErr:
		return &multiErr{
			errs: append(pathErrors(newErr), *err),
		}
	case *multiErr:
		err.errs = append(err.errs, pathErrors(newErr)...)

	default:
		return &multiErr{
			errs: append(pathErrors(newErr), pathErr{
				err: err,
			}),
		}
	}

	return current
}

func Prefix(err error, dirs ...string) error {
	if err == nil {
		return err
	}

	if len(dirs) == 0 {
		return err
	}

	switch e := err.(type) {
	case *pathErr:
		e.path = append(dirs, e.path...)

	case *multiErr:
		for i, pathErr := range e.errs {
			e.errs[i].path = append(dirs, pathErr.path...)
		}

	default:
		err = &pathErr{
			path: dirs,
			err:  err,
		}
	}

	return err
}

func pathErrors(e error) []pathErr {
	switch e := e.(type) {
	case *pathErr:
		return []pathErr{*e}
	case *multiErr:
		return e.errs
	default:
		return []pathErr{{err: e}}
	}
}

type pathErr struct {
	path []string
	err  error
}

//nolint:revive
func (e *pathErr) cutRoot() (first, second string, err error) {
	switch len(e.path) {
	case 0:
		return "", "", e
	case 1:
		return e.path[0], "", e.err
	case 2:
		return e.path[0], e.path[1], e.err
	}

	return e.path[0], e.path[1], &pathErr{
		path: e.path[2:],
		err:  e.err,
	}
}

func (e *pathErr) Error() string {
	if len(e.path) == 0 {
		return e.err.Error()
	}

	return fmt.Sprintf("%s: %s", strings.Join(e.path, "/"), e.err)
}

type multiErr struct {
	errs []pathErr
}

func (e *multiErr) Error() string {
	sb := strings.Builder{}

	for i, err := range e.errs {
		sb.WriteString("- " + err.Error())

		if i != len(e.errs)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func Len(err error) int {
	if err == nil {
		return 0
	}

	if e, ok := err.(*multiErr); ok {
		return len(e.errs)
	}

	return 1
}
