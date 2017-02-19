package mybot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/iwataka/mybot/models"
)

type Config interface {
	GetTwitterScreenNames() ([]string, error)
	GetTwitterTimelines() ([]TimelineConfig, error)
	SetTwitterTimelines(timelines []TimelineConfig) error
	AddTwitterTimeline(timeline *TimelineConfig) error
	GetTwitterFavorites() ([]FavoriteConfig, error)
	SetTwitterFavorites(favorites []FavoriteConfig) error
	AddTwitterFavorite(favorite *FavoriteConfig) error
	GetTwitterSearches() ([]SearchConfig, error)
	SetTwitterSearches(searches []SearchConfig) error
	AddTwitterSearch(search *SearchConfig) error
	GetTwitterAPIs() ([]APIConfig, error)
	SetTwitterAPIs(apis []APIConfig) error
	GetTwitterNotification() (*Notification, error)
	SetTwitterNotification(notification *Notification) error
	GetTwitterInteraction() (*InteractionConfig, error)
	SetTwitterInteraction(interaction *InteractionConfig) error
	GetTwitterDuration() (string, error)
	SetTwitterDuration(dur string) error
	GetSlackMessages() ([]MessageConfig, error)
	SetSlackMessages(msgs []MessageConfig) error
	Load() error
	Save() error
}

// FileConfig is a root of the all configurations of this applciation.
type FileConfig struct {
	// Twitter is a configuration related to Twitter.
	Twitter *TwitterConfig `json:"twitter" toml:"twitter"`
	// Slack is a configuration related to Slack
	Slack *SlackConfig `json:"slack" toml:"slack"`
	// source is a configuration file from which this was loaded. This is
	// needed to save the content to the same file.
	File string `json:"-" toml:"-"`
}

func (c *FileConfig) GetTwitterScreenNames() ([]string, error) {
	return c.Twitter.GetScreenNames(), nil
}

func (c *FileConfig) GetTwitterTimelines() ([]TimelineConfig, error) {
	return c.Twitter.Timelines, nil
}

func (c *FileConfig) SetTwitterTimelines(timelines []TimelineConfig) error {
	c.Twitter.Timelines = timelines
	return nil
}

func (c *FileConfig) AddTwitterTimeline(timeline *TimelineConfig) error {
	if timeline == nil {
		return nil
	}
	c.Twitter.Timelines = append(c.Twitter.Timelines, *timeline)
	return nil
}

func (c *FileConfig) GetTwitterFavorites() ([]FavoriteConfig, error) {
	return c.Twitter.Favorites, nil
}

func (c *FileConfig) SetTwitterFavorites(favorites []FavoriteConfig) error {
	c.Twitter.Favorites = favorites
	return nil
}

func (c *FileConfig) AddTwitterFavorite(favorite *FavoriteConfig) error {
	if favorite == nil {
		return nil
	}
	c.Twitter.Favorites = append(c.Twitter.Favorites, *favorite)
	return nil
}

func (c *FileConfig) GetTwitterSearches() ([]SearchConfig, error) {
	return c.Twitter.Searches, nil
}

func (c *FileConfig) SetTwitterSearches(searches []SearchConfig) error {
	c.Twitter.Searches = searches
	return nil
}

func (c *FileConfig) AddTwitterSearch(search *SearchConfig) error {
	if search == nil {
		return nil
	}
	c.Twitter.Searches = append(c.Twitter.Searches, *search)
	return nil
}

func (c *FileConfig) GetTwitterAPIs() ([]APIConfig, error) {
	return c.Twitter.APIs, nil
}

func (c *FileConfig) SetTwitterAPIs(apis []APIConfig) error {
	c.Twitter.APIs = apis
	return nil
}

func (c *FileConfig) GetTwitterNotification() (*Notification, error) {
	return c.Twitter.Notification, nil
}

func (c *FileConfig) SetTwitterNotification(notification *Notification) error {
	c.Twitter.Notification = notification
	return nil
}

