package main

import (
	"regexp"

	"github.com/iwataka/anaconda"
)

type TweetFilterConfig struct {
	Patterns           []string
	UrlPatterns        []string `toml:"url_patterns"`
	HasMedia           *bool    `toml:"has_media"`
	HasUrl             *bool    `toml:"has_url"`
	Retweeted          *bool
	FavoriteThreshold  *int `toml:"favorite_threshold"`
	RetweetedThreshold *int `toml:"retweeted_threshold"`
	Vision             *VisionCondition
}

func (c *TweetFilterConfig) GetChecker(a *VisionAPI) TweetChecker {
	return func(t anaconda.Tweet) (bool, error) {
		for _, p := range c.Patterns {
			match, err := regexp.MatchString(p, t.Text)
			if err != nil {
				return false, err
			}
			if !match {
				return false, nil
			}
		}
		for _, url := range c.UrlPatterns {
			match := false
			for _, u := range t.Entities.Urls {
				match, err := regexp.MatchString(url, u.Display_url)
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
		if c.HasUrl != nil && *c.HasUrl != (len(t.Entities.Urls) != 0) {
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
		if c.Vision != nil && a != nil && a.Images != nil {
			urls := make([]string, len(t.Entities.Media))
			for i, m := range t.Entities.Media {
				urls[i] = m.Media_url
			}
			match, err := a.MatchImage(urls, c.Vision)
			if err != nil {
				return false, err
			}
			if !match {
				return false, nil
			}
		}
		return true, nil
	}
}
