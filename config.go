package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type MybotConfig struct {
	GitHub *struct {
		Projects []GitHubProject
		Duration string
	} `toml:"github"`
	Twitter *struct {
		Accounts []struct {
			Name   string
			Filter *TweetFilterConfig
			Action *TwitterAction
		}
		Searches []struct {
			Query  string
			Filter *TweetFilterConfig
			Action *TwitterAction
		}
		Notification *Notification
		Duration     string
	}
	Interaction *struct {
		Duration  string
		AllowSelf bool `toml:"allowSelf"`
		Users     []string
	}
	Log *struct {
		AllowSelf bool `toml:"allowSelf"`
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
	md, err := toml.Decode(string(bytes), c)
	if err != nil {
		return nil, err
	}
	if len(md.Undecoded()) != 0 {
		return nil, errors.New(fmt.Sprintf("%v undecoded in %s", md.Undecoded(), path))
	}
	err = validateConfig(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func validateConfig(config *MybotConfig) error {
	for _, account := range config.Twitter.Accounts {
		if account.Action == nil {
			return errors.New(fmt.Sprintf("Account %s has no action", account.Name))
		}
	}
	return nil
}
