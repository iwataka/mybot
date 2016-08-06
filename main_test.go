package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCache(t *testing.T) {
	f, err := ioutil.TempFile("", "mybot")
	if err != nil {
		t.Fatal(err)
	}
	path := f.Name()
	unmarshalCache(path)
	marshalCache(path)
	unmarshalCache(path)
}

func TestLogger(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "mybot")
	defer os.RemoveAll(dir)

	tmp := filepath.Join(dir, "mybot-test-logger.log")
	_, err = newLogger(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(tmp); os.IsNotExist(err) || info.IsDir() {
		t.Fatalf("%s not exist.", tmp)
	}
}
