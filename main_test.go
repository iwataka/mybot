package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/iwataka/mybot/lib"
)

func TestCache(t *testing.T) {
	f, err := ioutil.TempFile("", "mybot")
	if err != nil {
		t.Fatal(err)
	}
	path := f.Name()
	c, err := mybot.NewFileCache(path)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Save()
	if err != nil {
		t.Fatal(err)
	}
	c, err = mybot.NewFileCache(path)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogger(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "mybot")
	defer os.RemoveAll(dir)

	tmp := filepath.Join(dir, "mybot-test-logger.log")
	_, err = mybot.NewTwitterLogger(tmp, -1, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(tmp); os.IsNotExist(err) || info.IsDir() {
		t.Fatalf("%s not exist.", tmp)
	}
}
