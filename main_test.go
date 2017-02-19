package main

import (
	"io/ioutil"
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
