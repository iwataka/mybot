package mybot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/iwataka/mybot/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Config interface {
	Savable
	Loadable
	GetProperties() *ConfigProperties
	GetTwitterScreenNames() []string
	GetTwitterTimelines() []TimelineConfig
	SetTwitterTimelines(timelines []TimelineConfig)
	AddTwitterTimeline(timeline TimelineConfig)
	GetTwitterFavorites() []FavoriteConfig
	SetTwitterFavorites(favorites []FavoriteConfig)
	AddTwitterFavorite(favorite FavoriteConfig)
	GetTwitterSearches() []SearchConfig
	SetTwitterSearches(searches []SearchConfig)
	AddTwitterSearch(search SearchConfig)
	GetTwitterNotification() Notification
	SetTwitterNotification(notification Notification)
	GetTwitterInteraction() InteractionConfig
	SetTwitterInteraction(interaction InteractionConfig)
	GetTwitterDuration() string
	SetTwitterDuration(dur string)
	GetSlackMessages() []MessageConfig
	SetSlackMessages(msgs []MessageConfig)
	AddSlackMessage(msg MessageConfig)
	Validate() error
	ValidateWithAPI(api *TwitterAPI) error
	Unmarshal(bytes []byte) error
	Marshal(indent, ext string) ([]byte, error)
}

// FileConfig is a root of the all configurations of this applciation.
type FileConfig struct {
	*ConfigProperties
	// source is a configuration file from which this was loaded. This is
	// needed to save the content to the same file.
	File string `json:"-" toml:"-" bson:"-"`
}

type ConfigProperties struct {
	// Twitter is a configuration related to Twitter.
	Twitter TwitterConfig `json:"twitter" toml:"twitter" bson:"twitter"`
	// Slack is a configuration related to Slack
	Slack SlackConfig `json:"slack" toml:"slack" bson:"slack"`
}

func newConfigProperties() *ConfigProperties {
	return &ConfigProperties{
		Twitter: NewTwitterConfig(),
		Slack:   NewSlackConfig(),
	}
}

func (c *ConfigProperties) GetProperties() *ConfigProperties {
	return c
}

func (c *ConfigProperties) GetTwitterScreenNames() []string {
	return c.Twitter.GetScreenNames()
}

func (c *ConfigProperties) GetTwitterTimelines() []TimelineConfig {
	return c.Twitter.Timelines
}

func (c *ConfigProperties) SetTwitterTimelines(timelines []TimelineConfig) {
	c.Twitter.Timelines = timelines
}

func (c *ConfigProperties) AddTwitterTimeline(timeline TimelineConfig) {
	c.Twitter.Timelines = append(c.Twitter.Timelines, timeline)
}

func (c *ConfigProperties) GetTwitterFavorites() []FavoriteConfig {
	return c.Twitter.Favorites
}

func (c *ConfigProperties) SetTwitterFavorites(favorites []FavoriteConfig) {
	c.Twitter.Favorites = favorites
}

func (c *ConfigProperties) AddTwitterFavorite(favorite FavoriteConfig) {
	c.Twitter.Favorites = append(c.Twitter.Favorites, favorite)
}

func (c *ConfigProperties) GetTwitterSearches() []SearchConfig {
	return c.Twitter.Searches
}

func (c *ConfigProperties) SetTwitterSearches(searches []SearchConfig) {
	c.Twitter.Searches = searches
}

func (c *ConfigProperties) AddTwitterSearch(search SearchConfig) {
	c.Twitter.Searches = append(c.Twitter.Searches, search)
}

func (c *ConfigProperties) GetTwitterNotification() Notification {
	return c.Twitter.Notification
}

func (c *ConfigProperties) SetTwitterNotification(notification Notification) {
	c.Twitter.Notification = notification
}

func (c *ConfigProperties) GetTwitterInteraction() InteractionConfig {
	return c.Twitter.Interaction
}

func (c *ConfigProperties) SetTwitterInteraction(interaction InteractionConfig) {
	c.Twitter.Interaction = interaction
}

func (c *ConfigProperties) GetSlackMessages() []MessageConfig {
	return c.Slack.Messages
}

func (c *ConfigProperties) SetSlackMessages(msgs []MessageConfig) {
	c.Slack.Messages = msgs
}

func (c *ConfigProperties) AddSlackMessage(msg MessageConfig) {
	c.Slack.Messages = append(c.Slack.Messages, msg)
}

func (c *ConfigProperties) GetTwitterDuration() string {
	return c.Twitter.Duration
}

func (c *ConfigProperties) SetTwitterDuration(dur string) {
	c.Twitter.Duration = dur
}

// NewFileConfig takes the configuration file path and returns a configuration
// instance.
func NewFileConfig(path string) (*FileConfig, error) {
	c := &FileConfig{newConfigProperties(), path}
	err := c.Load()
	if err != nil {
		return nil, WithStack(err)
	}

	err = c.Validate()
	if err != nil {
		return nil, WithStack(err)
	}

	return c, nil
}

