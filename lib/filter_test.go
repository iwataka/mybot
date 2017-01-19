package mybot

import (
	"testing"

	"github.com/iwataka/anaconda"
)

type EmptyVisionMatcher struct{}

func (m *EmptyVisionMatcher) MatchImages(urls []string, c *VisionCondition) ([]string, []bool, error) {
	results := make([]string, len(urls), len(urls))
	matches := make([]bool, len(urls), len(urls))
	for i, _ := range urls {
		matches[i] = true
	}
	return results, matches, nil
}

func (m *EmptyVisionMatcher) Enabled() bool {
	return true
}

var visionMatcher = &EmptyVisionMatcher{}

type EmptyLanguageMatcher struct{}

func (m *EmptyLanguageMatcher) MatchText(text string, c *LanguageCondition) (string, bool, error) {
	return "", true, nil
}

func (m *EmptyLanguageMatcher) Enabled() bool {
	return true
}

var languageMatcher = &EmptyLanguageMatcher{}

var cache *Cache

func init() {
	var err error
	cache, err = NewCache("")
	if err != nil {
		panic(err)
	}
}

func TestCheckPatternsMatched(t *testing.T) {
	tweet := anaconda.Tweet{
		Text: "foo is bar",
	}
	conf := &TweetFilterConfig{
		Patterns: []string{"foo"},
	}
	match, err := conf.check(tweet, visionMatcher, languageMatcher, cache)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("%v is expected but %v found", true, match)
	}
}

func TestCheckPatternsUnmatched(t *testing.T) {
	tweet := anaconda.Tweet{
		Text: "fizz buzz",
	}
	conf := &TweetFilterConfig{
		Patterns: []string{"foo"},
	}
	match, err := conf.check(tweet, visionMatcher, languageMatcher, cache)
	if err != nil {
		t.Fatal(err)
	}
	if match {
		t.Fatalf("%v is expected but %v found", false, match)
	}
}

func TestCheckFavoriteThresholdExceeded(t *testing.T) {
	tweet := anaconda.Tweet{
		FavoriteCount: 100,
	}
	threshold := 80
	conf := &TweetFilterConfig{
		FavoriteThreshold: &threshold,
	}
	match, err := conf.check(tweet, visionMatcher, languageMatcher, cache)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("%v is expected but %v found", true, match)
	}
}

func TestCheckFavoriteThresholdNotExceeded(t *testing.T) {
	tweet := anaconda.Tweet{
		FavoriteCount: 100,
	}
	threshold := 120
	conf := &TweetFilterConfig{
		FavoriteThreshold: &threshold,
	}
	match, err := conf.check(tweet, visionMatcher, languageMatcher, cache)
	if err != nil {
		t.Fatal(err)
	}
	if match {
		t.Fatalf("%v is expected but %v found", false, match)
	}
}

func TestCheckRetweetedThresholdExceeded(t *testing.T) {
	tweet := anaconda.Tweet{
		RetweetCount: 100,
	}
	threshold := 80
	conf := &TweetFilterConfig{
		RetweetedThreshold: &threshold,
	}
	match, err := conf.check(tweet, visionMatcher, languageMatcher, cache)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("%v is expected but %v found", true, match)
	}
}

func TestCheckRetweetedThresholdNotExceeded(t *testing.T) {
	tweet := anaconda.Tweet{
		RetweetCount: 100,
	}
	threshold := 120
	conf := &TweetFilterConfig{
		RetweetedThreshold: &threshold,
	}
	match, err := conf.check(tweet, visionMatcher, languageMatcher, cache)
	if err != nil {
		t.Fatal(err)
	}
	if match {
		t.Fatalf("%v is expected but %v found", false, match)
	}
}

func TestCheckVisionMatched(t *testing.T) {
	tweet := anaconda.Tweet{
		Entities: anaconda.Entities{
			Media: []anaconda.EntityMedia{
				anaconda.EntityMedia{},
			},
		},
	}
	conf := &TweetFilterConfig{
		Vision: &VisionCondition{
			Label: []string{"foo"},
		},
	}
	match, err := conf.check(tweet, visionMatcher, languageMatcher, cache)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatalf("%v is expected but %v found", true, match)
	}
}

func TestCheckVisionUnmatched(t *testing.T) {
	tweet := anaconda.Tweet{
		Entities: anaconda.Entities{
			Media: []anaconda.EntityMedia{},
		},
	}
	conf := &TweetFilterConfig{
		Vision: &VisionCondition{
			Label: []string{"foo"},
		},
	}
	match, err := conf.check(tweet, visionMatcher, languageMatcher, cache)
	if err != nil {
		t.Fatal(err)
	}
	if match {
		t.Fatalf("%v is expected but %v found", false, match)
	}
}
