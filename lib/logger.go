package mybot

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Logger interface {
	Println(v ...interface{}) error
	HandleError(err error)
	ReadString() string
}

// TwitterLogger represents a logger of this app
type TwitterLogger struct {
	logger     *log.Logger
	logFile    string
	twitterAPI *TwitterAPI
	config     Config
}

// NewTwitterLogger creates an instance of Logger
func NewTwitterLogger(path string, flag int, a *TwitterAPI, c Config) (*TwitterLogger, error) {
	if flag < 0 {
		flag = log.Ldate | log.Ltime | log.Lshortfile
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	l := log.New(f, "", flag)
	l.SetOutput(f)
	result := &TwitterLogger{l, path, a, c}
	return result, nil
}

// Println prints the specified values with a trailing newline
func (l *TwitterLogger) Println(v ...interface{}) error {
	if l.twitterAPI != nil {
		c, err := l.config.GetLog()
		if err != nil {
			return err
		}
		if c != nil {
			err := l.twitterAPI.PostDMToAll(fmt.Sprintln(v...), c.AllowSelf, c.Users)
			if err != nil {
				l.logger.Println(err)
			}
		}
	}
	l.logger.Println(v...)
	return nil
}

// HandleError prints the specified error if it is not nil
func (l *TwitterLogger) HandleError(err error) {
	if err != nil {
		l.Println(err)
	}
}

// ReadString returns a content of logging
func (l *TwitterLogger) ReadString() string {
	out, err := ioutil.ReadFile(l.logFile)
	if err != nil {
		return ""
	}
	return string(out)
}
