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
	if a.Name != "golang" {
		t.Fatalf("%s expected but %s found", "golang", a.Name)
	}
	if a.Filter.Patterns[0] != "is released!" {
		t.Fatalf("%s expected but %s found", "is released!", a.Filter.Patterns[0])
	}
	if a.Filter.Opts["retweeted"] != false {
		t.Fatalf("%v expected but %v found", false, a.Filter.Opts["retweeted"])
	}
	if a.Filter.Vision.Label[0] != "cartoon|clip art|artwork" {
		t.Fatalf("%s expected but %s found", "cartoon|clip art|artwork", a.Filter.Vision.Label[0])
	}
	if a.Actions[0] != "retweet" {
		t.Fatalf("%s expected but %s found", "retweet", a.Actions[0])
	}
	n := c.Twitter.Notification
	if n.Place.AllowSelf != true {
		t.Fatalf("%v expected but %v found", true, n.Place.AllowSelf)
	}
	if n.Place.Users[0] != "foo" {
		t.Fatalf("%s expected but %s found", "foo", n.Place.Users[0])
	}
}
