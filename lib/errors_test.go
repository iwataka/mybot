package mybot

import (
	"errors"
	"testing"
)

func TestWithStack(t *testing.T) {
	err := errors.New("")
	_, ok := WithStack(err).(stackTracer)
	if !ok {
		t.Fatal("WithStack should return stackTracer")
	}
}

func TestErrorf(t *testing.T) {
	_, ok := Errorf("").(stackTracer)
	if !ok {
		t.Fatal("NewError should return stackTracer")
	}
}

func TestNewError(t *testing.T) {
	_, ok := NewError("").(stackTracer)
	if !ok {
		t.Fatal("NewError should return stackTracer")
	}
}
