package mybot

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Logger represents a logger of this app
type Logger struct {
	logger     *log.Logger
	logFile    string
	twitterAPI *TwitterAPI
	config     Config
}

// NewLogger creates an instance of Logger
func NewLogger(path string, flag int, a *TwitterAPI, c *FileConfig) (*Logger, error) {
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

// Println prints the specified values with a trailing newline
func (l *Logger) Println(v ...interface{}) error {
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
func (l *Logger) HandleError(err error) {
	if err != nil {
		l.Println(err)
	}
}

// ReadString returns a content of logging
func (l *Logger) ReadString() string {
	out, err := ioutil.ReadFile(l.logFile)
	if err != nil {
		return ""
	}
	return string(out)
}
