package utils_test

import (
	. "github.com/iwataka/mybot/utils"
	"github.com/stretchr/testify/assert"

	"errors"
	"testing"
)

func TestWithStack(t *testing.T) {
	var err error
	var ok bool

	err = WithStack(errors.New(""))
	_, ok = err.(StackTracer)
	assert.True(t, ok)

	sTracer := WithStack(err)
	assert.Equal(t, err, sTracer)

	err = NewStreamInterruptedError()
	_, ok = WithStack(err).(StackTracer)
	assert.False(t, ok)
}
