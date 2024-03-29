package core

import (
	"fmt"
	"regexp"
	"time"

	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/utils"
	"github.com/slack-go/slack"
)

// Filter is a configuration to filter out tweets
type Filter struct {
	models.FilterProperties `yaml:",inline"`
	Patterns                []string                 `json:"patterns,omitempty" toml:"patterns,omitempty" bson:"patterns,omitempty" yaml:"patterns,omitempty"`
	URLPatterns             []string                 `json:"url_patterns,omitempty" toml:"url_patterns,omitempty" bson:"url_patterns,omitempty" yaml:"url_patterns,omitempty"`
	Vision                  models.VisionCondition   `json:"vision,omitempty" toml:"vision,omitempty" bson:"vision,omitempty" yaml:"vision,omitempty"`
	Language                models.LanguageCondition `json:"language,omitempty" toml:"language,omitempty" bson:"language,omitempty" yaml:"language,omitempty"`
}

func NewFilter() Filter {
	return Filter{
		Patterns:    []string{},
		URLPatterns: []string{},
		Vision:      models.NewVisionCondition(),
		Language:    models.LanguageCondition{},
	}
}

func (f Filter) Validate() error {
	flag := !f.Vision.IsEmpty() && (f.RetweetedThreshold != nil || f.FavoriteThreshold != nil)
	if flag {
		return fmt.Errorf("Do not use both Vision API and retweeted/favorite threshold: %#v", f)
	}
	return nil
}

func (c Filter) CheckTweet(
	t anaconda.Tweet,
	v VisionMatcher,
	l LanguageMatcher,
	cache data.Cache,
) (bool, error) {
	match, err := c.checkTweetWithoutWebAPI(t)
	if err != nil {
		return false, utils.WithStack(err)
	}
	if !match {
		return false, nil
	}

	// If the Vision condition is empty or Vision API is not available,
	// skip this check. Otherwise if there is at least one media to satisfy
	// condition, the tweet will pass this check.
	if v != nil && v.Enabled() {
		match, err := c.matchTweetImages(t, v, cache)
		if err != nil {
			return false, utils.WithStack(err)
		}
		if !match {
			return false, nil
		}
	}

	// If the Language condition is empty or Language API is not available,
	// skip this check.
	if l != nil && l.Enabled() {
		_, match, err := l.MatchText(t.Text, c.Language)
		if err != nil {
			return false, utils.WithStack(err)
		}
		if !match {
			return false, nil
		}
	}

	return true, nil
}

func (c Filter) checkTweetWithoutWebAPI(t anaconda.Tweet) (bool, error) {
	for _, p := range c.Patterns {
		match, err := regexp.MatchString(p, t.Text)
		if err != nil {
			return false, utils.WithStack(err)
		}
		if !match {
			return false, nil
		}
	}

	for _, url := range c.URLPatterns {
		match, err := checkTweetMatchesURLPattern(t, url)
		if err != nil {
			return false, utils.WithStack(err)
		}
		if !match {
			return false, nil
		}
	}

	if c.HasMedia != nil && *c.HasMedia != (len(t.Entities.Media) != 0) {
		return false, nil
	}
	if c.FavoriteThreshold != nil && *c.FavoriteThreshold > t.FavoriteCount {
		return false, nil
	}
	if c.RetweetedThreshold != nil && *c.RetweetedThreshold > t.RetweetCount {
		return false, nil
	}
	if len(c.Lang) != 0 && c.Lang != t.Lang {
		return false, nil
	}

	return true, nil
}

func checkTweetMatchesURLPattern(t anaconda.Tweet, pattern string) (bool, error) {
	for _, u := range t.Entities.Urls {
		match, err := regexp.MatchString(pattern, u.Display_url)
		if err != nil {
			return false, utils.WithStack(err)
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

func (c Filter) CheckSlackMsg(
	ev *slack.MessageEvent,
	v VisionMatcher,
	l LanguageMatcher,
	cache data.Cache,
) (bool, error) {
	for _, p := range c.Patterns {
		match, err := regexp.MatchString(p, ev.Text)
		if err != nil {
			return false, utils.WithStack(err)
		}
		if !match {
			return false, nil
		}
	}

	if c.HasMedia != nil {
		hasMedia, err := checkSlackMsgHasMedia(ev)
		if err != nil {
			return false, utils.WithStack(err)
		}
		if !hasMedia {
			return false, nil
		}
	}

	// If the Vision condition is empty or Vision API is not available,
	// skip this check. Otherwise if there is at least one media to satisfy
	// condition, the tweet will pass this check.
	if v != nil && v.Enabled() {
		match, err := c.matchSlackImages(ev.Attachments, v, cache)
		if err != nil {
			return false, utils.WithStack(err)
		}
		if !match {
			return false, nil
		}
	}

	// If the Language condition is empty or Language API is not available,
	// skip this check.
	if l != nil && l.Enabled() {
		_, match, err := l.MatchText(ev.Text, c.Language)
		if err != nil {
			return false, utils.WithStack(err)
		}
		if !match {
			return false, nil
		}
	}

	return true, nil
}

func checkSlackMsgHasMedia(ev *slack.MessageEvent) (bool, error) {
	for _, a := range ev.Attachments {
		if a.ImageURL != "" {
			return true, nil
		}
	}
	return false, nil
}

func (c Filter) matchTweetImages(t anaconda.Tweet, v VisionMatcher, cache data.Cache) (bool, error) {
	urls := make([]string, len(t.Entities.Media))
	for i, m := range t.Entities.Media {
		urls[i] = m.Media_url_https
	}
	return c.matchImageURLs(TwitterStatusURL(t), urls, v, cache)
}

func (c Filter) matchSlackImages(atts []slack.Attachment, v VisionMatcher, cache data.Cache) (bool, error) {
	urls := []string{}
	for _, a := range atts {
		if a.ImageURL == "" {
			continue
		}
		urls = append(urls, a.ImageURL)
	}
	return c.matchImageURLs("", urls, v, cache)
}

func (c Filter) matchImageURLs(src string, urls []string, v VisionMatcher, cache data.Cache) (bool, error) {
	if c.Vision.IsEmpty() {
		return true, nil
	}

	results, matches, err := v.MatchImages(urls, c.Vision, cache.GetLatestImages(-1))
	if err != nil {
		return false, utils.WithStack(err)
	}

	// Cache the results of matching images
	for i, result := range results {
		// Empty result means no analysis occurred
		if len(result) == 0 {
			continue
		}

		tweetSrc := src
		imgCache := models.ImageCacheData{}
		imgCache.URL = urls[i]
		imgCache.Src = tweetSrc
		imgCache.AnalysisResult = result
		imgCache.AnalysisDate = time.Now().Format(time.RubyDate)
		cache.SetImage(imgCache)
	}

	match := false
	for _, m := range matches {
		match = match || m
	}
	if !match {
		return false, nil
	}
	return true, nil
}

func (c Filter) ShouldRepeat() bool {
	return c.RetweetedThreshold != nil || c.FavoriteThreshold != nil
}
