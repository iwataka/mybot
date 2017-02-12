package mybot

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	c, err := NewFileConfig("test_assets/config.template.toml")
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	a := c.Twitter.Timelines[0]
	if a.ScreenNames[0] != "golang" {
		t.Fatalf("%s expected but %s found", "golang", a.ScreenNames[0])
	}
	f := a.Filter
	if f.Patterns[0] != "is released!" {
		t.Fatalf("%s expected but %s found", "is released!", f.Patterns[0])
	}
	if *f.HasURL != true {
		t.Fatalf("%v expected but %v found", true, *f.HasURL)
	}
	if *f.Retweeted != false {
		t.Fatalf("%v expected but %v found", false, *f.Retweeted)
	}
	if f.Lang != "en" {
		t.Fatalf("%s expected but %s found", "en", f.Lang)
	}
	if f.Vision.Label[0] != "cartoon|clip art|artwork" {
		t.Fatalf("%s expected but %s found", "cartoon|clip art|artwork", f.Vision.Label[0])
	}
	if a.Action.Twitter.Retweet != true {
		t.Fatalf("%v expected but %v found", true, a.Action.Twitter.Retweet)
	}
	if a.Action.Slack.Channels[0] != "foo" {
		t.Fatalf("%v expected but %v found", "foo", a.Action.Slack.Channels[0])
	}
	s := c.Twitter.Searches[0]
	if s.Queries[0] != "foo" {
		t.Fatalf("%s expected but %s found", "foo", s.Queries[0])
	}
	if s.Queries[1] != "bar" {
		t.Fatalf("%s expected but %s found", "bar", s.Queries[1])
	}
	if *s.Filter.RetweetedThreshold != 100 {
		t.Fatalf("%d expected but %d found", 100, *s.Filter.RetweetedThreshold)
	}
	if s.Action.Twitter.Retweet != true {
		t.Fatalf("%v expected but %v found", true, s.Action.Twitter.Retweet)
	}
	n := c.Twitter.Notification
	if n.Place.AllowSelf != true {
		t.Fatalf("%v expected but %v found", true, n.Place.AllowSelf)
	}
	if n.Place.Users[0] != "foo" {
		t.Fatalf("%s expected but %s found", "foo", n.Place.Users[0])
	}
	if !c.Log.AllowSelf {
		t.Fatalf("%v expected but %v found", true, c.Log.AllowSelf)
	}
	if c.Log.Users[0] != "foo" {
		t.Fatalf("%s expected but %s found", "foo", c.Log.Users[0])
	}
	if c.Log.Users[1] != "bar" {
		t.Fatalf("%s expected but %s found", "bar", c.Log.Users[1])
	}
	if c.Log.Linenum != 8 {
		t.Fatalf("%v expected but %v found", 8, c.Log.Linenum)
	}

	clone := *c
	err = clone.Validate()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&clone, c) {
		t.Fatalf("%v expected but %v found", c, &clone)
	}
}