func (c *FileConfig) GetTwitterInteraction() (*InteractionConfig, error) {
	return c.Twitter.Interaction, nil
}

func (c *FileConfig) SetTwitterInteraction(interaction *InteractionConfig) error {
	c.Twitter.Interaction = interaction
	return nil
}

func (c *FileConfig) GetSlackMessages() ([]MessageConfig, error) {
	return c.Slack.Messages, nil
}

func (c *FileConfig) SetSlackMessages(msgs []MessageConfig) error {
	c.Slack.Messages = msgs
	return nil
}

func (c *FileConfig) GetTwitterDuration() (string, error) {
	return c.Twitter.Duration, nil
}

func (c *FileConfig) SetTwitterDuration(dur string) error {
	c.Twitter.Duration = dur
	return nil
}

// NewFileConfig takes the configuration file path and returns a configuration
// instance.
func NewFileConfig(path string) (*FileConfig, error) {
	c := &FileConfig{
		Twitter: NewTwitterConfig(),
		Slack:   NewSlackConfig(),
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
func (c *FileConfig) Validate() error {
	// Validate timeline configurations
	for _, timeline := range c.Twitter.Timelines {
		err := timeline.Validate()
		if err != nil {
			return err
		}
	}

	// Validate favorite configurations
	for _, favorite := range c.Twitter.Favorites {
		err := favorite.Validate()
		if err != nil {
			return err
		}
	}

	// Validate search configurations
	for _, search := range c.Twitter.Searches {
		err := search.Validate()
		if err != nil {
			return err
		}
	}

	// Validate API configurations
	for _, api := range c.Twitter.APIs {
		err := api.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *FileConfig) ValidateWithAPI(api *TwitterAPI) error {
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
func (c *FileConfig) ToText(indent string) ([]byte, error) {
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

func (c *FileConfig) FromText(bytes []byte) error {
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
func (c *FileConfig) Save() error {
	// Make a directory before all.
	err := os.MkdirAll(filepath.Dir(c.File), 0751)
	if err != nil {
		return err
	}
	if c != nil {
		bs, err := c.ToText("")
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(c.File, bs, 0640)
		if err != nil {
			return err
		}
	}
	return nil
}

// Load loads the configuration from the source file. If the specified source
// file doesn't exist, this method does nothing and returns nil.
func (c *FileConfig) Load() error {
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
	Filter *Filter `json:"filter" toml:"filter"`
	// Action defines actions for data passing through filters.
	Action *Action `json:"action" toml:"action"`
}

func NewSourceConfig() SourceConfig {
	return SourceConfig{
		Filter: NewFilter(),
		Action: NewAction(),
	}
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

func (c *SourceConfig) Validate() error {
	if c.Action == nil || c.Action.IsEmpty() {
		return fmt.Errorf("%v has no action", c)
	}
	return nil
}

type Action struct {
	Twitter *TwitterAction `json:"twitter" toml:"twitter"`
	Slack   *SlackAction   `json:"slack" toml:"slack"`
}

func NewAction() *Action {
	return &Action{
		Twitter: NewTwitterAction(),
		Slack:   NewSlackAction(),
	}
}

func (a *Action) Add(action *Action) *Action {
	if action == nil {
		return a
	}

	result := *a

	if a.Twitter == nil {
		result.Twitter = action.Twitter
	} else {
		result.Twitter = a.Twitter.Add(action.Twitter)
	}

	if a.Slack == nil {
		result.Slack = action.Slack
	} else {
		result.Slack = a.Slack.Add(action.Slack)
	}

	return &result
}

func (a *Action) Sub(action *Action) *Action {
	if action == nil {
		return a
	}

	result := *a

	if a.Twitter != nil {
		result.Twitter = a.Twitter.Sub(action.Twitter)
	}
	if a.Slack != nil {
		result.Slack = a.Slack.Sub(action.Slack)
	}

	return &result
}

func (a *Action) IsEmpty() bool {
	return a.Twitter.IsEmpty() && a.Slack.IsEmpty()
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
	// Interaction is a configuration related to interaction with users
	// such as Twitter's direct message exchange.
	Interaction *InteractionConfig `json:"interaction" toml:"interaction"`
	// Duration is a duration for some periodic jobs such as fetching
	// users' favorites and searching by the specified condition.
	Duration string `json:"duration" toml:"duration"`
	// Debug is a flag for debugging, if it is true, additional information
	// is outputted.
	Debug bool `json:"debug" toml:"debug"`
}

func NewTwitterConfig() *TwitterConfig {
	return &TwitterConfig{
		Timelines:    []TimelineConfig{},
		Searches:     []SearchConfig{},
		Interaction:  &InteractionConfig{},
		Duration:     "1h",
		Notification: NewNotification(),
	}
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
	SourceConfig
	models.SourceProperties
	models.AccountProperties
	models.TimelineProperties
}

// NewTimelineConfig returns TimelineConfig instance, which is empty but has a
// non-nil filter and action.
func NewTimelineConfig() *TimelineConfig {
	return &TimelineConfig{
		SourceConfig: NewSourceConfig(),
	}
}

func (c *TimelineConfig) Validate() error {
	err := c.SourceConfig.Validate()
	if err != nil {
		return err
	}
	if len(c.ScreenNames) == 0 {
		return fmt.Errorf("%v has no screen names", c)
	}
	return c.Filter.Validate()
}

// FavoriteConfig is a configuration for Twitter favorites
type FavoriteConfig struct {
	SourceConfig
	models.SourceProperties
	models.AccountProperties
	models.FavoriteProperties
}

// NewFavoriteCnfig returns FavoriteConfig instance, which is empty but has a
// non-nil filter and action.
func NewFavoriteConfig() *FavoriteConfig {
	return &FavoriteConfig{
		SourceConfig: NewSourceConfig(),
	}
}

func (c *FavoriteConfig) Validate() error {
	err := c.SourceConfig.Validate()
	if err != nil {
		return err
	}
	if len(c.ScreenNames) == 0 {
		return fmt.Errorf("%v has no screen names", c)
	}
	return c.Filter.Validate()
}

// SearchConfig is a configuration for Twitter searches
type SearchConfig struct {
	SourceConfig
	models.SourceProperties
	models.SearchProperties
}

// NewSearchConfig returns SearchConfig instance, which is empty but has a
// non-nil filter and action.
func NewSearchConfig() *SearchConfig {
	return &SearchConfig{
		SourceConfig: NewSourceConfig(),
	}
}

func (c *SearchConfig) Validate() error {
	err := c.SourceConfig.Validate()
	if err != nil {
		return err
	}
	if len(c.Queries) == 0 {
		return fmt.Errorf("%v has no queries", c)
	}
	return c.Filter.Validate()
}

type APIConfig struct {
	SourceURL       string `json:"source_url,omitempty" toml:"source_url,omitempty"`
	MessageTemplate string `json:"message_template,omitempty" toml:"message_template,omitempty"`
}

func NewAPIConfig() *APIConfig {
	return &APIConfig{}
}

func (c *APIConfig) Validate() error {
	if len(c.SourceURL) == 0 {
		return fmt.Errorf("%v has no source URL", c)
	}
	if len(c.MessageTemplate) == 0 {
		return fmt.Errorf("%v has no message template", c)
	}
	return nil
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

type SlackConfig struct {
	Messages []MessageConfig `json:"messages" toml:"messages"`
}

func NewSlackConfig() *SlackConfig {
	return &SlackConfig{
		Messages: []MessageConfig{},
	}
}

type MessageConfig struct {
	SourceConfig
	Channels []string `json:"channels" toml:"channels"`
}

func NewMessageConfig() *MessageConfig {
	return &MessageConfig{
		SourceConfig: NewSourceConfig(),
	}
}

// InteractionConfig is a configuration for interaction through Twitter direct
// message
type InteractionConfig struct {
	AllowSelf bool     `json:"allow_self" toml:"allow_self"`
	Users     []string `json:"users,omitempty" toml:"users,omitempty"`
}
