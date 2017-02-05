package mybot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/BurntSushi/toml"
)

// Config is a root of the all configurations of this applciation.
type Config struct {
	// Twitter is a configuration related to Twitter.
	Twitter *TwitterConfig `json:"twitter" toml:"twitter"`
	// Interaction is a configuration related to interaction with users
	// such as Twitter's direct message exchange.
	Interaction *InteractionConfig `json:"interaction" toml:"interaction"`
	// Log is a configuration related to logging.
	Log *LogConfig `json:"log" toml:"log"`
	// source is a configuration file from which this was loaded. This is
	// needed to save the content to the same file.
	File string `json:"-" toml:"-"`
}

// NewConfig takes the configuration file path and returns a configuration
// instance.
func NewConfig(path string) (*Config, error) {
	c := &Config{
		Twitter: &TwitterConfig{
			Timelines: []TimelineConfig{},
			Searches:  []SearchConfig{},
			Duration:  "1h",
			Notification: &Notification{
				Place: &PlaceNotification{},
			},
		},
		Log: &LogConfig{
			Linenum: 10,
		},
		Interaction: &InteractionConfig{},
	}

	c.File = path
	err := c.Load()
	if err != nil {
		return nil, err
	}

	// Assign empty values to config instance to prevent nil pointer
	// reference error.
	for _, t := range c.Twitter.Timelines {
		t.Init()
	}
	for _, f := range c.Twitter.Favorites {
		f.Init()
	}
	for _, s := range c.Twitter.Searches {
		s.Init()
	}

	err = c.Validate()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Validate tries to validate the specified configuration. If invalid values
// are detected, this returns an error.
func (c *Config) Validate() error {
	// Validate timeline configurations
	for _, timeline := range c.Twitter.Timelines {
		if timeline.Action == nil {
			msg := fmt.Sprintf("%v has no action", timeline)
			return errors.New(msg)
		}
		if len(timeline.ScreenNames) == 0 {
			msg := fmt.Sprintf("%v has no name", timeline)
			return errors.New(msg)
		}
		filter := timeline.Filter
		if !filter.Vision.isEmpty() &&
			(filter.RetweetedThreshold != nil || filter.FavoriteThreshold != nil) {
			bytes, _ := json.Marshal(timeline)
			msg := fmt.Sprintf("%s\n%s",
				"Don't use both of Vision API and retweeted/favorite threshold",
				string(bytes),
			)
			return errors.New(msg)
		}
	}

	// Validate favorite configurations
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
			(filter.RetweetedThreshold != nil || filter.FavoriteThreshold != nil) {
			bytes, _ := json.Marshal(favorite)
			msg := fmt.Sprintf("%s\n%s",
				"Don't use both of Vision API and retweeted/favorite threshold",
				string(bytes),
			)
			return errors.New(msg)
		}
	}

	// Validate search configurations
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
			(filter.RetweetedThreshold != nil || filter.FavoriteThreshold != nil) {
			bytes, _ := json.Marshal(search)
			msg := fmt.Sprintf("%s\n%s",
				"Don't use both of Vision API and retweeted/favorite threshold",
				string(bytes),
			)
			return errors.New(msg)
		}
	}

	// Validate API configurations
	for _, api := range c.Twitter.APIs {
		if len(api.SourceURL) == 0 {
			return errors.New("API source URL shouldn't be empty")
		}
		if len(api.MessageTemplate) == 0 {
			return errors.New("API message template shouldn't be empty")
		}
	}

	return nil
}