func TestNewConfigWhenNotExist(t *testing.T) {
	_, err := NewFileConfig("config_not_exist.toml")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewTimelineConfig(t *testing.T) {
	tl := NewTimelineConfig()
	if tl.Filter == nil {
		t.Fatalf("Non-nil expected but nil found")
	}
	if tl.Filter.Vision == nil {
		t.Fatalf("Non-nil expected but nil found")
	}
	if tl.Action == nil {
		t.Fatalf("Non-nil expected but nil found")
	}
}

func TestNewFavoriteConfig(t *testing.T) {
	f := NewFavoriteConfig()
	if f.Filter == nil {
		t.Fatalf("Non-nil expected but nil found")
	}
	if f.Filter.Vision == nil {
		t.Fatalf("Non-nil expected but nil found")
	}
	if f.Action == nil {
		t.Fatalf("Non-nil expected but nil found")
	}
}

func TestNewSearchConfig(t *testing.T) {
	s := NewSearchConfig()
	if s.Filter == nil {
		t.Fatalf("Non-nil expected but nil found")
	}
	if s.Filter.Vision == nil {
		t.Fatalf("Non-nil expected but nil found")
	}
	if s.Action == nil {
		t.Fatalf("Non-nil expected but nil found")
	}
}

func TestTweetAction_AddNil(t *testing.T) {
	a1 := &TweetAction{
		Twitter: nil,
		Slack:   nil,
	}
	a2 := &TweetAction{
		Twitter: &TwitterAction{
			Collections: []string{"foo"},
		},
		Slack: &SlackAction{
			Channels: []string{"bar"},
		},
	}
	a2.Twitter.Retweet = true

	result1 := *a1
	result1.Add(a2)
	if !reflect.DeepEqual(&result1, a2) {
		t.Fatalf("Failed to add %v to %v: %v", a2, a1, result1)
	}

	result2 := *a2
	result2.Add(a1)
	if !reflect.DeepEqual(&result2, a2) {
		t.Fatalf("Failed to add %v to %v: %v", a1, a2, result2)
	}
}

func TestTweetAction_Add(t *testing.T) {
	a1 := &TweetAction{
		Twitter: &TwitterAction{
			Collections: []string{"twitter"},
		},
		Slack: &SlackAction{
			Channels: []string{"slack"},
		},
	}
	a1.Twitter.Retweet = true
	a2 := &TweetAction{
		Twitter: &TwitterAction{
			Collections: []string{"facebook"},
		},
		Slack: &SlackAction{
			Channels: []string{"mattermost"},
		},
	}
	a2.Twitter.Favorite = true
	a1.Add(a2)

	if !a1.Twitter.Retweet {
		t.Fatalf("%v expected but %v found", true, a1.Twitter.Retweet)
	}
	if !a1.Twitter.Favorite {
		t.Fatalf("%v expected but %v found", true, a1.Twitter.Favorite)
	}
	if a1.Twitter.Follow {
		t.Fatalf("%v expected but %v found", false, a1.Twitter.Follow)
	}
	if len(a1.Twitter.Collections) != 2 {
		t.Fatalf("%d expected but %d found", 2, len(a1.Twitter.Collections))
	}
	if len(a1.Slack.Channels) != 2 {
		t.Fatalf("%d expected but %d found", 2, len(a1.Slack.Channels))
	}
}

func TestSaveLoad(t *testing.T) {
	c, err := NewFileConfig("test_assets/config.template.toml")
	if err != nil {
		t.Fatal(err)
	}
	dir, err := ioutil.TempDir(os.TempDir(), "mybot_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	jsonCfg := *c
	jsonCfg.File = filepath.Join(dir, "config.json")
	err = jsonCfg.Save()
	if err != nil {
		t.Fatal(err)
	}
	err = jsonCfg.Load()
	if err != nil {
		t.Fatal(err)
	}
	jsonCfg.File = c.File
	if !reflect.DeepEqual(&jsonCfg, c) {
		t.Fatalf("%v expected but %v found", c, jsonCfg)
	}

	tomlCfg := *c
	tomlCfg.File = filepath.Join(dir, "config.toml")
	err = tomlCfg.Save()
	if err != nil {
		t.Fatal(err)
	}
	err = tomlCfg.Load()
	if err != nil {
		t.Fatal(err)
	}
	tomlCfg.File = c.File
	if !reflect.DeepEqual(&tomlCfg, c) {
		t.Fatalf("%v expected but %v found", c, tomlCfg)
	}
}

func TestFileConfigTwitterTimelines(t *testing.T) {
	c, err := NewFileConfig("")
	if err != nil {
		t.Fatal(err)
	}
	testConfigTwitterTimelines(t, c)
}

func testConfigTwitterTimelines(t *testing.T, c Config) {
	action := NewTweetAction()
	action.Twitter.Retweet = true
	timeline := TimelineConfig{}
	timeline.Action = action
	timeline.ScreenNames = []string{"foo"}
	timelines := []TimelineConfig{timeline}
	err := c.SetTwitterTimelines(timelines)
	if err != nil {
		t.Fatal(err)
	}
	ts, err := c.GetTwitterTimelines()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(timelines, ts) {
		t.Fatalf("%v is not set properly", timelines)
	}
}

func TestFileConfigTwitterFavorites(t *testing.T) {
	c, err := NewFileConfig("")
	if err != nil {
		t.Fatal(err)
	}
	testConfigTwitterFavorites(t, c)
}

func testConfigTwitterFavorites(t *testing.T, c Config) {
	action := NewTweetAction()
	action.Twitter.Retweet = true
	favorite := FavoriteConfig{}
	favorite.Action = action
	favorite.ScreenNames = []string{"foo"}
	favorites := []FavoriteConfig{favorite}
	err := c.SetTwitterFavorites(favorites)
	if err != nil {
		t.Fatal(err)
	}
	fs, err := c.GetTwitterFavorites()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(favorites, fs) {
		t.Fatalf("%v is not set properly", favorites)
	}
}

func TestFileConfigTwitterSearches(t *testing.T) {
	c, err := NewFileConfig("")
	if err != nil {
		t.Fatal(err)
	}
	testConfigTwitterSearches(t, c)
}

func testConfigTwitterSearches(t *testing.T, c Config) {
	action := NewTweetAction()
	action.Twitter.Retweet = true
	search := SearchConfig{}
	search.Action = action
	search.Queries = []string{"foo"}
	searches := []SearchConfig{search}
	err := c.SetTwitterSearches(searches)
	if err != nil {
		t.Fatal(err)
	}
	ss, err := c.GetTwitterSearches()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(searches, ss) {
		t.Fatalf("%v is not set properly", searches)
	}
}

func TestFileConfigTwitterAPIs(t *testing.T) {
	c, err := NewFileConfig("")
	if err != nil {
		t.Fatal(err)
	}
	testConfigTwitterAPIs(t, c)
}

func testConfigTwitterAPIs(t *testing.T, c Config) {
	apis := []APIConfig{
		APIConfig{
			SourceURL:       "foo",
			MessageTemplate: "bar",
		},
	}
	err := c.SetTwitterAPIs(apis)
	if err != nil {
		t.Fatal(err)
	}
	as, err := c.GetTwitterAPIs()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(apis, as) {
		t.Fatalf("%v is not set properly", apis)
	}
}

func TestFileConfigTwitterNotification(t *testing.T) {
	c, err := NewFileConfig("")
	if err != nil {
		t.Fatal(err)
	}
	testConfigTwitterNotification(t, c)
}

func testConfigTwitterNotification(t *testing.T, c Config) {
	notification := &Notification{
		Place: &PlaceNotification{
			Users: []string{"foo"},
		},
	}
	err := c.SetTwitterNotification(notification)
	if err != nil {
		t.Fatal(err)
	}
	n, err := c.GetTwitterNotification()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(notification, n) {
		t.Fatalf("%v is not set properly", notification)
	}
}

func TestFileConfigInteraction(t *testing.T) {
	c, err := NewFileConfig("")
	if err != nil {
		t.Fatal(err)
	}
	testConfigInteraction(t, c)
}

func testConfigInteraction(t *testing.T, c Config) {
	interaction := &InteractionConfig{}
	interaction.Users = []string{"foo"}
	err := c.SetInteraction(interaction)
	if err != nil {
		t.Fatal(err)
	}
	i, err := c.GetInteraction()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(interaction, i) {
		t.Fatalf("%v is not set properly", interaction)
	}
}

func TestFileConfigLog(t *testing.T) {
	c, err := NewFileConfig("")
	if err != nil {
		t.Fatal(err)
	}
	testConfigLog(t, c)
}

func testConfigLog(t *testing.T, c Config) {
	log := &LogConfig{}
	log.Users = []string{"foo"}
	err := c.SetLog(log)
	if err != nil {
		t.Fatal(err)
	}
	l, err := c.GetLog()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(log, l) {
		t.Fatalf("%v is not set properly", log)
	}
}

func TestFileConfigTwitterDuration(t *testing.T) {
	c, err := NewFileConfig("")
	if err != nil {
		t.Fatal(err)
	}
	testConfigTwitterDuration(t, c)
}

func testConfigTwitterDuration(t *testing.T, c Config) {
	duration := "20m"
	err := c.SetTwitterDuration(duration)
	if err != nil {
		t.Fatal(err)
	}
	dur, err := c.GetTwitterDuration()
	if err != nil {
		t.Fatal(err)
	}
	if duration != dur {
		t.Fatalf("%v is not set properly", duration)
	}
}
