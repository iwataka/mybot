package mybot

import (
	"fmt"
	"regexp"
	"time"

	"github.com/iwataka/anaconda"
)

// TweetFilterConfig is a configuration to filter out tweets
type TweetFilterConfig struct {
	Patterns           []string           `toml:"patterns,omitempty"`
	URLPatterns        []string           `toml:"url_patterns,omitempty"`
	HasMedia           *bool              `toml:"has_media,omitempty"`
	HasURL             *bool              `toml:"has_url,omitempty"`
	Retweeted          *bool              `toml:"retweeted,omitempty"`
	FavoriteThreshold  int                `toml:"favorite_threshold"`
	RetweetedThreshold int                `toml:"retweeted_threshold"`
	Lang               string             `toml:"lang,omitempty"`
	Vision             *VisionCondition   `toml:"vision,omitempty"`
	Language           *LanguageCondition `toml:"language,omitempty"`
}

func (c *TweetFilterConfig) check(t anaconda.Tweet, v *VisionAPI, l *LanguageAPI) (bool, error) {
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

	if c.Vision != nil && v != nil && v.api != nil {
		urls := make([]string, len(t.Entities.Media))
		for i, m := range t.Entities.Media {
			urls[i] = m.Media_url
		}
		if len(urls) != 0 {
			results, matches, err := v.MatchImages(urls, c.Vision)
			if err != nil {
				return false, err
			}

			result := results[len(results)-1]
			// empty result means no Vision API analysis occurred.
			if len(result) != 0 {
				v.cache.ImageAnalysisDates =
					append(v.cache.ImageAnalysisDates, time.Now().Format(time.RubyDate))
				v.cache.ImageAnalysisResults =
					append(v.cache.ImageAnalysisResults, results[len(results)-1])
				srcFmt := "https://twitter.com/%s/status/%s"
				tweetSrc := fmt.Sprintf(srcFmt, t.User.IdStr, t.IdStr)
				v.cache.ImageSources = append(v.cache.ImageSources, tweetSrc)
				v.cache.ImageURLs = append(v.cache.ImageURLs, urls[len(urls)-1])
			}

			for _, m := range matches {
				if m {
					return true, nil
				}
			}
		}
	}

	if c.Language != nil && l != nil && l.api != nil {
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

func (c *TweetFilterConfig) shouldRepeat() bool {
	return c.RetweetedThreshold > 0 || c.FavoriteThreshold > 0
}