func (c *Config) ValidateWithAPI(api *TwitterAPI) error {
	for _, name := range c.Twitter.GetScreenNames() {
		_, err := api.api.GetUsersShow(name, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// ToText returns a configuration content as a toml text. If error occurs while
// encoding, this returns an empty string. This return value is not same as the
// source file's content.
func (c *Config) ToText(indent string) ([]byte, error) {
	ext := filepath.Ext(c.File)
	buf := new(bytes.Buffer)
	switch ext {
	case ".json":
		b := new(bytes.Buffer)
		enc := json.NewEncoder(b)
		err := enc.Encode(c)
		if err != nil {
			return []byte{}, err
		}
		// go1.6 or lower doesn't support json.Encoder#SetIndent.
		err = json.Indent(buf, b.Bytes(), "", indent)
		if err != nil {
			return []byte{}, err
		}
	case ".toml":
		enc := toml.NewEncoder(buf)
		enc.Indent = indent
		err := enc.Encode(c)
		if err != nil {
			return []byte{}, err
		}
	}
	return buf.Bytes(), nil
}

func (c *Config) FromText(bytes []byte) error {
	ext := filepath.Ext(c.File)
	switch ext {
	case ".json":
		err := json.Unmarshal(bytes, c)
		if err != nil {
			return err
		}
	case ".toml":
		md, err := toml.Decode(string(bytes), c)
		if err != nil {
			return err
		}
		if len(md.Undecoded()) != 0 {
			return fmt.Errorf("%v undecoded in %s", md.Undecoded(), c.File)
		}
	}
	return nil
}

// Save saves the specified configuration to the source file.
func (c *Config) Save() error {
	// Make a directory before all.
	err := os.MkdirAll(filepath.Dir(c.File), 0751)
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
		err = ioutil.WriteFile(c.File, writer.Bytes(), 0640)
		if err != nil {
			return err
		}
	}
	return nil
}

// Load loads the configuration from the source file. If the specified source
// file doesn't exist, this method does nothing and returns nil.
func (c *Config) Load() error {
	if info, err := os.Stat(c.File); err == nil && !info.IsDir() {
		bytes, err := ioutil.ReadFile(c.File)
		if err != nil {
			return err
		}
		err = c.FromText(bytes)
		if err != nil {
			return err
		}
	}
	return nil
}

// SourceConfig is a configuration for common data sources such as Twitter's
// timelines, favorites and searches. Sources should have filters and actions.
type SourceConfig struct {
	// Filter filters out incoming data from sources.
	Filter *TweetFilter `json:"filter" toml:"filter"`
	// Action defines actions for data passing through filters.
	Action *TweetAction `json:"action" toml:"action"`
}

func (c *SourceConfig) Init() {
	if c.Filter.Vision == nil {
		c.Filter.Vision = new(VisionCondition)
	}
	if c.Filter.Vision.Face == nil {
		c.Filter.Vision.Face = new(VisionFaceCondition)
	}
	if c.Filter.Language == nil {
		c.Filter.Language = new(LanguageCondition)
	}

	if c.Action.Twitter == nil {
		c.Action.Twitter = NewTwitterAction()
	}
	if c.Action.Slack == nil {
		c.Action.Slack = NewSlackAction()
	}
}

type TweetAction struct {
	Twitter *TwitterAction `json:"twitter" toml:"twitter"`
	Slack   *SlackAction   `json:"slack" toml:"slack"`
}

func NewTweetAction() *TweetAction {
	return &TweetAction{
		Twitter: NewTwitterAction(),
		Slack:   NewSlackAction(),
	}
}

func (a *TweetAction) Add(action *TweetAction) {
	if a.Twitter == nil {
		a.Twitter = action.Twitter
	} else {
		a.Twitter.Add(action.Twitter)
	}

	if a.Slack == nil {
		a.Twitter = action.Twitter
	} else {
		a.Slack.Add(action.Slack)
	}
}

func (a *TweetAction) Sub(action *TweetAction) {
	if a.Twitter != nil {
		a.Twitter.Sub(action.Twitter)
	}

	if a.Slack != nil {
		a.Slack.Sub(action.Slack)
	}
}

// TwitterConfig is a configuration related to Twitter.
type TwitterConfig struct {
	Timelines []TimelineConfig `json:"timelines" toml:"timelines"`
	Favorites []FavoriteConfig `json:"favorites" toml:"favorites"`
	Searches  []SearchConfig   `json:"searches" toml:"searches"`
	APIs      []APIConfig      `json:"api" toml:"api"`
	// Notification is a configuration related to notification for users.
	// Currently only place notification is supported, which means that
	// when a tweet with place information is detected, it is notified to
	// the specified users.
	Notification *Notification `json:"notification" toml:"notification"`
	// Duration is a duration for some periodic jobs such as fetching
	// users' favorites and searching by the specified condition.
	Duration string `json:"duration" toml:"duration"`
	// Debug is a flag for debugging, if it is true, additional information
	// is outputted.
	Debug bool `json:"debug" toml:"debug"`
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
	ScreenNames    []string `json:"screen_names" toml:"screen_names"`
	ExcludeReplies *bool    `json:"exclude_replies" toml:"exclude_replies"`
	IncludeRts     *bool    `json:"include_rts" toml:"include_rts"`
	Count          *int     `json:"count,omitempty" toml:"count,omitempty"`
}

// NewTimelineConfig returns TimelineConfig instance, which is empty but has a
// non-nil filter and action.
func NewTimelineConfig() *TimelineConfig {
	return &TimelineConfig{
		SourceConfig: &SourceConfig{
			Filter: NewTweetFilter(),
			Action: NewTweetAction(),
		},
	}
}

// FavoriteConfig is a configuration for Twitter favorites
type FavoriteConfig struct {
	*SourceConfig
	ScreenNames []string `json:"screen_names" toml:"screen_names"`
	Count       *int     `json:"count,omitempty" toml:"count,omitempty"`
}

// NewFavoriteCnfig returns FavoriteConfig instance, which is empty but has a
// non-nil filter and action.
func NewFavoriteConfig() *FavoriteConfig {
	return &FavoriteConfig{
		SourceConfig: &SourceConfig{
			Filter: NewTweetFilter(),
			Action: NewTweetAction(),
		},
	}
}

// SearchConfig is a configuration for Twitter searches
type SearchConfig struct {
	*SourceConfig
	Queries    []string `json:"queries" toml:"queries"`
	ResultType string   `json:"result_type,omitempty" toml:"result_type,omitempty"`
	Count      *int     `json:"count,omitempty" toml:"count,omitempty"`
}

// NewSearchConfig returns SearchConfig instance, which is empty but has a
// non-nil filter and action.
func NewSearchConfig() *SearchConfig {
	return &SearchConfig{
		SourceConfig: &SourceConfig{
			Filter: NewTweetFilter(),
			Action: NewTweetAction(),
		},
	}
}

type APIConfig struct {
	SourceURL       string `json:"source_url,omitempty" toml:"source_url,omitempty"`
	MessageTemplate string `json:"message_template,omitempty" toml:"message_template,omitempty"`
}

func NewAPIConfig() *APIConfig {
	return &APIConfig{}
}

func (c *APIConfig) Message() (string, error) {
	// Send GET request to retrieve JSON data
	client := &http.Client{}
	res, err := client.Get(c.SourceURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	// Unmarshal json data
	var data interface{}
	err = json.Unmarshal(bs, &data)
	if err != nil {
		return "", err
	}

	// Returns a message generated by json data and a message template
	tmpl, err := template.New("message").Parse(c.MessageTemplate)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// InteractionConfig is a configuration for interaction through Twitter direct
// message
type InteractionConfig struct {
	AllowSelf bool     `json:"allow_self" toml:"allow_self"`
	Users     []string `json:"users,omitempty" toml:"users,omitempty"`
}

// LogConfig is a configuration for logging
type LogConfig struct {
	AllowSelf bool     `json:"allow_self" toml:"allow_self"`
	Users     []string `json:"users,omitempty" toml:"users,omitempty"`
	Linenum   int      `json:"linenum,omitempty" toml:"linenum,omitempty"`
}
