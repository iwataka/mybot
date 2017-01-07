package mybot

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

// MybotConfig is a root of the all configurations of this applciation.
type MybotConfig struct {
	// Twitter is a configuration related to Twitter.
	Twitter *TwitterConfig `toml:"twitter"`
	// DB is a configuration related to DB.
	DB *DBConfig `toml:"db"`
	// Interaction is a configuration related to interaction with users
	// such as Twitter's direct message exchange.
	Interaction *InteractionConfig `toml:"interaction"`
	// Log is a configuration related to logging.
	Log *LogConfig `toml:"log"`
	// Server is a server configuration.
	Server *ServerConfig `toml:"server"`
	// source is a configuration file from which this was loaded. This is
	// needed to save the content to the same file.
	source string `toml:"-"`
}

// NewMybotConfig takes the configuration file path and returns a configuration
// instance.
func NewMybotConfig(path string, vision *VisionAPI) (*MybotConfig, error) {
	c := &MybotConfig{
		Twitter: &TwitterConfig{
			Timelines: []TimelineConfig{},
			Searches:  []SearchConfig{},
			Duration:  "1h",
			Notification: &Notification{
				Place: &PlaceNotification{},
			},
		},
		Server: &ServerConfig{
			Host: "localhost",
			Port: "3256",
		},
		DB:          &DBConfig{},
		Log:         &LogConfig{},
		Interaction: &InteractionConfig{},
	}

	c.source = path
	err := c.Load()
	if err != nil {
		return nil, err
	}

	err = c.Validate()
	if err != nil {
		return nil, err
	}

	// Assign empty values to config instance to prevent nil pointer
	// reference error.
	for _, t := range c.Twitter.Timelines {
		if t.Filter.Vision == nil {
			t.Filter.Vision = new(VisionCondition)
		}
		if t.Filter.Vision.Face == nil {
			t.Filter.Vision.Face = new(VisionFaceCondition)
		}
	}
	for _, f := range c.Twitter.Favorites {
		if f.Filter.Vision == nil {
			f.Filter.Vision = new(VisionCondition)
		}
		if f.Filter.Vision.Face == nil {
			f.Filter.Vision.Face = new(VisionFaceCondition)
		}
	}
	for _, s := range c.Twitter.Searches {
		if s.Filter.Vision == nil {
			s.Filter.Vision = new(VisionCondition)
		}
		if s.Filter.Vision.Face == nil {
			s.Filter.Vision.Face = new(VisionFaceCondition)
		}
	}
	return c, nil
}

// Validate tries to validate the specified configuration. If invalid values
// are detected, this returns an error.
func (c *MybotConfig) Validate() error {
	for _, account := range c.Twitter.Timelines {
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

	for _, favorite := range c.Twitter.Favorites {
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

	for _, search := range c.Twitter.Searches {
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

// Read returns a configuration content as a toml text. If error occurs while
// encoding, this returns an empty string. This return value is not same as the
// source file's content.
func (c *MybotConfig) Read(indent string) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := toml.NewEncoder(buf)
	enc.Indent = indent
	err := enc.Encode(c)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// Save saves the specified configuration to the source file.
func (c *MybotConfig) Save() error {
	// Make a directory before all.
	err := os.MkdirAll(filepath.Dir(c.source), 0751)
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
		err = ioutil.WriteFile(c.source, writer.Bytes(), 0640)
		if err != nil {
			return err
		}
	}
	return nil
}

// Load loads the configuration from the source file. If the specified source
// file doesn't exist, this method does nothing and returns nil.
func (c *MybotConfig) Load() error {
	if info, err := os.Stat(c.source); err == nil && !info.IsDir() {
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
	}
	return nil
}

// SourceConfig is a configuration for common data sources such as Twitter's
// timelines, favorites and searches. Sources should have filters and actions.
type SourceConfig struct {
	// Filter filters out incoming data from sources.
	Filter *TweetFilterConfig `toml:"filter"`
	// Action defines actions for data passing through filters.
	Action *TwitterAction `toml:"action"`
}

// TwitterConfig is a configuration related to Twitter.
type TwitterConfig struct {
	Timelines []TimelineConfig `toml:"timelines,omitempty"`
	Favorites []FavoriteConfig `toml:"favorites,omitempty"`
	Searches  []SearchConfig   `toml:"searches,omitempty"`
	// Notification is a configuration related to notification for users.
	// Currently only place notification is supported, which means that
	// when a tweet with place information is detected, it is notified to
	// the specified users.
	Notification *Notification `toml:"notification"`
	// Duration is a duration for some periodic jobs such as fetching
	// users' favorites and searching by the specified condition.
	Duration string `toml:"duration"`
	// Debug is a flag for debugging, if it is true, additional information
	// is outputted.
	Debug bool `toml:"debug"`
}

// GetScreenNames returns all screen names in the TwitterConfig instance. This
// is useful to do something for all related users.
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

// TimelineConfig is a configuration for Twitter timelines
type TimelineConfig struct {
	*SourceConfig
	ScreenNames    []string `toml:"screen_names"`
	ExcludeReplies *bool    `toml:"exclude_replies"`
	IncludeRts     *bool    `toml:"include_rts"`
	Count          int      `toml:"count"`
}

// NewTimelineConfig returns TimelineConfig instance, which is empty but has a
// non-nil filter and action.
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

// NewFavoriteCnfig returns FavoriteConfig instance, which is empty but has a
// non-nil filter and action.
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

// NewSearchConfig returns SearchConfig instance, which is empty but has a
// non-nil filter and action.
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

type ServerConfig struct {
	Name     string `toml:"name"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	LogLines int    `toml:"log_lines,omitempty"`
}
