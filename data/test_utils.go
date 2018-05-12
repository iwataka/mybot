package data

import "testing"

func NewTestFileCache(path string, t *testing.T) *FileCache {
	c, err := NewFileCache(path)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
