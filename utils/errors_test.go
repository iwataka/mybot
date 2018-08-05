package utils_test

import (
	. "github.com/iwataka/mybot/utils"
	"github.com/stretchr/testify/require"

	"errors"
	"testing"
)

func Test_WithStack(t *testing.T) {
	var err error
	var ok bool

	err = WithStack(errors.New(""))
	_, ok = err.(StackTracer)
	require.True(t, ok)

	sTracer := WithStack(err)
	require.Equal(t, err, sTracer)

	err = NewStreamInterruptedError()
	_, ok = WithStack(err).(StackTracer)
	require.False(t, ok)
}
