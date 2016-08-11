package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var config = &mybotConfig{}

type mybotConfig struct {
	GitHub         *githubConfig `yaml:"github"`
	Retweet        *retweetConfig
	Interaction    *interactionConfig
	Log            *logConfig
	Authentication *authenticationConfig
	Option         *optionConfig
}

type githubConfig struct {
	Projects []projectConfig
	Duration string
}

type projectConfig struct {
	User string
	Repo string
}

type retweetConfig struct {
	Accounts     []accountConfig
	Notification *notificationConfig
	Duration     string
}

type accountConfig struct {
	Name     string
	Patterns []string
	Opts     map[string]bool
}

type notificationConfig struct {
	Place *placeConfig
}

type placeConfig struct {
	AllowSelf bool `yaml:"allowSelf"`
	Users     []string
}

type interactionConfig struct {
	Duration  string
	AllowSelf bool `yaml:"allowSelf"`
	Users     []string
}

type logConfig struct {
	AllowSelf bool `yaml:"allowSelf"`
	Users     []string
}

type authenticationConfig struct {
	ConsumerKey       string `yaml:"consumerKey"`
	ConsumerSecret    string `yaml:"consumerSecret"`
	AccessToken       string `yaml:"accessToken"`
	AccessTokenSecret string `yaml:"accessTokenSecret"`
}

type optionConfig struct {
	Name string
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
