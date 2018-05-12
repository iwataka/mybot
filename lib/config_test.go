package mybot

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/iwataka/mybot/data"
)

func TestNewConfig(t *testing.T) {
	c, err := NewFileConfig("testdata/config.template.toml")
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
	ch := c.Slack.Messages[0].Channels[0]
	if ch != "foo" {
		t.Fatalf("%s expected but %v found", "foo", ch)
	}
	n := c.Twitter.Notification
	if n.Place.AllowSelf != true {
		t.Fatalf("%v expected but %v found", true, n.Place.AllowSelf)
	}
	if n.Place.Users[0] != "foo" {
		t.Fatalf("%s expected but %s found", "foo", n.Place.Users[0])
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

func TestSaveLoad(t *testing.T) {
	c, err := NewFileConfig("testdata/config.template.toml")
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
	action := data.NewAction()
	action.Twitter.Retweet = true
	timeline := TimelineConfig{}
	timeline.Action = action
	timeline.ScreenNames = []string{"foo"}
	timelines := []TimelineConfig{timeline}
	c.SetTwitterTimelines(timelines)
	ts := c.GetTwitterTimelines()
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
	action := data.NewAction()
	action.Twitter.Retweet = true
	favorite := FavoriteConfig{}
	favorite.Action = action
	favorite.ScreenNames = []string{"foo"}
	favorites := []FavoriteConfig{favorite}
	c.SetTwitterFavorites(favorites)
	fs := c.GetTwitterFavorites()
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
	action := data.NewAction()
	action.Twitter.Retweet = true
	search := SearchConfig{}
	search.Action = action
	search.Queries = []string{"foo"}
	searches := []SearchConfig{search}
	c.SetTwitterSearches(searches)
	ss := c.GetTwitterSearches()
	if !reflect.DeepEqual(searches, ss) {
		t.Fatalf("%v is not set properly", searches)
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
	notification := Notification{
		Place: PlaceNotification{
			Users: []string{"foo"},
		},
	}
	c.SetTwitterNotification(notification)
	n := c.GetTwitterNotification()
	if !reflect.DeepEqual(notification, n) {
		t.Fatalf("%v is not set properly", notification)
	}
}

func TestFileConfigSlackMessages(t *testing.T) {
	c, err := NewFileConfig("")
	if err != nil {
		t.Fatal(err)
	}
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
	if !reflect.DeepEqual(msgs, ms) {
		t.Fatalf("%v expected but %v found", msgs, ms)
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
	interaction := InteractionConfig{}
	interaction.Users = []string{"foo"}
	c.SetTwitterInteraction(interaction)
	i := c.GetTwitterInteraction()
	if !reflect.DeepEqual(interaction, i) {
		t.Fatalf("%v is not set properly", interaction)
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
	c.SetTwitterDuration(duration)
	dur := c.GetTwitterDuration()
	if duration != dur {
		t.Fatalf("%v is not set properly", duration)
	}
}
