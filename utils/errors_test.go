package utils_test

import (
	"github.com/iwataka/mybot/utils"
	"github.com/stretchr/testify/require"

	"errors"
	"testing"
)

func Test_WithStack(t *testing.T) {
	var err error
	var ok bool

	err = utils.WithStack(errors.New(""))
	_, ok = err.(utils.StackTracer)
	require.True(t, ok)

	sTracer := utils.WithStack(err)
	require.Equal(t, err, sTracer)
}
