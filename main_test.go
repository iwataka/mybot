package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

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
	_, err = mybot.NewLogger(tmp, -1, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(tmp); os.IsNotExist(err) || info.IsDir() {
		t.Fatalf("%s not exist.", tmp)
	}
}

func TestMonitorFile(t *testing.T) {
	file, err := ioutil.TempFile("", "mybot")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	ch := make(chan bool)
	defer close(ch)
	dur := time.Second / 10
	go monitorFile(file.Name(), dur, func() {
		ch <- true
	})
	time.Sleep(dur)
	_, err = file.WriteString("foo")
	if err != nil {
		t.Fatal(err)
	}
	<-ch
}
