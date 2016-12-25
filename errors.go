package main

import (
	"errors"
)

type KillError error

func NewKillError(msg string) KillError {
	return errors.New(msg)
}
