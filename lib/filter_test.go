package mybot

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/mocks"
	"github.com/iwataka/mybot/models"
	"github.com/nlopes/slack"
)

func TestCheckTweetPatternsMatched(t *testing.T) {
	tweet := anaconda.Tweet{
		Text: "foo is bar",
	}
	conf := &Filter{
		Patterns: []string{"foo"},
	}
	cache := NewTestFileCache("", t)
	match, err := conf.CheckTweet(tweet, nil, nil, cache)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("%v expected but %v found", true, match)
	}
}

func TestCheckSlackMsgPatternsMatched(t *testing.T) {
	conf := &Filter{
		Patterns: []string{"foo"},
	}
	ev := &slack.MessageEvent{}
	ev.Attachments = []slack.Attachment{}
	ev.Text = "foo is bar"
	cache := NewTestFileCache("", t)
	match, err := conf.CheckSlackMsg(ev, nil, nil, cache)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("%v expected but %v found", true, match)
	}
}

func TestCheckTweetPatternsUnmatched(t *testing.T) {
	tweet := anaconda.Tweet{
		Text: "fizz buzz",
	}
	conf := &Filter{
		Patterns: []string{"foo"},
	}
	cache := NewTestFileCache("", t)
	match, err := conf.CheckTweet(tweet, nil, nil, cache)
	if err != nil {
		t.Fatal(err)
	}
	if match {
		t.Fatalf("%v expected but %v found", false, match)
	}
}

func TestCheckSlackMsgPatternsUnmatched(t *testing.T) {
	conf := &Filter{
		Patterns: []string{"foo"},
	}
	ev := &slack.MessageEvent{}
	ev.Attachments = []slack.Attachment{}
	ev.Text = "fizz buzz"
	cache := NewTestFileCache("", t)
	match, err := conf.CheckSlackMsg(ev, nil, nil, cache)
	if err != nil {
		t.Fatal(err)
	}
	if match {
		t.Fatalf("%v expected but %v found", false, match)
	}
}

func TestCheckTweetFavoriteThresholdExceeded(t *testing.T) {
	tweet := anaconda.Tweet{
		FavoriteCount: 100,
	}
	threshold := 80
	conf := NewFilter()
	conf.FavoriteThreshold = &threshold
	cache := NewTestFileCache("", t)
	match, err := conf.CheckTweet(tweet, nil, nil, cache)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("%v expected but %v found", true, match)
	}
}

func TestCheckTweetFavoriteThresholdNotExceeded(t *testing.T) {
	tweet := anaconda.Tweet{
		FavoriteCount: 100,
	}
	threshold := 120
	conf := NewFilter()
	conf.FavoriteThreshold = &threshold
	cache := NewTestFileCache("", t)
	match, err := conf.CheckTweet(tweet, nil, nil, cache)
	if err != nil {
		t.Fatal(err)
	}
	if match {
		t.Fatalf("%v expected but %v found", false, match)
	}
}

func TestCheckTweetRetweetedThresholdExceeded(t *testing.T) {
	tweet := anaconda.Tweet{
		RetweetCount: 100,
	}
	threshold := 80
	conf := NewFilter()
	conf.RetweetedThreshold = &threshold
	cache := NewTestFileCache("", t)
	match, err := conf.CheckTweet(tweet, nil, nil, cache)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("%v expected but %v found", true, match)
	}
}

func TestCheckTweetRetweetedThresholdNotExceeded(t *testing.T) {
	tweet := anaconda.Tweet{
		RetweetCount: 100,
	}
	threshold := 120
	conf := NewFilter()
	conf.RetweetedThreshold = &threshold
	cache := NewTestFileCache("", t)
	match, err := conf.CheckTweet(tweet, nil, nil, cache)
	if err != nil {
		t.Fatal(err)
	}
	if match {
		t.Fatalf("%v expected but %v found", false, match)
	}
}

func TestCheckTweetVisionMatched(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := mocks.NewMockVisionMatcher(ctrl)
	v.EXPECT().Enabled().Return(true)
	v.EXPECT().MatchImages(gomock.Any(), gomock.Any()).Return([]string{""}, []bool{true}, nil)

	tweet := anaconda.Tweet{
		Entities: anaconda.Entities{
			Media: []anaconda.EntityMedia{
				{},
			},
		},
	}
	conf := &Filter{
		Vision: models.VisionCondition{
			Label: []string{"foo"},
		},
	}

	cache := NewTestFileCache("", t)
	match, err := conf.CheckTweet(tweet, v, nil, cache)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("%v expected but %v found", true, match)
	}
}

func TestCheckSlackMsgVisionMatched(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := mocks.NewMockVisionMatcher(ctrl)
	v.EXPECT().Enabled().Return(true)
	v.EXPECT().MatchImages(gomock.Any(), gomock.Any()).Return([]string{""}, []bool{true}, nil)

	att := slack.Attachment{
		ImageURL: "url",
	}
	conf := &Filter{
		Vision: models.VisionCondition{
			Label: []string{"foo"},
		},
	}
	ev := &slack.MessageEvent{}
	ev.Attachments = []slack.Attachment{att}

	cache := NewTestFileCache("", t)
	match, err := conf.CheckSlackMsg(ev, v, nil, cache)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("%v expected but %v found", true, match)
	}
}

func TestCheckTweetVisionUnmatched(t *testing.T) {
	ctrl := gomock.NewController(t)
	v := mocks.NewMockVisionMatcher(ctrl)
	v.EXPECT().Enabled().Return(true)
	v.EXPECT().MatchImages(gomock.Any(), gomock.Any()).Return([]string{""}, []bool{false}, nil)

	conf := &Filter{
		Vision: models.VisionCondition{
			Label: []string{"foo"},
		},
	}
	ev := &slack.MessageEvent{}
	ev.Attachments = []slack.Attachment{}

	cache := NewTestFileCache("", t)
	match, err := conf.CheckSlackMsg(ev, v, nil, cache)
	if err != nil {
		t.Fatal(err)
	}
	if match {
		t.Fatalf("%v expected but %v found", false, match)
	}
}

func TestFilterValidate(t *testing.T) {
	threshold := 100
	f := NewFilter()
	f.FavoriteThreshold = &threshold
	f.Vision.Label = []string{"foo"}
	err := f.Validate()
	if err == nil {
		t.Fatalf("%v should be invalid but not", f)
	}
}
