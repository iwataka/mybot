package mybot

import (
	"testing"

	"github.com/iwataka/anaconda"
)

func TestSlackActionAdd(t *testing.T) {
	a1 := &SlackAction{
		Channels: []string{"foo", "bar"},
	}
	a2 := &SlackAction{
		Channels: []string{"foo", "fizz"},
	}
	result := a1.Add(a2)
	if len(result.Channels) != 3 {
		t.Fatalf("%v expected but %v found", 3, len(result.Channels))
	}
}

func TestSlackActionSub(t *testing.T) {
	a1 := &SlackAction{
		Channels: []string{"foo", "bar"},
	}
	a2 := &SlackAction{
		Channels: []string{"foo", "fizz"},
	}
	result := a1.Sub(a2)
	if len(result.Channels) != 1 {
		t.Fatalf("%v expected but %v found", 1, len(result.Channels))
	}
}

func TestSlackConvertFromTweet(t *testing.T) {
	slack := NewSlackAPI("")
	tweet := anaconda.Tweet{
		IdStr: "1",
		User: anaconda.User{
			IdStr: "1",
		},
	}
	text, params := slack.convertFromTweet(tweet)
	if text != TwitterStatusURL(tweet) {
		t.Fatal("Text is invalid")
	}
	if !params.UnfurlLinks || !params.UnfurlMedia {
		t.Fatal("Should unfurl all kinds of things")
	}
}

func TestNewSlackAPI(t *testing.T) {
	api := NewSlackAPI("")
	if api.Enabled() {
		t.Fatalf("%v is expected to be disabled but not", api)
	}
}

func TestSlackActionIsEmpty(t *testing.T) {
	a := NewSlackAction()
	if !a.IsEmpty() {
		t.Fatalf("%v should be empty", a)
	}
	a.Channels = []string{"foo"}
	if a.IsEmpty() {
		t.Fatalf("%v should not be empty", a)
	}
}
