package core_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/core"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/mocks"
	"github.com/stretchr/testify/require"
)

const (
	defaultTestConfigFilePath = "testdata/config.yaml"
)

func TestNewFileConfig(t *testing.T) {
	c := core.NewTestFileConfig(defaultTestConfigFilePath, t)

	a := c.GetTwitterTimelines()[0]
	require.Equal(t, "Golang Release", a.Name)
	require.Equal(t, "golang", a.ScreenNames[0])
	f := a.Filter
	require.Equal(t, "is released!", f.Patterns[0])
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

	require.NoError(t, c.Validate())
}

func TestNewFileConfig_withUnexistingFile(t *testing.T) {
	_, err := core.NewFileConfig("config_not_exist.toml")
	require.NoError(t, err)
}

func TestNewFileConfig_withInvalidJsonFile(t *testing.T) {
	fpath := filepath.Join(os.TempDir(), "invalidformat.json")
	defer func() { require.NoError(t, os.RemoveAll(fpath)) }()
	require.NoError(t, ioutil.WriteFile(fpath, []byte("foo"), os.FileMode(0777)))
	_, err := core.NewFileConfig(fpath)
	require.Error(t, err)
}

func TestNewFileConfig_withInvalidTomlFile(t *testing.T) {
	fpath := filepath.Join(os.TempDir(), "invalidformat.toml")
	defer func() { require.NoError(t, os.RemoveAll(fpath)) }()
	require.NoError(t, ioutil.WriteFile(fpath, []byte("[[[]]]"), os.FileMode(0777)))
	_, err := core.NewFileConfig(fpath)
	require.Error(t, err)
}

func TestNewFileConfig_withUnknownFileExtension(t *testing.T) {
	fpath := filepath.Join(os.TempDir(), "invalid.txt")
	require.NoError(t, ioutil.WriteFile(fpath, []byte(""), os.FileMode(0777)))
	_, err := core.NewFileConfig(fpath)
	defer func() { require.NoError(t, os.RemoveAll(fpath)) }()
	require.Error(t, err)
}

func TestNewFileConfig_withInvalidDataFile(t *testing.T) {
	c, err := core.NewFileConfig(defaultTestConfigFilePath)
	require.NoError(t, err)
	c.GetTwitterTimelines()[0].ScreenNames = []string{}
	bs, err := json.Marshal(c)
	require.NoError(t, err)
	fpath := filepath.Join(os.TempDir(), "invaliddata.json")
	defer func() { require.NoError(t, os.RemoveAll(fpath)) }()
	require.NoError(t, ioutil.WriteFile(fpath, bs, os.FileMode(0777)))
	_, err = core.NewFileConfig(fpath)
	require.Error(t, err)
}

func TestFileConfig_SaveLoad_withJson(t *testing.T) {
	defaultConfig, err := core.NewFileConfig(defaultTestConfigFilePath)
	require.NoError(t, err)
	c := core.FileConfig{
		ConfigProperties: defaultConfig.GetProperties(),
		File:             filepath.Join(os.TempDir(), "config.json"),
	}
	require.NoError(t, c.Save())
	require.NoError(t, c.Load())
	require.Equal(t, defaultConfig.GetProperties(), c.GetProperties())
}

func TestFileConfig_SaveLoad_withToml(t *testing.T) {
	defaultConfig, err := core.NewFileConfig(defaultTestConfigFilePath)
	require.NoError(t, err)
	c := core.FileConfig{
		ConfigProperties: defaultConfig.GetProperties(),
		File:             filepath.Join(os.TempDir(), "config.toml"),
	}
	require.NoError(t, c.Save())
	require.NoError(t, c.Load())
	require.Equal(t, defaultConfig.GetProperties(), c.GetProperties())
}

func TestFileConfig_Save_withUnknownExtension(t *testing.T) {
	c, err := core.NewFileConfig(defaultTestConfigFilePath)
	require.NoError(t, err)
	c.File = filepath.Join(os.TempDir(), "invalid_extension.txt")
	require.Error(t, c.Save())
}

func TestFileConfig_Save_withUnwritableFile(t *testing.T) {
	c, err := core.NewFileConfig(defaultTestConfigFilePath)
	require.NoError(t, err)
	fpath := filepath.Join(os.TempDir(), "unwritable.json")
	_, err = os.Create(fpath)
	require.NoError(t, err)
	require.NoError(t, os.Chmod(fpath, 0440))
	c.File = fpath
	defer func() { require.NoError(t, c.Delete()) }()
	require.Error(t, c.Save())
}

func TestFileConfig_Load_withUnreadableFile(t *testing.T) {
	fpath := filepath.Join(os.TempDir(), "unreadable.json")
	require.NoError(t, ioutil.WriteFile(fpath, []byte("{}"), 0640))
	c, err := core.NewFileConfig(fpath)
	require.NoError(t, err)
	require.NoError(t, os.Chmod(fpath, 0220))
	defer func() { require.NoError(t, c.Delete()) }()
	require.Error(t, c.Load())
}

func TestNewDBConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	col := mocks.NewMockMgoCollection(ctrl)
	query := mocks.NewMockMgoQuery(ctrl)
	col.EXPECT().Find(gomock.Any()).Return(query)
	query.EXPECT().Count().Return(0, nil)
	c, err := core.NewDBConfig(col, "foo")
	require.NoError(t, err)
	col.EXPECT().RemoveAll(gomock.Any()).Return(nil, nil)
	require.NoError(t, c.Delete())
}

func TestNewDBConfig_withCountError(t *testing.T) {
	ctrl := gomock.NewController(t)
	col := mocks.NewMockMgoCollection(ctrl)
	query := mocks.NewMockMgoQuery(ctrl)
	col.EXPECT().Find(gomock.Any()).Return(query)
	query.EXPECT().Count().Return(0, errors.New("error"))
	_, err := core.NewDBConfig(col, "foo")
	require.Error(t, err)
}

func TestNewDBConfig_withExistingData(t *testing.T) {
	ctrl := gomock.NewController(t)
	col := mocks.NewMockMgoCollection(ctrl)
	query := mocks.NewMockMgoQuery(ctrl)
	col.EXPECT().Find(gomock.Any()).Return(query)
	query.EXPECT().Count().Return(1, nil)
	query.EXPECT().One(gomock.Any()).Return(nil)
	_, err := core.NewDBConfig(col, "foo")
	require.NoError(t, err)
}

func TestDBConfig_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	col := mocks.NewMockMgoCollection(ctrl)
	query := mocks.NewMockMgoQuery(ctrl)
	col.EXPECT().Find(gomock.Any()).Return(query)
	query.EXPECT().Count().Return(0, nil)
	col.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil, nil)
	c, err := core.NewDBConfig(col, "foo")
	require.NoError(t, err)
	require.NoError(t, c.Save())
}

func TestConfig_GetConfigProperties(t *testing.T) {
	c := core.NewTestFileConfig(defaultTestConfigFilePath, t)
	props := c.GetProperties()
	require.Equal(t, &c.ConfigProperties, &props)
}

func TestConfig_GetTwitterScreenNames(t *testing.T) {
	c := core.NewTestFileConfig(defaultTestConfigFilePath, t)
	require.Equal(t, []string{"golang", "foo"}, c.GetTwitterScreenNames())
}

func TestConfig_GetTwitterTimelinesByScreenName(t *testing.T) {
	c := core.NewTestFileConfig(defaultTestConfigFilePath, t)
	require.Len(t, c.GetTwitterTimelinesByScreenName("golang"), 1)
}

func TestConfig_AddTwitterTimeline(t *testing.T) {
	c := core.NewTestFileConfig(defaultTestConfigFilePath, t)
	tc := core.TimelineConfig{}
	c.AddTwitterTimeline(tc)
	tcs := c.GetTwitterTimelines()
	require.Equal(t, 3, len(tcs))
	require.Equal(t, tc, tcs[2])
}

func TestConfig_AddTwitterFavorite(t *testing.T) {
	c := core.NewTestFileConfig(defaultTestConfigFilePath, t)
	fav := core.FavoriteConfig{}
	c.AddTwitterFavorite(fav)
	favs := c.GetTwitterFavorites()
	require.Equal(t, 2, len(favs))
	require.Equal(t, fav, favs[1])
}

func TestConfig_AddTwitterSearch(t *testing.T) {
	c := core.NewTestFileConfig(defaultTestConfigFilePath, t)
	s := core.SearchConfig{}
	c.AddTwitterSearch(s)
	ss := c.GetTwitterSearches()
	require.Equal(t, 2, len(ss))
	require.Equal(t, s, ss[1])
}

func TestConfig_AddSlackMessage(t *testing.T) {
	c := core.NewTestFileConfig(defaultTestConfigFilePath, t)
	msg := core.MessageConfig{}
	c.AddSlackMessage(msg)
	msgs := c.GetSlackMessages()
	require.Equal(t, 2, len(msgs))
	require.Equal(t, msg, msgs[1])
}

func TestConfig_Validate(t *testing.T) {
	var c core.Config

	c = core.NewTestFileConfig(defaultTestConfigFilePath, t)
	c.GetTwitterSearches()[0].Queries = []string{}
	require.Error(t, c.Validate())

	c = core.NewTestFileConfig(defaultTestConfigFilePath, t)
	c.GetSlackMessages()[0].Channels = []string{}
	require.Error(t, c.Validate())
}

