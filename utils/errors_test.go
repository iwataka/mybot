package utils_test

import (
	. "github.com/iwataka/mybot/utils"

	"errors"
	"testing"
)

func TestWithStack(t *testing.T) {
	err := errors.New("")
	_, ok := WithStack(err).(StackTracer)
	if !ok {
		t.Fatal("WithStack should return StackTracer")
	}
}

func TestWithStackForInterruptedError(t *testing.T) {
	err := NewStreamInterruptedError()
	_, ok := WithStack(err).(StackTracer)
	if ok {
		t.Fatal("WithStack with InterruptedError should not return StackTracer")
	}
}
