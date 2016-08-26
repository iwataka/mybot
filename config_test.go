package main

import "testing"

func TestNewMybotConfig(t *testing.T) {
	c, err := NewMybotConfig("config.template.toml")
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	p := c.GitHub.Projects[0]
	if p.User != "golang" {
		t.Fatalf("%s expected but %s found", "golang", p.User)
	}
	if p.Repo != "go" {
		t.Fatalf("%s expected but %s found", "go", p.Repo)
	}
	if c.GitHub.Duration != "12h" {
		t.Fatalf("%s expected but %s found", "30m", c.GitHub.Duration)
	}
	a := c.Twitter.Accounts[0]
	if *a.Name != "golang" {
		t.Fatalf("%s expected but %s found", "golang", a.Name)
	}
	f := a.Filter
	if f.Patterns[0] != "is released!" {
		t.Fatalf("%s expected but %s found", "is released!", f.Patterns[0])
	}
	if *f.HasUrl != true {
		t.Fatalf("%v expected but %v found", true, *f.HasUrl)
	}
	if *f.Retweeted != false {
		t.Fatalf("%v expected but %v found", false, *f.Retweeted)
	}
	if *f.Lang != "en" {
		t.Fatalf("%s expected but %s found", "en", *f.Lang)
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
	if *s.Filter.RetweetedThreshold != 100 {
		t.Fatalf("%d expected but %d found", 100, *s.Filter.RetweetedThreshold)
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
