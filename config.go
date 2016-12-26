package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// MybotConfig is a root of the all configurations.
type MybotConfig struct {
	GitHub      *GitHubConfig      `toml:"github"`
	Twitter     *TwitterConfig     `toml:"twitter"`
	DB          *DBConfig          `toml:"db"`
	Interaction *InteractionConfig `toml:"interaction"`
	Log         *LogConfig         `toml:"log"`
	HTTP        *HTTPConfig        `toml:"http"`
	source      string             `toml:"-"`
}

// GitHubConfig is a configuration of GitHub projects
type GitHubConfig struct {
	Projects []GitHubProject `toml:"projects,omitempty"`
	Duration string          `toml:"duration"`
}

// SourceConfig is a configuration for common sources
type SourceConfig struct {
	Filter *TweetFilterConfig `toml:"filter"`
	Action *TwitterAction     `toml:"action"`
}

// TwitterConfig is a configuration for Twitter
type TwitterConfig struct {
	Timelines    []TimelineConfig `toml:"timelines,omitempty"`
	Favorites    []FavoriteConfig `toml:"favorites,omitempty"`
	Searches     []SearchConfig   `toml:"searches,omitempty"`
	Notification *Notification    `toml:"notification"`
	Duration     string           `toml:"duration"`
	Debug        bool             `toml:"debug"`
}

// TimelineConfig is a configuration for Twitter timelines
type TimelineConfig struct {
	*SourceConfig
	ScreenNames    []string `toml:"screen_names"`
	ExcludeReplies *bool    `toml:"exclude_replies"`
	IncludeRts     *bool    `toml:"include_rts"`
	Count          int      `toml:"count"`
}

func NewTimelineConfig() *TimelineConfig {
	return &TimelineConfig{
		SourceConfig: &SourceConfig{
			Filter: &TweetFilterConfig{
				Vision: &VisionCondition{
					Face: &VisionFaceCondition{},
				},
			},
			Action: &TwitterAction{},
		},
	}
}

// FavoriteConfig is a configuration for Twitter favorites
type FavoriteConfig struct {
	*SourceConfig
	ScreenNames []string `toml:"screen_names"`
	Count       int      `toml:"count"`
}

func NewFavoriteConfig() *FavoriteConfig {
	return &FavoriteConfig{
		SourceConfig: &SourceConfig{
			Filter: &TweetFilterConfig{
				Vision: &VisionCondition{
					Face: &VisionFaceCondition{},
				},
			},
			Action: &TwitterAction{},
		},
	}
}

// SearchConfig is a configuration for Twitter searches
type SearchConfig struct {
	*SourceConfig
	Queries    []string `toml:"queries"`
	ResultType string   `toml:"result_type,omitempty"`
	Count      int      `toml:"count"`
}

func NewSearchConfig() *SearchConfig {
	return &SearchConfig{
		SourceConfig: &SourceConfig{
			Filter: &TweetFilterConfig{
				Vision: &VisionCondition{
					Face: &VisionFaceCondition{},
				},
			},
			Action: &TwitterAction{},
		},
	}
}

type DBConfig struct {
	Driver      string `toml:"driver,omitempty"`
	DataSource  string `toml:"data_source,omitempty"`
	VisionTable string `toml:"vision_table,omitempty"`
}

// InteractionConfig is a configuration for interaction through Twitter direct
// message
type InteractionConfig struct {
	Duration  string   `toml:"duration"`
	AllowSelf bool     `toml:"allow_self"`
	Users     []string `toml:"users,omitempty"`
	Count     int      `toml:"count"`
}

// LogConfig is a configuration for logging
type LogConfig struct {
	AllowSelf bool     `toml:"allow_self"`
	Users     []string `toml:"users,omitempty"`
}

