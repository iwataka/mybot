package mybot

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	c, err := NewConfig("test_assets/config.template.toml")
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
}

func TestNewConfigWhenNotExist(t *testing.T) {
	_, err := NewConfig("config_not_exist.toml")
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
			Retweet:     true,
			Collections: []string{"foo"},
		},
		Slack: &SlackAction{
			Channels: []string{"bar"},
		},
	}

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
			Retweet:     true,
			Collections: []string{"twitter"},
		},
		Slack: &SlackAction{
			Channels: []string{"slack"},
		},
	}
	a2 := &TweetAction{
		Twitter: &TwitterAction{
			Favorite:    true,
			Collections: []string{"facebook"},
		},
		Slack: &SlackAction{
			Channels: []string{"mattermost"},
		},
	}
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
