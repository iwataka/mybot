package mybot

import (
	"strings"
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
	tweet := anaconda.Tweet{}
	tweet.Text = "some texts"
	m1 := anaconda.EntityMedia{}
	m1.Media_url = "foo"
	m2 := anaconda.EntityMedia{}
	m2.Media_url = "bar"
	tweet.Entities.Media = []anaconda.EntityMedia{m1, m2}

	api := NewSlackAPI("")
	_, params := api.convertFromTweet(tweet)
	if len(params.Attachments) != 2 {
		t.Fatalf("%d expected but %d found", 2, len(params.Attachments))
	}
	att1 := params.Attachments[0]
	if !strings.HasSuffix(att1.Text, m1.Media_url) {
		t.Fatalf(`"%s" should have the suffix "%s"`, att1.Text, m1.Media_url)
	}
	if !strings.HasPrefix(att1.Text, tweet.Text) {
		t.Fatalf(`"%s" should have the prefix "%s"`, att1.Text, tweet.Text)
	}
	att2 := params.Attachments[1]
	if !strings.HasSuffix(att2.Text, m2.Media_url) {
		t.Fatalf(`"%s" should have the suffix "%s"`, att2.Text, m2.Media_url)
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
