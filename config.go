package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var config = &mybotConfig{}

type mybotConfig struct {
	Tweet        *tweetConfig
	Talk         *talkConfig
	UserGroup    *userGroupConfig `yaml:"userGroup"`
	Notification *notificationConfig
}

type tweetConfig struct {
	Github            []githubConfig
	Retweet           []retweetConfig
	Interval          int
	ConsumerKey       string `yaml:"consumerKey"`
	ConsumerSecret    string `yaml:"consumerSecret"`
	AccessToken       string `yaml:"accessToken"`
	AccessTokenSecret string `yaml:"accessTokenSecret"`
}

type githubConfig struct {
	User string
	Repo string
}

type retweetConfig struct {
	Name     string
	Patterns []string
	Opts     map[string]bool
}

type talkConfig struct {
	Enabled  bool
	Interval int
}

type userGroupConfig struct {
	IncludeSelf bool `yaml:"includeSelf"`
	Users       []string
}

type notificationConfig struct {
	Place bool
}

func unmarshalConfig(path string) error {
	if path == "" {
		path = "config.yml"
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return err
	}
	return nil
}
