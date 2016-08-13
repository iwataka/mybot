package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type MultiLogger struct {
	*log.Logger
	logFile string
	loggers []Logger
}

type Logger func(string) error

func NewLogger(path string, prefix string, flag int, ls []Logger) (*MultiLogger, error) {
	if flag < 0 {
		flag = log.Ldate | log.Ltime | log.Lshortfile
	}
	if info, err := os.Stat(path); os.IsExist(err) && info.IsDir() {
		path = filepath.Join(path, ".mybot-debug.log")
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	l := log.New(f, prefix, flag)
	return &MultiLogger{l, path, ls}, nil
}

func (l *MultiLogger) Info(msg string) {
	if l.loggers != nil {
		for _, logger := range l.loggers {
			err := logger(msg)
			if err != nil {
				l.Println(err)
			}
		}
	}
	l.Println(msg)
}

func (l *MultiLogger) InfoIfError(err error) {
	if err != nil {
		l.Info(err.Error())
	}
}

func (l *MultiLogger) FatalIfError(err error) {
	l.InfoIfError(err)
	if err != nil {
		panic(err)
	}
}

func (l *MultiLogger) ReadString() string {
	out, err := ioutil.ReadFile(l.logFile)
	if err != nil {
		return ""
	}
	return string(out)
}