// Validate tries to validate the specified configuration. If invalid values
// are detected, this returns an error.
func (c *ConfigProperties) Validate() error {
	// Validate timeline configurations
	for _, timeline := range c.Twitter.Timelines {
		if err := timeline.Validate(); err != nil {
			return WithStack(err)
		}
	}

	// Validate favorite configurations
	for _, favorite := range c.Twitter.Favorites {
		if err := favorite.Validate(); err != nil {
			return WithStack(err)
		}
	}

	// Validate search configurations
	for _, search := range c.Twitter.Searches {
		if err := search.Validate(); err != nil {
			return WithStack(err)
		}
	}

	for _, msg := range c.Slack.Messages {
		if err := msg.Validate(); err != nil {
			return WithStack(err)
		}
	}

	return nil
}

func (c *ConfigProperties) ValidateWithAPI(api *TwitterAPI) error {
	for _, name := range c.Twitter.GetScreenNames() {
		_, err := api.API.GetUsersShow(name, nil)
		if err != nil {
			return WithStack(err)
		}
	}
	return nil
}

// Marshal returns a configuration content as a toml text. If error occurs while
// encoding, this returns an empty string. This return value is not same as the
// source file's content.
func (c *ConfigProperties) Marshal(indent, ext string) ([]byte, error) {
	buf := new(bytes.Buffer)
	switch ext {
	case ".json":
		b := new(bytes.Buffer)
		enc := json.NewEncoder(b)
		err := enc.Encode(c)
		if err != nil {
			return []byte{}, WithStack(err)
		}
		// go1.6 or lower doesn't support json.Encoder#SetIndent.
		err = json.Indent(buf, b.Bytes(), "", indent)
		if err != nil {
			return []byte{}, WithStack(err)
		}
	case ".toml":
		enc := toml.NewEncoder(buf)
		enc.Indent = indent
		err := enc.Encode(c)
		if err != nil {
			return []byte{}, WithStack(err)
		}
	}
	return buf.Bytes(), nil
}

func (c *ConfigProperties) Unmarshal(bytes []byte) error {
	err := json.Unmarshal(bytes, c)
	if err == nil {
		return nil
	}

	_, err = toml.Decode(string(bytes), c)
	if err == nil {
		return nil
	}

	return Errorf("Configuration must be written in either json or toml format")
}

// Save saves the specified configuration to the source file.
func (c *FileConfig) Save() error {
	// Make a directory before all.
	err := os.MkdirAll(filepath.Dir(c.File), 0751)
	if err != nil {
		return WithStack(err)
	}
	if c != nil {
		bs, err := c.Marshal("", filepath.Ext(c.File))
		if err != nil {
			return WithStack(err)
		}
		err = ioutil.WriteFile(c.File, bs, 0640)
		if err != nil {
			return WithStack(err)
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
			return WithStack(err)
		}
		err = c.Unmarshal(bytes)
		if err != nil {
			return WithStack(err)
		}
	}
	return nil
}

// Source is a configuration for common data sources such as Twitter's
// timelines, favorites and searches. Sources should have filters and actions.
type Source struct {
	// Filter filters out incoming data from sources.
	Filter Filter `json:"filter" toml:"filter" bson:"filter"`
	// Action defines actions for data passing through filters.
	Action Action `json:"action" toml:"action" bson:"action"`
}

func NewSource() Source {
	return Source{
		Filter: NewFilter(),
		Action: NewAction(),
	}
}

func (c *Source) Validate() error {
	if c.Action.IsEmpty() {
		return Errorf("%v has no action", c)
	}
	return nil
}

type Action struct {
	Twitter TwitterAction `json:"twitter" toml:"twitter" bson:"twitter"`
	Slack   SlackAction   `json:"slack" toml:"slack" bson:"slack"`
}

func NewAction() Action {
	return Action{
		Twitter: NewTwitterAction(),
		Slack:   NewSlackAction(),
	}
}

func (a Action) Add(action Action) Action {
	result := a
	result.Twitter = a.Twitter.Add(action.Twitter)
	result.Slack = a.Slack.Add(action.Slack)
	return result
}

func (a Action) Sub(action Action) Action {
	result := a
	result.Twitter = a.Twitter.Sub(action.Twitter)
	result.Slack = a.Slack.Sub(action.Slack)
	return result
}

func (a Action) IsEmpty() bool {
	return a.Twitter.IsEmpty() && a.Slack.IsEmpty()
}

