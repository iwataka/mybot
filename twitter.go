package main

import (
	"fmt"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

var twitterApi *anaconda.TwitterApi

var twitterSelf string

func getTwitterSelf() (string, error) {
	if twitterSelf == "" {
		self, err := twitterApi.GetSelf(nil)
		if err != nil {
			return "", err
		}
		twitterSelf = self.ScreenName
	}
	return twitterSelf, nil
}

func twitterCheckUser(user string) (bool, error) {
	if config.UserGroup.IncludeSelf {
		self, err := getTwitterSelf()
		if err != nil {
			return false, err
		}
		if user == self {
			return true, nil
		}
	}
	for _, u := range config.UserGroup.Users {
		if user == u {
			return true, nil
		}
	}
	return false, nil
}

func twitterRetweet(name string, trimUser bool, check func(anaconda.Tweet) bool) error {
	v := url.Values{}
	v.Set("screen_name", name)
	latestId, exists := cache.LatestTweetId[name]
	if exists {
		v.Set("since_id", fmt.Sprintf("%d", latestId))
	}
	tweets, err := twitterApi.GetUserTimeline(v)
	if err != nil {
		return err
	}
	for i := len(tweets) - 1; i >= 0; i-- {
		t := tweets[i]
		if check(t) {
			cache.LatestTweetId[name] = t.Id
			_, err := twitterApi.Retweet(t.Id, trimUser)
			if err != nil {
				return err
			}
			err = twitterPostInfo(t)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func twitterPostInfo(t anaconda.Tweet) error {
	if config.Notification.Place && t.HasCoordinates() {
		msg := fmt.Sprintf("ID: %s\nCountry: %s\nCreatedAt: %s", t.IdStr, t.Place.Country, t.CreatedAt)
		return twitterPost(msg)
	}
	return nil
}

func twitterPost(msg string) error {
	for _, user := range config.UserGroup.Users {
		twitterApi.PostDMToScreenName(msg, user)
	}
	if config.UserGroup.IncludeSelf {
		self, err := getTwitterSelf()
		if err != nil {
			return err
		}
		twitterApi.PostDMToScreenName(msg, self)
	}
	return nil
}

func twitterTalk() error {
	if !config.Talk.Enabled {
		return nil
	}
	dms, err := twitterApi.GetDirectMessages(nil)
	if err != nil {
		return err
	}
	userToDM := make(map[string]anaconda.DirectMessage)
	for _, dm := range dms {
		sender := dm.SenderScreenName
		allowed, err := twitterCheckUser(sender)
		if err != nil {
			return err
		}
		if allowed {
			_, exists := userToDM[sender]
			if !exists {
				userToDM[sender] = dm
			}
		}
	}
	for user, dm := range userToDM {
		latest, exists := cache.LatestDM[user]
		if !exists || latest != dm.Id {
			res, err := twitterApi.PostDMToScreenName(dm.Text, user)
			if err != nil {
				return err
			}
			cache.LatestDM[user] = res.Id
		}
	}
	return nil
}
