package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type MybotConfig struct {
	GitHub *struct {
		Projects []GitHubProject
		Duration string
	} `yaml:"github"`
	Retweet *struct {
		Accounts []struct {
			Name   string
			Filter TweetFilterConfig
		}
		Searches []struct {
			Query  string
			Filter TweetFilterConfig
		}
		Notification *Notification
		Duration     string
	}
	Interaction *struct {
		Duration  string
		AllowSelf bool `yaml:"allowSelf"`
		Users     []string
	}
	Log *struct {
		AllowSelf bool `yaml:"allowSelf"`
		Users     []string
	}
	Authentication *TwitterAuth
	Option         *HTTPServer
}

func NewMybotConfig(path string) (*MybotConfig, error) {
	c := &MybotConfig{
		Option: &HTTPServer{Port: "8080"},
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
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