// TwitterConfig is a configuration related to Twitter.
type TwitterConfig struct {
	Timelines []TimelineConfig `json:"timelines" toml:"timelines" bson:"timelines"`
	Favorites []FavoriteConfig `json:"favorites" toml:"favorites" bson:"favorites"`
	Searches  []SearchConfig   `json:"searches" toml:"searches" bson:"searches"`
	// Notification is a configuration related to notification for users.
	// Currently only place notification is supported, which means that
	// when a tweet with place information is detected, it is notified to
	// the specified users.
	Notification Notification `json:"notification" toml:"notification" bson:"notification"`
	// Interaction is a configuration related to interaction with users
	// such as Twitter's direct message exchange.
	Interaction InteractionConfig `json:"interaction" toml:"interaction" bson:"interaction"`
	// Duration is a duration for some periodic jobs such as fetching
	// users' favorites and searching by the specified condition.
	Duration string `json:"duration" toml:"duration" bson:"duration"`
	// Debug is a flag for debugging, if it is true, additional information
	// is outputted.
	Debug bool `json:"debug" toml:"debug" bson:"debug"`
}

func NewTwitterConfig() TwitterConfig {
	return TwitterConfig{
		Timelines:    []TimelineConfig{},
		Searches:     []SearchConfig{},
		Interaction:  InteractionConfig{},
		Duration:     "1h",
		Notification: NewNotification(),
	}
}

// GetScreenNames returns all screen names in the TwitterConfig instance. This
// is useful to do something for all related users.
func (tc TwitterConfig) GetScreenNames() []string {
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
	Source
	models.SourceProperties
	models.AccountProperties
	models.TimelineProperties
}

// NewTimelineConfig returns TimelineConfig instance, which is empty but has a
// non-nil filter and action.
func NewTimelineConfig() TimelineConfig {
	return TimelineConfig{
		Source: NewSource(),
	}
}

func (c *TimelineConfig) Validate() error {
	err := c.Source.Validate()
	if err != nil {
		return WithStack(err)
	}
	if len(c.ScreenNames) == 0 {
		return Errorf("%v has no screen names", c)
	}
	return c.Filter.Validate()
}

// FavoriteConfig is a configuration for Twitter favorites
type FavoriteConfig struct {
	Source
	models.SourceProperties
	models.AccountProperties
	models.FavoriteProperties
}

// NewFavoriteCnfig returns FavoriteConfig instance, which is empty but has a
// non-nil filter and action.
func NewFavoriteConfig() FavoriteConfig {
	return FavoriteConfig{
		Source: NewSource(),
	}
}

func (c *FavoriteConfig) Validate() error {
	err := c.Source.Validate()
	if err != nil {
		return WithStack(err)
	}
	if len(c.ScreenNames) == 0 {
		return Errorf("%v has no screen names", c)
	}
	return c.Filter.Validate()
}

// SearchConfig is a configuration for Twitter searches
type SearchConfig struct {
	Source
	models.SourceProperties
	models.SearchProperties
}

// NewSearchConfig returns SearchConfig instance, which is empty but has a
// non-nil filter and action.
func NewSearchConfig() SearchConfig {
	return SearchConfig{
		Source: NewSource(),
	}
}

func (c *SearchConfig) Validate() error {
	err := c.Source.Validate()
	if err != nil {
		return WithStack(err)
	}
	if len(c.Queries) == 0 {
		return Errorf("%v has no queries", c)
	}
	return c.Filter.Validate()
}

// InteractionConfig is a configuration for interaction through Twitter direct
// message
type InteractionConfig struct {
	AllowSelf bool     `json:"allow_self" toml:"allow_self" bson:"allow_self"`
	Users     []string `json:"users,omitempty" toml:"users,omitempty" bson:"users,omitempty"`
}
type SlackConfig struct {
	Messages []MessageConfig `json:"messages" toml:"messages" bson:"messages"`
}

func NewSlackConfig() SlackConfig {
	return SlackConfig{
		Messages: []MessageConfig{},
	}
}

type MessageConfig struct {
	Source
	Channels []string `json:"channels" toml:"channels" bson:"channels"`
}

func (c MessageConfig) Validate() error {
	err := c.Source.Validate()
	if err != nil {
		return WithStack(err)
	}
	if len(c.Channels) == 0 {
		return Errorf("At least one channel required")
	}
	return nil
}

func NewMessageConfig() MessageConfig {
	return MessageConfig{
		Source: NewSource(),
	}
}

type DBConfig struct {
	*ConfigProperties
	col *mgo.Collection `json:"-" toml:"-" bson:"-"`
	ID  string          `json:"id" toml:"id" bson:"id"`
}

func NewDBConfig(col *mgo.Collection, id string) (*DBConfig, error) {
	c := &DBConfig{newConfigProperties(), col, id}
	err := c.Load()
	return c, WithStack(err)
}

func (c *DBConfig) Load() error {
	query := c.col.Find(bson.M{"id": c.ID})
	count, err := query.Count()
	if err != nil {
		return WithStack(err)
	}
	if count > 0 {
		tmpCol := c.col
		err := query.One(c)
		c.col = tmpCol
		return WithStack(err)
	}
	return nil
}

func (c *DBConfig) Save() error {
	_, err := c.col.Upsert(bson.M{"id": c.ID}, c)
	return WithStack(err)
}
