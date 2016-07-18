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

	tmpFile, err := ioutil.TempFile("", "mybot")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	_defaultOutput := defaultOutput
	defaultOutput = tmpFile
	defer func() {
		defaultOutput = _defaultOutput
	}()
	logger, err := newLogger("")
	if err != nil {
		t.Fatal(err)
	}
	logger.Print("foo")
	out, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if string(out) == "" {
		t.Fatalf("%s is empty", tmpFile.Name())
	}
}
