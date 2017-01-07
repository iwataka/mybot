package mybot

import (
	"testing"
)

func TestNewMybotConfig(t *testing.T) {
	c, err := NewMybotConfig("config.template.toml")
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
	if a.Action.Retweet != true {
		t.Fatalf("%v expected but %v found", true, a.Action.Retweet)
	}
	s := c.Twitter.Searches[0]
	if s.Queries[0] != "foo" {
		t.Fatalf("%s expected but %s found", "foo", s.Queries[0])
	}
	if s.Queries[1] != "bar" {
		t.Fatalf("%s expected but %s found", "bar", s.Queries[1])
	}
	if s.Filter.RetweetedThreshold != 100 {
		t.Fatalf("%d expected but %d found", 100, s.Filter.RetweetedThreshold)
	}
	if s.Action.Retweet != true {
		t.Fatalf("%v expected but %v found", true, s.Action.Retweet)
	}
	n := c.Twitter.Notification
	if n.Place.AllowSelf != true {
		t.Fatalf("%v expected but %v found", true, n.Place.AllowSelf)
	}
	if n.Place.Users[0] != "foo" {
		t.Fatalf("%s expected but %s found", "foo", n.Place.Users[0])
	}
}

func TestNewMybotConfigWhenNotExist(t *testing.T) {
	_, err := NewMybotConfig("config_not_exist.toml")
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
