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
	return c, nil
}
