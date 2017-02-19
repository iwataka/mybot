package mybot

import (
	"testing"

	"github.com/iwataka/anaconda"
	"github.com/nlopes/slack"
)

func init() {
	visionMatcher = &EmptyVisionMatcher{}
	languageMatcher = &EmptyLanguageMatcher{}

	var err error
	cache, err = NewFileCache("")
	if err != nil {
		panic(err)
	}
}

var (
	visionMatcher   *EmptyVisionMatcher
	languageMatcher *EmptyLanguageMatcher
	cache           Cache
)

type (
	EmptyVisionMatcher   struct{}
	EmptyLanguageMatcher struct{}
)

func (m *EmptyVisionMatcher) MatchImages(urls []string, c *VisionCondition) ([]string, []bool, error) {
	results := make([]string, len(urls), len(urls))
	matches := make([]bool, len(urls), len(urls))
	for i := range urls {
		matches[i] = true
	}
	return results, matches, nil
}

func (m *EmptyVisionMatcher) Enabled() bool {
	return true
}

func (m *EmptyLanguageMatcher) MatchText(text string, c *LanguageCondition) (string, bool, error) {
	return "", true, nil
}

func (m *EmptyLanguageMatcher) Enabled() bool {
	return true
}

func TestCheckTweetPatternsMatched(t *testing.T) {
	tweet := anaconda.Tweet{
		Text: "foo is bar",
	}
	conf := &Filter{
		Patterns: []string{"foo"},
	}
	match, err := conf.CheckTweet(tweet, visionMatcher, languageMatcher, cache)
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
	match, err := conf.CheckSlackMsg(ev, visionMatcher, languageMatcher, cache)
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
	match, err := conf.CheckTweet(tweet, visionMatcher, languageMatcher, cache)
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
	match, err := conf.CheckSlackMsg(ev, visionMatcher, languageMatcher, cache)
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
	match, err := conf.CheckTweet(tweet, visionMatcher, languageMatcher, cache)
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
	match, err := conf.CheckTweet(tweet, visionMatcher, languageMatcher, cache)
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
	match, err := conf.CheckTweet(tweet, visionMatcher, languageMatcher, cache)
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
	match, err := conf.CheckTweet(tweet, visionMatcher, languageMatcher, cache)
	if err != nil {
		t.Fatal(err)
	}
	if match {
		t.Fatalf("%v expected but %v found", false, match)
	}
}

func TestCheckTweetVisionMatched(t *testing.T) {
	tweet := anaconda.Tweet{
		Entities: anaconda.Entities{
			Media: []anaconda.EntityMedia{
				{},
			},
		},
	}
	conf := &Filter{
		Vision: &VisionCondition{
			Label: []string{"foo"},
		},
	}
	match, err := conf.CheckTweet(tweet, visionMatcher, languageMatcher, cache)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("%v expected but %v found", true, match)
	}
}

func TestCheckSlackMsgVisionMatched(t *testing.T) {
	att := slack.Attachment{
		ImageURL: "url",
	}
	conf := &Filter{
		Vision: &VisionCondition{
			Label: []string{"foo"},
		},
	}
	ev := &slack.MessageEvent{}
	ev.Attachments = []slack.Attachment{att}
	match, err := conf.CheckSlackMsg(ev, visionMatcher, languageMatcher, cache)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("%v expected but %v found", true, match)
	}
}

func TestCheckTweetVisionUnmatched(t *testing.T) {
	conf := &Filter{
		Vision: &VisionCondition{
			Label: []string{"foo"},
		},
	}
	ev := &slack.MessageEvent{}
	ev.Attachments = []slack.Attachment{}
	match, err := conf.CheckSlackMsg(ev, visionMatcher, languageMatcher, cache)
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
