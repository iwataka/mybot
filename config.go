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
		Timelines []struct {
			ScreenName  *string  `toml:"screen_name"`
			ScreenNames []string `toml:"screen_names"`
			Count       *int
			Filter      *TweetFilterConfig
			Action      *TwitterAction
		}
		Searches []struct {
			Query   *string
			Queries []string
			Count   *int
			Filter  *TweetFilterConfig
			Action  *TwitterAction
		}
		Notification *Notification
		Duration     string
	}
	Interaction *struct {
		Duration  string
		AllowSelf bool `toml:"allow_self"`
		Users     []string
	}
	Log *struct {
		AllowSelf bool `toml:"allow_self"`
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
	for _, account := range config.Twitter.Timelines {
		if account.Action == nil {
			msg := fmt.Sprintf("%v has no action", account)
			return errors.New(msg)
		}
		if account.ScreenName == nil && (account.ScreenNames == nil || len(account.ScreenNames) == 0) {
			msg := fmt.Sprintf("%v has no name", account)
			return errors.New(msg)
		}
		if account.ScreenName != nil && account.ScreenNames != nil && len(account.ScreenNames) != 0 {
			msg := fmt.Sprintf("%v has name and names properties, use `names` only.")
			return errors.New(msg)
		}
	}
	for _, search := range config.Twitter.Searches {
		if search.Action == nil {
			msg := fmt.Sprintf("%v has no action", search)
			return errors.New(msg)
		}
		if search.Query == nil && (search.Queries == nil || len(search.Queries) == 0) {
			msg := fmt.Sprintf("%v has no query", search)
			return errors.New(msg)
		}
		if search.Query != nil && search.Queries != nil && len(search.Queries) != 0 {
			msg := fmt.Sprintf("%v has query and queries properties, use `query` only.")
			return errors.New(msg)
		}
	}
	return nil
}
