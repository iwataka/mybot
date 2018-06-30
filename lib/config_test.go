package mybot_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/iwataka/mybot/data"
	. "github.com/iwataka/mybot/lib"
	"github.com/stretchr/testify/require"
)

const (
	defaultTestConfigFilePath = "testdata/config.template.toml"
)

func TestNewConfig(t *testing.T) {
	c := NewTestFileConfig(defaultTestConfigFilePath, t)

	a := c.GetTwitterTimelines()[0]
	require.Equal(t, "golang", a.ScreenNames[0])
	f := a.Filter
	require.Equal(t, "is released!", f.Patterns[0])
	require.False(t, *f.Retweeted)
	require.Equal(t, "en", f.Lang)
	require.Equal(t, "cartoon|clip art|artwork", f.Vision.Label[0])
	require.True(t, a.Action.Twitter.Retweet)
	require.Equal(t, "foo", a.Action.Slack.Channels[0])

	s := c.GetTwitterSearches()[0]
	require.Equal(t, "foo", s.Queries[0])
	require.Equal(t, "bar", s.Queries[1])
	require.Equal(t, 100, *s.Filter.RetweetedThreshold)
	require.True(t, s.Action.Twitter.Retweet)

	ch := c.GetSlackMessages()[0].Channels[0]
	require.Equal(t, "foo", ch)
	n := c.Twitter.Notification
	require.True(t, n.Place.AllowSelf)
	require.Equal(t, "foo", n.Place.Users[0])

	clone := *c
	require.NoError(t, clone.Validate())
	require.True(t, reflect.DeepEqual(&clone, c))
}

func TestNewConfigWhenNotExist(t *testing.T) {
	_, err := NewFileConfig("config_not_exist.toml")
	require.NoError(t, err)
}

func TestSaveLoad(t *testing.T) {
	c, err := NewFileConfig("testdata/config.template.toml")
	require.NoError(t, err)
	dir, err := ioutil.TempDir(os.TempDir(), "mybot_")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	jsonCfg := *c
	jsonCfg.File = filepath.Join(dir, "config.json")
	err = jsonCfg.Save()
	require.NoError(t, err)
	err = jsonCfg.Load()
	require.NoError(t, err)
	jsonCfg.File = c.File
	require.True(t, reflect.DeepEqual(&jsonCfg, c))

	tomlCfg := *c
	tomlCfg.File = filepath.Join(dir, "config.toml")
	err = tomlCfg.Save()
	require.NoError(t, err)
	err = tomlCfg.Load()
	require.NoError(t, err)
	tomlCfg.File = c.File
	require.True(t, reflect.DeepEqual(&tomlCfg, c))
}

func TestFileConfigTwitterTimelines(t *testing.T) {
	c, err := NewFileConfig("")
	require.NoError(t, err)
	testConfigTwitterTimelines(t, c)
}

func testConfigTwitterTimelines(t *testing.T, c Config) {
	action := data.NewAction()
	action.Twitter.Retweet = true
	timeline := TimelineConfig{}
	timeline.Action = action
	timeline.ScreenNames = []string{"foo"}
	timelines := []TimelineConfig{timeline}
	c.SetTwitterTimelines(timelines)
	ts := c.GetTwitterTimelines()
	require.True(t, reflect.DeepEqual(timelines, ts))
}

func TestFileConfigTwitterFavorites(t *testing.T) {
	c, err := NewFileConfig("")
	require.NoError(t, err)
	testConfigTwitterFavorites(t, c)
}

func testConfigTwitterFavorites(t *testing.T, c Config) {
	action := data.NewAction()
	action.Twitter.Retweet = true
	favorite := FavoriteConfig{}
	favorite.Action = action
	favorite.ScreenNames = []string{"foo"}
	favorites := []FavoriteConfig{favorite}
	c.SetTwitterFavorites(favorites)
	fs := c.GetTwitterFavorites()
	require.True(t, reflect.DeepEqual(favorites, fs))
}

func TestFileConfigTwitterSearches(t *testing.T) {
	c, err := NewFileConfig("")
	require.NoError(t, err)
	testConfigTwitterSearches(t, c)
}

func testConfigTwitterSearches(t *testing.T, c Config) {
	action := data.NewAction()
	action.Twitter.Retweet = true
	search := SearchConfig{}
	search.Action = action
	search.Queries = []string{"foo"}
	searches := []SearchConfig{search}
	c.SetTwitterSearches(searches)
	ss := c.GetTwitterSearches()
	require.True(t, reflect.DeepEqual(searches, ss))
}

func TestFileConfigTwitterNotification(t *testing.T) {
	c, err := NewFileConfig("")
	require.NoError(t, err)
	testConfigTwitterNotification(t, c)
}

func testConfigTwitterNotification(t *testing.T, c Config) {
	notification := Notification{
		Place: PlaceNotification{
			Users: []string{"foo"},
		},
	}
	c.SetTwitterNotification(notification)
	n := c.GetTwitterNotification()
	require.True(t, reflect.DeepEqual(notification, n))
}

func TestFileConfigSlackMessages(t *testing.T) {
	c, err := NewFileConfig("")
	require.NoError(t, err)
	testConfigSlackMessages(t, c)
}

func testConfigSlackMessages(t *testing.T, c Config) {
	filter := NewFilter()
	filter.Lang = "ja"
	action := data.NewAction()
	action.Slack.Channels = []string{"foo"}
	action.Slack.Reactions = []string{":smile:"}
	msg := MessageConfig{
		Channels: []string{"foo"},
	}
	msg.Filter = filter
	msg.Action = action
	msgs := []MessageConfig{msg}

	c.SetSlackMessages(msgs)
	ms := c.GetSlackMessages()
	require.True(t, reflect.DeepEqual(msgs, ms))
}

func TestFileConfigInteraction(t *testing.T) {
	c, err := NewFileConfig("")
	require.NoError(t, err)
	testConfigInteraction(t, c)
}

func testConfigInteraction(t *testing.T, c Config) {
	interaction := InteractionConfig{}
	interaction.Users = []string{"foo"}
	c.SetTwitterInteraction(interaction)
	i := c.GetTwitterInteraction()
	require.True(t, reflect.DeepEqual(interaction, i))
}

func TestFileConfigTwitterDuration(t *testing.T) {
	c, err := NewFileConfig("")
	require.NoError(t, err)
	testConfigTwitterDuration(t, c)
}

func testConfigTwitterDuration(t *testing.T, c Config) {
	duration := "20m"
	c.SetTwitterDuration(duration)
	dur := c.GetTwitterDuration()
	require.Equal(t, duration, dur)
}
