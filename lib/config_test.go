package mybot_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/data"
	. "github.com/iwataka/mybot/lib"
	"github.com/iwataka/mybot/mocks"
	"github.com/stretchr/testify/require"
)

const (
	defaultTestConfigFilePath = "testdata/config.template.toml"
)

func TestNewConfig(t *testing.T) {
	c := NewTestFileConfig(defaultTestConfigFilePath, t)

	a := c.GetTwitterTimelines()[0]
	require.Equal(t, "Golang Release", a.Name)
	require.Equal(t, "golang", a.ScreenNames[0])
	f := a.Filter
	require.Equal(t, "is released!", f.Patterns[0])
	require.False(t, *f.Retweeted)
	require.Equal(t, "en", f.Lang)
	require.Equal(t, "cartoon|clip art|artwork", f.Vision.Label[0])
	require.True(t, a.Action.Twitter.Retweet)
	require.Equal(t, "foo", a.Action.Slack.Channels[0])

	s := c.GetTwitterSearches()[0]
	require.Equal(t, "foo bar", s.Name)
	require.Equal(t, "foo", s.Queries[0])
	require.Equal(t, "bar", s.Queries[1])
	require.Equal(t, 100, *s.Filter.RetweetedThreshold)
	require.True(t, s.Action.Twitter.Retweet)

	msg := c.GetSlackMessages()[0]
	require.Equal(t, "foo", msg.Name)
	require.Equal(t, "foo", msg.Channels[0])

	n := c.Twitter.Notification
	require.True(t, n.Place.AllowSelf)
	require.Equal(t, "foo", n.Place.Users[0])

	clone := *c
	require.NoError(t, clone.Validate())
	require.Equal(t, &clone, c)
}

func TestConfigGetConfigProperties(t *testing.T) {
	c := NewTestFileConfig(defaultTestConfigFilePath, t)
	require.Equal(t, c.ConfigProperties, c.GetProperties())
}

func TestConfigGetTwitterScreenNames(t *testing.T) {
	c := NewTestFileConfig(defaultTestConfigFilePath, t)
	require.Equal(t, []string{"golang", "foo"}, c.GetTwitterScreenNames())
}

func TestAddTwitterTimeline(t *testing.T) {
	c := NewTestFileConfig(defaultTestConfigFilePath, t)
	tc := TimelineConfig{}
	c.AddTwitterTimeline(tc)
	tcs := c.GetTwitterTimelines()
	require.Equal(t, 3, len(tcs))
	require.Equal(t, tc, tcs[2])
}

func TestAddTwitterFavorite(t *testing.T) {
	c := NewTestFileConfig(defaultTestConfigFilePath, t)
	fav := FavoriteConfig{}
	c.AddTwitterFavorite(fav)
	favs := c.GetTwitterFavorites()
	require.Equal(t, 1, len(favs))
	require.Equal(t, fav, favs[0])
}

func TestAddTwitterSearch(t *testing.T) {
	c := NewTestFileConfig(defaultTestConfigFilePath, t)
	s := SearchConfig{}
	c.AddTwitterSearch(s)
	ss := c.GetTwitterSearches()
	require.Equal(t, 2, len(ss))
	require.Equal(t, s, ss[1])
}

func TestAddSlackMessage(t *testing.T) {
	c := NewTestFileConfig(defaultTestConfigFilePath, t)
	msg := MessageConfig{}
	c.AddSlackMessage(msg)
	msgs := c.GetSlackMessages()
	require.Equal(t, 2, len(msgs))
	require.Equal(t, msg, msgs[1])
}

func TestNewConfigWhenNotExist(t *testing.T) {
	_, err := NewFileConfig("config_not_exist.toml")
	require.NoError(t, err)
}

func TestNewConfigForInvalidFormatJSONFile(t *testing.T) {
	fname := "testdata/invalidformat.json"
	err := ioutil.WriteFile(fname, []byte("foo"), os.FileMode(0777))
	require.NoError(t, err)
	defer os.Remove(fname)
	_, err = NewFileConfig(fname)
	require.Error(t, err)
}

