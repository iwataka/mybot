package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type MybotConfig struct {
	GitHub         *GitHubConfig `toml:"github"`
	Twitter        *TwitterConfig
	Interaction    *InteractionConfig
	Log            *LogConfig
	Authentication *TwitterAuth
	HTTP           *HTTPServer `toml:"http"`
}

type GitHubConfig struct {
	Projects []GitHubProject
	Duration string
}

type TwitterConfig struct {
	Timelines    []TimelineConfig
	Searches     []SearchConfig
	Notification *Notification
	Duration     string
}

type TimelineConfig struct {
	ScreenName     *string  `toml:"screen_name"`
	ScreenNames    []string `toml:"screen_names"`
	ExcludeReplies *bool    `toml:"exclude_replies"`
	IncludeRts     *bool    `toml:"include_rts"`
	Count          *int
	Filter         *TweetFilterConfig
	Action         *TwitterAction
}

type SearchConfig struct {
	Query      *string
	Queries    []string
	ResultType *string `toml:"result_type"`
	Count      *int
	Filter     *TweetFilterConfig
	Action     *TwitterAction
}

type InteractionConfig struct {
	Duration  string
	AllowSelf bool `toml:"allow_self"`
	Users     []string
	Count     *int
}

type LogConfig struct {
	AllowSelf bool `toml:"allow_self"`
	Users     []string
}

func NewMybotConfig(path string) (*MybotConfig, error) {
	c := &MybotConfig{
		GitHub: &GitHubConfig{
			Projects: []GitHubProject{},
			Duration: "12h",
		},
		Twitter: &TwitterConfig{
			Timelines: []TimelineConfig{},
			Searches:  []SearchConfig{},
			Duration:  "1h",
		},
		HTTP: &HTTPServer{
			Host: "127.0.0.1",
			Port: "8080",
		},
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
		return nil, fmt.Errorf("%v undecoded in %s", md.Undecoded(), path)
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
			msg := fmt.Sprintf("%v has name and names properties, use `names` only.", account)
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
			msg := fmt.Sprintf("%v has query and queries properties, use `query` only.", search)
			return errors.New(msg)
		}
	}
	return nil
}
