package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// MybotConfig is a root of the all configurations.
type MybotConfig struct {
	GitHub         *GitHubConfig      `toml:"github"`
	Twitter        *TwitterConfig     `toml:"twitter"`
	Interaction    *InteractionConfig `toml:"interaction"`
	Log            *LogConfig         `toml:"log"`
	Authentication *TwitterAuth       `toml:"authentication"`
	HTTP           *HTTPServer        `toml:"http"`
}

type GitHubConfig struct {
	Projects []GitHubProject `toml:"projects"`
	Duration string          `toml:"duration"`
}

type SourceConfig struct {
	Filter *TweetFilterConfig `toml:"filter"`
	Action *TwitterAction     `toml:"action"`
}

type TwitterConfig struct {
	Timelines    []TimelineConfig `toml:"timelines"`
	Favorites    []FavoriteConfig `toml:"favorites"`
	Searches     []SearchConfig   `toml:"searches"`
	Notification *Notification    `toml:"notification"`
	Duration     string           `toml:"duration"`
}

type TimelineConfig struct {
	*SourceConfig
	ScreenName     *string  `toml:"screen_name"`
	ScreenNames    []string `toml:"screen_names"`
	ExcludeReplies *bool    `toml:"exclude_replies"`
	IncludeRts     *bool    `toml:"include_rts"`
	Count          *int     `toml:"count"`
}

type FavoriteConfig struct {
	*SourceConfig
	ScreenName  *string  `toml:"screen_name"`
	ScreenNames []string `toml:"screen_names"`
	Count       *int     `toml:"count"`
}

type SearchConfig struct {
	*SourceConfig
	Query      *string  `toml:"query"`
	Queries    []string `toml:"queries"`
	ResultType *string  `toml:"result_type"`
	Count      *int     `toml:"count"`
}

type InteractionConfig struct {
	Duration  string   `toml:"duration"`
	AllowSelf bool     `toml:"allow_self"`
	Users     []string `toml:"users"`
	Count     *int     `toml:"count"`
}

type LogConfig struct {
	AllowSelf bool     `toml:"allow_self"`
	Users     []string `toml:"users"`
}

// NewMybotConfig takes the configuration file path and returns a configuration
// instance.
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
		filter := account.Filter
		if filter.Vision != nil && (filter.RetweetedThreshold != nil || filter.FavoriteThreshold != nil) {
			msg := "Don't use Vision API and retweeted/favorite threshold"
			return errors.New(msg)
		}
	}
	for _, favorite := range config.Twitter.Favorites {
		if favorite.Action == nil {
			msg := fmt.Sprintf("%v has no action", favorite)
			return errors.New(msg)
		}
		if favorite.ScreenName == nil && (favorite.ScreenNames == nil || len(favorite.ScreenNames) == 0) {
			msg := fmt.Sprintf("%v has no name", favorite)
			return errors.New(msg)
		}
		if favorite.ScreenName != nil && favorite.ScreenNames != nil && len(favorite.ScreenNames) != 0 {
			msg := fmt.Sprintf("%v has name and names properties, use `names` only.", favorite)
			return errors.New(msg)
		}
		filter := favorite.Filter
		if filter.Vision != nil && (filter.RetweetedThreshold != nil || filter.FavoriteThreshold != nil) {
			msg := "Don't use Vision API and retweeted/favorite threshold"
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
		filter := search.Filter
		if filter.Vision != nil && (filter.RetweetedThreshold != nil || filter.FavoriteThreshold != nil) {
			msg := "Don't use Vision API and retweeted/favorite threshold"
			return errors.New(msg)
		}
	}
	return nil
}

// TomlText returns a toml text encoded from the configuration. If error occurs
// while encoding, this returns an empty string.
func (c *MybotConfig) TomlText(indent string) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := toml.NewEncoder(buf)
	enc.Indent = indent
	err := enc.Encode(c)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}
