package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/iwataka/mybot/src"
)

func TestCache(t *testing.T) {
	f, err := ioutil.TempFile("", "mybot")
	if err != nil {
		t.Fatal(err)
	}
	path := f.Name()
	c, err := mybot.NewMybotCache(path)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Save(path)
	if err != nil {
		t.Fatal(err)
	}
	c, err = mybot.NewMybotCache(path)
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
	defer os.Remove(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	modified := false
	ch := make(chan bool)
	dur := time.Second / 10.0
	go monitorFile(file.Name(), dur, func() {
		modified = true
	})
	var e string = ""
	go func() {
		time.Sleep(dur * 3)
		if modified {
			e = fmt.Sprintf("%s is not modified", file.Name())
			ch <- true
			return
		}
		_, err = file.WriteString("foo")
		if err != nil {
			e = err.Error()
			ch <- true
			return
		}
		time.Sleep(dur * 3)
		if !modified {
			e = fmt.Sprintf("%s is now modified", file.Name())
			ch <- true
			return
		}
		modified = false
		time.Sleep(dur * 3)
		if modified {
			e = fmt.Sprintf("%s is not modified", file.Name())
			ch <- true
			return
		}
		ch <- true
	}()
	<-ch
	if e != "" {
		t.Fatalf(e)
	}
}
