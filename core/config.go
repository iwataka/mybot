package core

import (
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// TODO: move this module to data package

type Config interface {
	utils.Savable
	utils.Loadable
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
	GetPollingDuration() string
	SetPollingDuration(dur string)
	GetSlackMessages() []MessageConfig
	SetSlackMessages(msgs []MessageConfig)
	AddSlackMessage(msg MessageConfig)
	Validate() error
	ValidateWithAPI(api models.TwitterAPI) error
	Unmarshal(ext string, bytes []byte) error
	Marshal(ext string) ([]byte, error)
}

// FileConfig is a root of the all configurations of this applciation.
type FileConfig struct {
	ConfigProperties `yaml:",inline"`
	// source is a configuration file from which this was loaded. This is
	// needed to save the content to the same file.
	File string `json:"-" toml:"-" bson:"-" yaml:"-"`
}

type ConfigProperties struct {
	// Twitter is a configuration related to Twitter.
	Twitter TwitterConfig `json:"twitter" toml:"twitter" bson:"twitter" yaml:"twitter"`
	// Slack is a configuration related to Slack
	Slack SlackConfig `json:"slack" toml:"slack" bson:"slack" yaml:"slack"`
	// Duration is a duration for some periodic jobs such as fetching
	// users' favorites and searching by the specified condition.
	Duration string `json:"duration" toml:"duration" bson:"duration" yaml:"duration"`
}

func newConfigProperties() ConfigProperties {
	return ConfigProperties{
		Twitter:  NewTwitterConfig(),
		Slack:    NewSlackConfig(),
		Duration: "10m",
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

func (c *ConfigProperties) GetSlackMessages() []MessageConfig {
	return c.Slack.Messages
}

func (c *ConfigProperties) SetSlackMessages(msgs []MessageConfig) {
	c.Slack.Messages = msgs
}

func (c *ConfigProperties) AddSlackMessage(msg MessageConfig) {
	c.Slack.Messages = append(c.Slack.Messages, msg)
}

func (c *ConfigProperties) GetPollingDuration() string {
	return c.Duration
}

func (c *ConfigProperties) SetPollingDuration(dur string) {
	c.Duration = dur
}

// NewFileConfig takes the configuration file path and returns a configuration
// instance.
func NewFileConfig(path string) (*FileConfig, error) {
	c := &FileConfig{newConfigProperties(), path}
	err := c.Load()
	if err != nil {
		return nil, utils.WithStack(err)
	}

	err = c.Validate()
	if err != nil {
		return nil, utils.WithStack(err)
	}

	return c, nil
}

// Validate tries to validate the specified configuration. If invalid values
// are detected, this returns an error.
func (c *ConfigProperties) Validate() error {
	// Validate timeline configurations
	for _, timeline := range c.Twitter.Timelines {
		if err := timeline.Validate(); err != nil {
			return utils.WithStack(err)
		}
	}

	// Validate favorite configurations
	for _, favorite := range c.Twitter.Favorites {
		if err := favorite.Validate(); err != nil {
			return utils.WithStack(err)
		}
	}

	// Validate search configurations
	for _, search := range c.Twitter.Searches {
		if err := search.Validate(); err != nil {
			return utils.WithStack(err)
		}
	}

	for _, msg := range c.Slack.Messages {
		if err := msg.Validate(); err != nil {
			return utils.WithStack(err)
		}
	}

	return nil
}

func (c *ConfigProperties) ValidateWithAPI(api models.TwitterAPI) error {
	for _, name := range c.Twitter.GetScreenNames() {
		_, err := api.GetUsersShow(name, nil)
		if err != nil {
			return utils.WithStack(err)
		}
	}
	return nil
}

// Marshal returns a configuration content as a toml text. If error occurs while
// encoding, this returns an empty string. This return value is not same as the
// source file's content.
func (c *ConfigProperties) Marshal(ext string) ([]byte, error) {
	return utils.Encode(ext, c)
}

// TODO: Make error message clearer
func (c *ConfigProperties) Unmarshal(ext string, bytes []byte) error {
	return utils.Decode(ext, bytes, c)
}

// Save saves the specified configuration to the source file.
func (c *FileConfig) Save() error {
	// Make a directory before all.
	err := os.MkdirAll(filepath.Dir(c.File), 0751)
	if err != nil {
		return utils.WithStack(err)
	}
	if c != nil {
		bs, err := c.Marshal(filepath.Ext(c.File))
		if err != nil {
			return utils.WithStack(err)
		}
		err = ioutil.WriteFile(c.File, bs, 0640)
		if err != nil {
			return utils.WithStack(err)
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
			return utils.WithStack(err)
		}
		err = c.Unmarshal(filepath.Ext(c.File), bytes)
		if err != nil {
			return utils.WithStack(err)
		}
	}
	return nil
}

// Source is a configuration for common data sources such as Twitter's
// timelines, favorites and searches. Sources should have filters and actions.
type Source struct {
	// Name is a label to identify each source
	Name string `json:"name,omitempty" toml:"name,omitempty" bson:"name,omitempty" yaml:"name,omitempty"`
	// Filter filters out incoming data from sources.
	Filter Filter `json:"filter" toml:"filter" bson:"filter" yaml:"filter"`
	// Action defines actions for data passing through filters.
	Action data.Action `json:"action" toml:"action" bson:"action" yaml:"action"`
}

func NewSource() Source {
	return Source{
		Name:   "New",
		Filter: NewFilter(),
		Action: data.NewAction(),
	}
}

func (c *Source) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("No name")
	}
	if c.Action.IsEmpty() {
		return fmt.Errorf("No action")
	}
	return nil
}

func (c *Source) String() string {
	return c.Name
}

// TwitterConfig is a configuration related to Twitter.
type TwitterConfig struct {
	Timelines []TimelineConfig `json:"timelines" toml:"timelines" bson:"timelines" yaml:"timelines"`
	Favorites []FavoriteConfig `json:"favorites" toml:"favorites" bson:"favorites" yaml:"favorites"`
	Searches  []SearchConfig   `json:"searches" toml:"searches" bson:"searches" yaml:"searches"`
	// Debug is a flag for debugging, if it is true, additional information
	// is outputted.
	Debug bool `json:"debug" toml:"debug" bson:"debug" yaml:"debug"`
}

func NewTwitterConfig() TwitterConfig {
	return TwitterConfig{
		Timelines: []TimelineConfig{},
		Favorites: []FavoriteConfig{},
		Searches:  []SearchConfig{},
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
	return utils.UniqStrSlice(result)
}

// TimelineConfig is a configuration for Twitter timelines
type TimelineConfig struct {
	Source                    `yaml:",inline"`
	models.SourceProperties   `yaml:",inline"`
	models.AccountProperties  `yaml:",inline"`
	models.TimelineProperties `yaml:",inline"`
}

// NewTimelineConfig returns TimelineConfig instance, which is empty but has a
// non-nil filter and action.
func NewTimelineConfig() TimelineConfig {
	return TimelineConfig{
		Source:             NewSource(),
		TimelineProperties: models.NewTimelineProperties(),
	}
}

func (c *TimelineConfig) Validate() error {
	err := c.Source.Validate()
	if err != nil {
		return utils.WithStack(err)
	}
	if len(c.ScreenNames) == 0 {
		return fmt.Errorf("%v has no screen names", c)
	}
	return c.Filter.Validate()
}

// FavoriteConfig is a configuration for Twitter favorites
type FavoriteConfig struct {
	Source                    `yaml:",inline"`
	models.SourceProperties   `yaml:",inline"`
	models.AccountProperties  `yaml:",inline"`
	models.FavoriteProperties `yaml:",inline"`
}

// NewFavoriteCnfig returns FavoriteConfig instance, which is empty but has a
// non-nil filter and action.
func NewFavoriteConfig() FavoriteConfig {
	return FavoriteConfig{
		Source:             NewSource(),
		FavoriteProperties: models.NewFavoriteProperties(),
	}
}

func (c *FavoriteConfig) Validate() error {
	err := c.Source.Validate()
	if err != nil {
		return utils.WithStack(err)
	}
	if len(c.ScreenNames) == 0 {
		return fmt.Errorf("%v has no screen names", c)
	}
	return c.Filter.Validate()
}

// SearchConfig is a configuration for Twitter searches
type SearchConfig struct {
	Source                  `yaml:",inline"`
	models.SourceProperties `yaml:",inline"`
	models.SearchProperties `yaml:",inline"`
}

// NewSearchConfig returns SearchConfig instance, which is empty but has a
// non-nil filter and action.
func NewSearchConfig() SearchConfig {
	return SearchConfig{
		Source:           NewSource(),
		SearchProperties: models.NewSearchProperties(),
	}
}

func (c *SearchConfig) Validate() error {
	err := c.Source.Validate()
	if err != nil {
		return utils.WithStack(err)
	}
	if len(c.Queries) == 0 {
		return fmt.Errorf("%v has no queries", c)
	}
	return c.Filter.Validate()
}

type SlackConfig struct {
	Messages []MessageConfig `json:"messages" toml:"messages" bson:"messages" yaml:"messages"`
}

func NewSlackConfig() SlackConfig {
	return SlackConfig{
		Messages: []MessageConfig{},
	}
}

type MessageConfig struct {
	Source   `yaml:",inline"`
	Channels []string `json:"channels" toml:"channels" bson:"channels" yaml:"channels"`
}

func (c MessageConfig) Validate() error {
	err := c.Source.Validate()
	if err != nil {
		return utils.WithStack(err)
	}
	if len(c.Channels) == 0 {
		return fmt.Errorf("At least one channel required")
	}
	return nil
}

func NewMessageConfig() MessageConfig {
	return MessageConfig{
		Source: NewSource(),
	}
}

// TwitterNotification contains some notification settings.
type TwitterNotification struct {
	Place NotificationProperties
}

// NewTwitterNotification returns a new empty Notification.
func NewTwitterNotification() TwitterNotification {
	return TwitterNotification{
		Place: NotificationProperties{},
	}
}

type NotificationProperties struct {
	TwitterAllowSelf bool     `json:"twitter_allow_self" toml:"twitter_allow_self" bson:"twitter_allow_self" yaml:"twitter_allow_self"`
	TwitterUsers     []string `json:"twitter_users,omitempty" toml:"twitter_users,omitempty" bson:"twitter_users,omitempty" yaml:"twitter_users,omitempty"`
	SlackChannels    []string `json:"slack_channels,omitempty" toml:"slack_channels,omitempty" bson:"slack_channels,omitempty" yaml:"slack_channels,omitempty"`
}

// Notify sends a specified messages to certain users according to properties p.
// This returns false if sending the messages to no one.
func (p NotificationProperties) Notify(twitterAPI *TwitterAPI, slackAPI *SlackAPI, msg string) (bool, error) {
	sendsSomeone := false
	allowSelf := p.TwitterAllowSelf
	users := p.TwitterUsers
	for _, user := range users {
		_, err := twitterAPI.api.PostDMToScreenName(msg, user)
		if err != nil {
			return sendsSomeone, utils.WithStack(err)
		}
		sendsSomeone = true
	}
	if allowSelf {
		self, err := twitterAPI.GetSelf()
		if err != nil {
			return sendsSomeone, utils.WithStack(err)
		}
		_, err = twitterAPI.api.PostDMToScreenName(msg, self.ScreenName)
		if err != nil {
			return sendsSomeone, utils.WithStack(err)
		}
		sendsSomeone = true
	}
	chans := p.SlackChannels
	for _, ch := range chans {
		err := slackAPI.PostMessage(ch, msg, nil, false)
		if err != nil {
			return sendsSomeone, utils.WithStack(err)
		}
		sendsSomeone = true
	}
	return sendsSomeone, nil
}

type DBConfig struct {
	ConfigProperties `yaml:",inline"`
	col              *mgo.Collection
	ID               string `json:"id" toml:"id" bson:"id" yaml:"id"`
}

func NewDBConfig(col *mgo.Collection, id string) (*DBConfig, error) {
	c := &DBConfig{newConfigProperties(), col, id}
	err := c.Load()
	return c, utils.WithStack(err)
}

func (c *DBConfig) Load() error {
	query := c.col.Find(bson.M{"id": c.ID})
	count, err := query.Count()
	if err != nil {
		return utils.WithStack(err)
	}
	if count > 0 {
		tmpCol := c.col
		err := query.One(c)
		c.col = tmpCol
		return utils.WithStack(err)
	}
	return nil
}

func (c *DBConfig) Save() error {
	_, err := c.col.Upsert(bson.M{"id": c.ID}, c)
	return utils.WithStack(err)
}
