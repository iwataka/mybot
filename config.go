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
			Name    string
			Filter  TweetFilterConfig
			Actions []string
		}
		Searches []struct {
			Query   string
			Filter  TweetFilterConfig
			Actions []string
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
		for _, a := range account.Actions {
			if a != "retweet" && a != "favorite" {
				return errors.New(fmt.Sprintf("Invalid action: %s", a))
			}
		}
		if len(account.Actions) == 0 {
			return errors.New(fmt.Sprintf("Account %s has no action", account.Name))
		}
	}
	return nil
}