func TestNewConfigForInvalidFormatTomlFile(t *testing.T) {
	fname := "testdata/invalidformat.toml"
	err := ioutil.WriteFile(fname, []byte("[[[]]]"), os.FileMode(0777))
	require.NoError(t, err)
	defer os.Remove(fname)
	_, err = NewFileConfig(fname)
	require.Error(t, err)
}

func TestNewConfigForUnknownExtensionFile(t *testing.T) {
	fname := "testdata/invalid.txt"
	err := ioutil.WriteFile(fname, []byte(""), os.FileMode(0777))
	require.NoError(t, err)
	defer os.Remove(fname)
	_, err = NewFileConfig(fname)
	require.Error(t, err)
}

func TestNewConfigForInvalidDataFile(t *testing.T) {
	c, err := NewFileConfig(defaultTestConfigFilePath)
	c.GetTwitterTimelines()[0].ScreenNames = []string{}
	bs, err := json.Marshal(c)
	require.NoError(t, err)
	fname := "testdata/invaliddata.json"
	err = ioutil.WriteFile(fname, bs, os.FileMode(0777))
	require.NoError(t, err)
	defer os.Remove(fname)
	_, err = NewFileConfig(fname)
	require.Error(t, err)
}

func TestConfigValidate(t *testing.T) {
	var c Config

	c = NewTestFileConfig(defaultTestConfigFilePath, t)
	c.GetTwitterSearches()[0].Queries = []string{}
	require.Error(t, c.Validate())

	c = NewTestFileConfig(defaultTestConfigFilePath, t)
	c.GetSlackMessages()[0].Channels = []string{}
	require.Error(t, c.Validate())
}

func TestValidateWithAPI(t *testing.T) {
	c := NewTestFileConfig(defaultTestConfigFilePath, t)
	var api *mocks.MockTwitterAPI
	ctrl := gomock.NewController(t)

	api = mocks.NewMockTwitterAPI(ctrl)
	api.EXPECT().GetUsersShow("golang", gomock.Any()).Return(anaconda.User{}, errors.New("foo"))
	require.Error(t, c.ValidateWithAPI(api))

	api = mocks.NewMockTwitterAPI(ctrl)
	api.EXPECT().GetUsersShow(gomock.Any(), gomock.Any()).AnyTimes().Return(anaconda.User{}, nil)
	require.NoError(t, c.ValidateWithAPI(api))
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
	require.Equal(t, &jsonCfg, c)

	tomlCfg := *c
	tomlCfg.File = filepath.Join(dir, "config.toml")
	err = tomlCfg.Save()
	require.NoError(t, err)
	err = tomlCfg.Load()
	require.NoError(t, err)
	tomlCfg.File = c.File
	require.Equal(t, &tomlCfg, c)
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
	require.Equal(t, timelines, ts)
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
	require.Equal(t, favorites, fs)
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
	require.Equal(t, searches, ss)
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
	require.Equal(t, notification, n)
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
	require.Equal(t, msgs, ms)
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
	require.Equal(t, interaction, i)
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

func TestSourceValidate(t *testing.T) {
	s := NewSource()
	require.Error(t, s.Validate())
	s.Name = "foo"
	require.Error(t, s.Validate())
}

func TestTimelineConfigValidate(t *testing.T) {
	tc := NewTimelineConfig()
	require.Error(t, tc.Validate())
}

func TestFavoriteConfigValidate(t *testing.T) {
	fav := NewFavoriteConfig()
	require.Error(t, fav.Validate())
}

func TestSearchConfigValidate(t *testing.T) {
	s := NewSearchConfig()
	require.Error(t, s.Validate())
}

func TestMessageConfigValidate(t *testing.T) {
	msg := NewMessageConfig()
	require.Error(t, msg.Validate())
}