type HTTPConfig struct {
	Name     string `toml:"name"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Enabled  bool   `toml:"enabled"`
	LogLines int    `toml:"log_lines,omitempty"`
}

// GetScreenNames returns all screen names in the configuration
func (tc *TwitterConfig) GetScreenNames() []string {
	result := []string{}
	for _, t := range tc.Timelines {
		result = append(result, t.ScreenNames...)
	}
	for _, f := range tc.Favorites {
		result = append(result, f.ScreenNames...)
	}
	return result
}

// NewMybotConfig takes the configuration file path and returns a configuration
// instance.
func NewMybotConfig(path string, vision *VisionAPI) (*MybotConfig, error) {
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
		HTTP: &HTTPConfig{
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
	err = ValidateConfig(c)
	if err != nil {
		return nil, err
	}
	for _, t := range c.Twitter.Timelines {
		if t.Filter.Vision == nil {
			t.Filter.Vision = new(VisionCondition)
		}
		if t.Filter.Vision.Face == nil {
			t.Filter.Vision.Face = new(VisionFaceCondition)
		}
		t.Filter.visionAPI = vision
	}
	for _, f := range c.Twitter.Favorites {
		if f.Filter.Vision == nil {
			f.Filter.Vision = new(VisionCondition)
		}
		if f.Filter.Vision.Face == nil {
			f.Filter.Vision.Face = new(VisionFaceCondition)
		}
		f.Filter.visionAPI = vision
	}
	for _, s := range c.Twitter.Searches {
		if s.Filter.Vision == nil {
			s.Filter.Vision = new(VisionCondition)
		}
		if s.Filter.Vision.Face == nil {
			s.Filter.Vision.Face = new(VisionFaceCondition)
		}
		s.Filter.visionAPI = vision
	}
	c.source = path
	return c, nil
}

func (c *MybotConfig) SetVisionoAPI(vision *VisionAPI) {
	for _, t := range c.Twitter.Timelines {
		t.Filter.visionAPI = vision
	}
	for _, f := range c.Twitter.Favorites {
		f.Filter.visionAPI = vision
	}
	for _, s := range c.Twitter.Searches {
		s.Filter.visionAPI = vision
	}
}

func ValidateConfig(config *MybotConfig) error {
	for _, account := range config.Twitter.Timelines {
		if account.Action == nil {
			msg := fmt.Sprintf("%v has no action", account)
			return errors.New(msg)
		}
		if len(account.ScreenNames) == 0 {
			msg := fmt.Sprintf("%v has no name", account)
			return errors.New(msg)
		}
		filter := account.Filter
		if (filter.Vision != nil && !filter.Vision.isEmpty()) &&
			(filter.RetweetedThreshold > 0 || filter.FavoriteThreshold > 0) {
			bytes, _ := json.Marshal(account)
			msg := fmt.Sprintf("%s\n%s",
				"Don't use both of Vision API and retweeted/favorite threshold",
				string(bytes),
			)
			return errors.New(msg)
		}
	}
	for _, favorite := range config.Twitter.Favorites {
		if favorite.Action == nil {
			msg := fmt.Sprintf("%v has no action", favorite)
			return errors.New(msg)
		}
		if len(favorite.ScreenNames) == 0 {
			msg := fmt.Sprintf("%v has no name", favorite)
			return errors.New(msg)
		}
		filter := favorite.Filter
		if (filter.Vision != nil && !filter.Vision.isEmpty()) &&
			(filter.RetweetedThreshold > 0 || filter.FavoriteThreshold > 0) {
			bytes, _ := json.Marshal(favorite)
			msg := fmt.Sprintf("%s\n%s",
				"Don't use both of Vision API and retweeted/favorite threshold",
				string(bytes),
			)
			return errors.New(msg)
		}
	}
	for _, search := range config.Twitter.Searches {
		if search.Action == nil {
			msg := fmt.Sprintf("%v has no action", search)
			return errors.New(msg)
		}
		if len(search.Queries) == 0 {
			msg := fmt.Sprintf("%v has no query", search)
			return errors.New(msg)
		}
		filter := search.Filter
		if (filter.Vision != nil && !filter.Vision.isEmpty()) &&
			(filter.RetweetedThreshold > 0 || filter.FavoriteThreshold > 0) {
			bytes, _ := json.Marshal(search)
			msg := fmt.Sprintf("%s\n%s",
				"Don't use both of Vision API and retweeted/favorite threshold",
				string(bytes),
			)
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

// Save saves the config data to the specified file
func (c *MybotConfig) Save() error {
	err := os.MkdirAll(filepath.Dir(c.source), 0600)
	if err != nil {
		return err
	}
	if c != nil {
		writer := new(bytes.Buffer)
		enc := toml.NewEncoder(writer)
		err := enc.Encode(c)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(c.source, writer.Bytes(), 0600)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *MybotConfig) Reload() error {
	bytes, err := ioutil.ReadFile(c.source)
	if err != nil {
		return err
	}
	md, err := toml.Decode(string(bytes), c)
	if err != nil {
		return err
	}
	if len(md.Undecoded()) != 0 {
		return fmt.Errorf("%v undecoded in %s", md.Undecoded(), c.source)
	}
	return nil
}
