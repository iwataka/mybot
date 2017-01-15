package mybot

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestInfo(t *testing.T) {
	tmp, err := ioutil.TempFile("", "mybot")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	_, err = NewLogger(tmp.Name(), -1, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if info, _ := os.Stat(tmp.Name()); info == nil {
		t.Fatalf("%s expected to exist but not", tmp.Name())
	}
}
