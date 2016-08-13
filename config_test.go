package main

import "testing"

func TestNewMybotConfig(t *testing.T) {
	c, err := NewMybotConfig("fixtures/config.yml")
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
	if c.GitHub.Duration != "30m" {
		t.Fatalf("%s expected but %s found", "30m", c.GitHub.Duration)
	}
	a := c.Retweet.Accounts[0]
	if a.Name != "golang" {
		t.Fatalf("%s expected but %s found", "golang", a.Name)
	}
	if a.Patterns[0] != "is released!" {
		t.Fatalf("%s expected but %s found", "is released!", a.Patterns[0])
	}
	if a.Opts["retweeted"] != false {
		t.Fatalf("%v expected but %v found", false, a.Opts["retweeted"])
	}
	n := c.Retweet.Notification
	if n.Place.AllowSelf != true {
		t.Fatalf("%v expected but %v found", true, n.Place.AllowSelf)
	}
	if n.Place.Users[0] != "John" {
		t.Fatalf("%s expected but %s found", "John", n.Place.Users[0])
	}
}
