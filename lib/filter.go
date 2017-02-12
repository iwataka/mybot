package mybot

import (
	"fmt"
	"regexp"
	"time"

	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/models"
)

// TweetFilter is a configuration to filter out tweets
type TweetFilter struct {
	models.TweetFilterProperties
	Patterns    []string           `json:"patterns,omitempty" toml:"patterns,omitempty"`
	URLPatterns []string           `json:"url_patterns,omitempty" toml:"url_patterns,omitempty"`
	Vision      *VisionCondition   `json:"vision,omitempty" toml:"vision,omitempty"`
	Language    *LanguageCondition `json:"language,omitempty" toml:"language,omitempty"`
}

func NewTweetFilter() *TweetFilter {
	return &TweetFilter{
		Vision:   NewVisionCondition(),
		Language: &LanguageCondition{},
	}
}

func (f *TweetFilter) Validate() error {
	flag := (f.Vision != nil && !f.Vision.isEmpty()) &&
		(f.RetweetedThreshold != nil || f.FavoriteThreshold != nil)
	if flag {
		return fmt.Errorf("%v use both of Vision API and retweeted/favorite threshold", f)
	}
	return nil
}

func (c *TweetFilter) check(t anaconda.Tweet, v VisionMatcher, l LanguageMatcher, cache Cache) (bool, error) {
	for _, p := range c.Patterns {
		match, err := regexp.MatchString(p, t.Text)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	for _, url := range c.URLPatterns {
		match := false
		var err error
		for _, u := range t.Entities.Urls {
			match, err = regexp.MatchString(url, u.Display_url)
			if err != nil {
				return false, err
			}
			if match {
				break
			}
		}
		if !match {
			return false, nil
		}
	}
	if c.HasMedia != nil && *c.HasMedia != (len(t.Entities.Media) != 0) {
		return false, nil
	}
	if c.HasURL != nil && *c.HasURL != (len(t.Entities.Urls) != 0) {
		return false, nil
	}
	if c.Retweeted != nil && *c.Retweeted != (t.RetweetedStatus != nil) {
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

	// If the Vision condition is empty or Vision API is not available,
	// skip this check. Otherwise if there is at least one media to satisfy
	// condition, the tweet will pass this check.
	if c.Vision != nil && !c.Vision.isEmpty() && v != nil && v.Enabled() {
		urls := make([]string, len(t.Entities.Media))
		for i, m := range t.Entities.Media {
			urls[i] = m.Media_url
		}

		results, matches, err := v.MatchImages(urls, c.Vision)
		if err != nil {
			return false, err
		}

		// Cache the results of matching images
		for i, result := range results {
			// Empty result means no analysis occurred
			if len(result) == 0 {
				continue
			}

			tweetSrc := TwitterStatusURL(t)
			imgCache := ImageCacheData{}
			imgCache.URL = urls[i]
			imgCache.Src = tweetSrc
			imgCache.AnalysisResult = result
			imgCache.AnalysisDate = time.Now().Format(time.RubyDate)
			err := cache.SetImage(imgCache)
			if err != nil {
				return false, err
			}
		}

		match := false
		for _, m := range matches {
			match = match || m
		}
		if !match {
			return false, nil
		}
	}

	// If the Language condition is empty or Language API is not available,
	// skip this check.
	if c.Language != nil && !c.Language.isEmpty() && l != nil && l.Enabled() {
		text := t.Text
		_, match, err := l.MatchText(text, c.Language)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}

	return true, nil
}

func (c *TweetFilter) shouldRepeat() bool {
	return c.RetweetedThreshold != nil || c.FavoriteThreshold != nil
}
