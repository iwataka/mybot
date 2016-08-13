package main

import (
	"io/ioutil"
	"log"
	"os"
)

type Logger struct {
	*log.Logger
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
	return &Logger{l, path, a, c}, nil
}

func (l *Logger) Info(msg string) {
	if l.twitterAPI != nil {
		c := l.config.Log
		if c != nil {
			err := l.twitterAPI.PostDMToAll(msg, c.AllowSelf, c.Users)
			if err != nil {
				l.Println(err)
			}
		}
	}
	l.Println(msg)
}

func (l *Logger) InfoIfError(err error) {
	if err != nil {
		l.Info(err.Error())
	}
}

func (l *Logger) FatalIfError(err error) {
	l.InfoIfError(err)
	if err != nil {
		panic(err)
	}
}

func (l *Logger) ReadString() string {
	out, err := ioutil.ReadFile(l.logFile)
	if err != nil {
		return ""
	}
	return string(out)
}
