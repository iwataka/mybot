package mybot

import (
	"testing"
)

func NewTestFileCache(path string, t *testing.T) *FileCache {
	c, err := NewFileCache(path)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func NewTestFileConfig(path string, t *testing.T) *FileConfig {
	c, err := NewFileConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
