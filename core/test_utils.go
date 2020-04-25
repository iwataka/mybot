package core

import (
	"testing"
)

func NewTestFileConfig(path string, t *testing.T) *FileConfig {
	c, err := NewFileConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
