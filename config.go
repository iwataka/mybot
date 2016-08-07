package main

import (
	"io/ioutil"
	"os"

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
		if info, _ := os.Stat("config.yml"); info != nil && !info.IsDir() {
			path = "config.yml"
		} else {
			path = "config.template.yml"
		}
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
