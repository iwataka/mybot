package main

import (
	"net/http"
	"regexp"

	"github.com/ChimeraCoder/anaconda"
)

type TweetCheckConfig struct {
	Name     string
	Patterns []string
	Opts     map[string]bool
	Image    *struct{ Descriptions []string }
}

func (c *TweetCheckConfig) GetChecker(a *VisionAPI) TweetChecker {
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
		for key, val := range c.Opts {
			if key == "hasMedia" {
				if val != (len(t.Entities.Media) != 0) {
					return false, nil
				}
			} else if key == "hasUrl" {
				if val != (len(t.Entities.Urls) != 0) {
					return false, nil
				}
			} else if key == "retweeted" {
				if val != t.Retweeted {
					return false, nil
				}
			}
		}
		if c.Image != nil && a != nil {
			imgData := make([][]byte, len(t.Entities.Media))
			for i, m := range t.Entities.Media {
				resp, err := http.Get(m.Media_url)
				_, err = resp.Body.Read(imgData[i])
				if err != nil {
					return false, err
				}
			}
			match, err := a.MatchImageDescription(imgData, c.Image.Descriptions)
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
