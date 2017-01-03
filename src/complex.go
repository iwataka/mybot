package mybot

import (
	"fmt"
	"regexp"

	"github.com/iwataka/anaconda"
)

// TweetFilterConfig is a configuration to filter out tweets
type TweetFilterConfig struct {
	Patterns           []string         `toml:"patterns,omitempty"`
	URLPatterns        []string         `toml:"url_patterns,omitempty"`
	HasMedia           *bool            `toml:"has_media"`
	HasURL             *bool            `toml:"has_url"`
	Retweeted          *bool            `toml:"retweeted"`
	FavoriteThreshold  int              `toml:"favorite_threshold"`
	RetweetedThreshold int              `toml:"retweeted_threshold"`
	Lang               string           `toml:"lang,omitempty"`
	Vision             *VisionCondition `toml:"vision"`
	VisionAPI          *VisionAPI
}

func (c *TweetFilterConfig) check(t anaconda.Tweet) (bool, error) {
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
	if c.FavoriteThreshold > 0 && c.FavoriteThreshold > t.FavoriteCount {
		return false, nil
	}
	if c.RetweetedThreshold > 0 && c.RetweetedThreshold > t.RetweetCount {
		return false, nil
	}
	if len(c.Lang) != 0 && c.Lang != t.Lang {
		return false, nil
	}
	if c.Vision != nil && c.VisionAPI != nil && c.VisionAPI.api != nil {
		urls := make([]string, len(t.Entities.Media))
		for i, m := range t.Entities.Media {
			urls[i] = m.Media_url
		}
		match := false
		var err error
		if len(urls) != 0 {
			match, err = c.VisionAPI.MatchImages(urls, c.Vision)
			if err != nil {
				return false, err
			}
			c.VisionAPI.cache.ImageSource = fmt.Sprintf("https://twitter.com/%s/status/%s", t.User.IdStr, t.IdStr)
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

func (c *TweetFilterConfig) shouldRepeat() bool {
	return c.RetweetedThreshold > 0 || c.FavoriteThreshold > 0
}