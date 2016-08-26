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
	if *a.Filter.HasUrl != true {
		t.Fatalf("%v expected but %v found", true, *a.Filter.HasUrl)
	}
	if *a.Filter.Retweeted != false {
		t.Fatalf("%v expected but %v found", false, *a.Filter.Retweeted)
	}
	if *a.Filter.RetweetedThreshold != 100 {
		t.Fatalf("%d expected but %d found", 100, *a.Filter.RetweetedThreshold)
	}
	if *a.Filter.Lang != "en" {
		t.Fatalf("%s expected but %s found", "en", *a.Filter.Lang)
	}
	if a.Filter.Vision.Label[0] != "cartoon|clip art|artwork" {
		t.Fatalf("%s expected but %s found", "cartoon|clip art|artwork", a.Filter.Vision.Label[0])
	}
	if a.Action.Retweet == false {
		t.Fatalf("%v expected but %v found", true, a.Action.Retweet)
	}
	n := c.Twitter.Notification
	if n.Place.AllowSelf != true {
		t.Fatalf("%v expected but %v found", true, n.Place.AllowSelf)
	}
	if n.Place.Users[0] != "foo" {
		t.Fatalf("%s expected but %s found", "foo", n.Place.Users[0])
	}
}