func TestConfig_ValidateWithAPI(t *testing.T) {
	c := core.NewTestFileConfig(defaultTestConfigFilePath, t)
	var api *mocks.MockTwitterAPI
	ctrl := gomock.NewController(t)

	api = mocks.NewMockTwitterAPI(ctrl)
	api.EXPECT().GetUsersShow("golang", gomock.Any()).Return(anaconda.User{}, errors.New("foo"))
	require.Error(t, c.ValidateWithAPI(api))

	api = mocks.NewMockTwitterAPI(ctrl)
	api.EXPECT().GetUsersShow(gomock.Any(), gomock.Any()).AnyTimes().Return(anaconda.User{}, nil)
	require.NoError(t, c.ValidateWithAPI(api))
}

func TestConfig_Unmarshal(t *testing.T) {
	fpath := filepath.Join(os.TempDir(), "unexisting.yaml")
	c, err := core.NewFileConfig(fpath)
	require.NoError(t, err)
	bs, err := ioutil.ReadFile(defaultTestConfigFilePath)
	require.NoError(t, err)
	require.NoError(t, c.Unmarshal(".yaml", bs))
	require.Equal(t, "golang", c.GetTwitterTimelines()[0].ScreenNames[0])
}

func TestFileConfig_TwitterTimelines(t *testing.T) {
	c, err := core.NewFileConfig("")
	require.NoError(t, err)
	testConfigTwitterTimelines(t, c)
}

func testConfigTwitterTimelines(t *testing.T, c core.Config) {
	timeline := core.TimelineConfig{}
	action := data.NewAction()
	action.Twitter.Retweet = true
	timeline.Action = action
	timeline.ScreenNames = []string{"foo"}
	timelines := []core.TimelineConfig{timeline}

	c.SetTwitterTimelines(timelines)
	require.Equal(t, timelines, c.GetTwitterTimelines())
}

func TestFileConfig_TwitterFavorites(t *testing.T) {
	c, err := core.NewFileConfig("")
	require.NoError(t, err)
	testConfigTwitterFavorites(t, c)
}

func testConfigTwitterFavorites(t *testing.T, c core.Config) {
	favorite := core.FavoriteConfig{}
	action := data.NewAction()
	action.Twitter.Retweet = true
	favorite.Action = action
	favorite.ScreenNames = []string{"foo"}
	favorites := []core.FavoriteConfig{favorite}

	c.SetTwitterFavorites(favorites)
	require.Equal(t, favorites, c.GetTwitterFavorites())
}

func TestFileConfig_TwitterSearches(t *testing.T) {
	c, err := core.NewFileConfig("")
	require.NoError(t, err)
	testConfigTwitterSearches(t, c)
}

func testConfigTwitterSearches(t *testing.T, c core.Config) {
	search := core.SearchConfig{}
	action := data.NewAction()
	action.Twitter.Retweet = true
	search.Action = action
	search.Queries = []string{"foo"}
	searches := []core.SearchConfig{search}

	c.SetTwitterSearches(searches)
	require.Equal(t, searches, c.GetTwitterSearches())
}

func TestFileConfig_SlackMessages(t *testing.T) {
	c, err := core.NewFileConfig("")
	require.NoError(t, err)
	testConfigSlackMessages(t, c)
}

func testConfigSlackMessages(t *testing.T, c core.Config) {
	filter := core.NewFilter()
	filter.Lang = "ja"
	action := data.NewAction()
	action.Slack.Channels = []string{"foo"}
	action.Slack.Reactions = []string{":smile:"}
	msg := core.MessageConfig{
		Channels: []string{"foo"},
	}
	msg.Filter = filter
	msg.Action = action
	msgs := []core.MessageConfig{msg}

	c.SetSlackMessages(msgs)
	ms := c.GetSlackMessages()
	require.Equal(t, msgs, ms)
}

func TestFileConfig_TwitterDuration(t *testing.T) {
	c, err := core.NewFileConfig("")
	require.NoError(t, err)
	testConfigTwitterDuration(t, c)
}

func testConfigTwitterDuration(t *testing.T, c core.Config) {
	duration := "20m"
	c.SetPollingDuration(duration)
	dur := c.GetPollingDuration()
	require.Equal(t, duration, dur)
}

func TestSource_Validate(t *testing.T) {
	s := core.NewSource()
	require.Error(t, s.Validate())
	s.Name = "foo"
	require.Error(t, s.Validate())
}

func TestTimelineConfig_Validate(t *testing.T) {
	tc := core.NewTimelineConfig()
	require.Error(t, tc.Validate())
}

func TestFavoriteConfig_Validate(t *testing.T) {
	fav := core.NewFavoriteConfig()
	require.Error(t, fav.Validate())
}

func TestSearchConfig_Validate(t *testing.T) {
	s := core.NewSearchConfig()
	require.Error(t, s.Validate())
}

func TestMessageConfig_Validate(t *testing.T) {
	msg := core.NewMessageConfig()
	require.Error(t, msg.Validate())
}
