package utils

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"

	"fmt"
)

// StackTracer provides a feature to get stack trace.
type StackTracer interface {
	StackTrace() errors.StackTrace
}

// WithStack return an error with the current stack trace.
func WithStack(err error) error {
	switch err.(type) {
	case StackTracer:
		return err
	default:
		return errors.WithStack(err)
	}
}

// TomlUndecodedKeysError is an error indicating that there are some undecoded
// keys in File. In other words, there are some keys which exist in File but
// not in a destination object.
type TomlUndecodedKeysError struct {
	Undecoded []toml.Key
	File      string
}

// Error returns a message of this error.
func (e *TomlUndecodedKeysError) Error() string {
	return fmt.Sprintf("%v undecoded in %s", e.Undecoded, e.File)
}
