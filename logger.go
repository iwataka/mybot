package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Logger struct {
	logger     *log.Logger
	logFile    string
	twitterAPI *TwitterAPI
	config     *MybotConfig
}

func NewLogger(path string, flag int, a *TwitterAPI, c *MybotConfig) (*Logger, error) {
	if flag < 0 {
		flag = log.Ldate | log.Ltime | log.Lshortfile
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	l := log.New(f, "", flag)
	l.SetOutput(f)
	result := &Logger{l, path, a, c}
	return result, nil
}

func (l *Logger) Println(v ...interface{}) {
	if l.twitterAPI != nil {
		c := l.config.Log
		if c != nil {
			err := l.twitterAPI.PostDMToAll(fmt.Sprintln(v...), c.AllowSelf, c.Users)
			if err != nil {
				l.logger.Println(err)
			}
		}
	}
	l.logger.Println(v...)
}

func (l *Logger) HandleError(err error) {
	if err != nil {
		l.Println(err)
	}
}

func (l *Logger) ReadString() string {
	out, err := ioutil.ReadFile(l.logFile)
	if err != nil {
		return ""
	}
	return string(out)
}
