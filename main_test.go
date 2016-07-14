package main

import (
	"io/ioutil"
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
