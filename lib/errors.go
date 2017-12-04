package mybot

import (
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func WithStack(err error) error {
	switch err.(type) {
	case stackTracer:
		return err
	case *InterruptedError:
		return err
	default:
		return errors.WithStack(err)
	}
}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

func NewError(msg string) error {
	return errors.New(msg)
}

type InterruptedError struct {
	msg string
}

func NewInterruptedError() error {
	return &InterruptedError{"Interrupted"}
}

func (e InterruptedError) Error() string {
	return e.msg
}
