package mybot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCacheSave(t *testing.T) {
	dir := os.TempDir()
	fname := "cache.json"
	path := filepath.Join(dir, fname)
	c, err := NewCache(path)
	if err != nil {
		t.Fatal(err)
	}
	c.Save(path)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("%s expected to exist but not", path)
	}
}
