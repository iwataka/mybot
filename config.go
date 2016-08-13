package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

type MybotConfig struct {
	GitHub *struct {
		Projects []GitHubProject
		Duration string
	} `yaml:"github"`
	Retweet *struct {
		Accounts     []TweetCheckConfig
		Notification *Notification
		Duration     string
	}
	Interaction *struct {
		Duration  string
		AllowSelf bool `yaml:"allowSelf"`
		Users     []string
	}
	Log            *TwitterLogConfig
	Authentication *TwitterAuth
	Option         *HTTPServer
}

func NewMybotConfig(path string) (*MybotConfig, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := &MybotConfig{
		Option: &HTTPServer{Port: "8080"},
	}
	err = yaml.Unmarshal(bytes, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *MybotConfig) Save(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0600)
	if err != nil {
		return err
	}
	if c != nil {
		data, err := json.Marshal(c)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, data, 0600)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *MybotConfig) GetReloadDuration() (time.Duration, error) {
	gd, err := time.ParseDuration(c.GitHub.Duration)
	if err != nil {
		return 0, err
	}
	rd, err := time.ParseDuration(c.Retweet.Duration)
	if err != nil {
		return 0, err
	}
	if gd < rd {
		return time.Duration(gd), nil
	}
	return time.Duration(rd), nil
}
