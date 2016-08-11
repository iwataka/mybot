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

func twitterCheckUser(user string, allowSelf bool, users []string) (bool, error) {
	if allowSelf {
		self, err := getTwitterSelf()
		if err != nil {
			return false, err
		}
		if user == self {
			return true, nil
		}
	}
	for _, u := range users {
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
			err = twitterPostInfo(t, config.Retweet)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func twitterPostInfo(t anaconda.Tweet, c *retweetConfig) error {
	if c.Notification.Place != nil && t.HasCoordinates() {
		msg := fmt.Sprintf("ID: %s\nCountry: %s\nCreatedAt: %s", t.IdStr, t.Place.Country, t.CreatedAt)
		allowSelf := c.Notification.Place.AllowSelf
		users := c.Notification.Place.Users
		return twitterPost(msg, allowSelf, users)
	}
	return nil
}

func twitterPost(msg string, allowSelf bool, users []string) error {
	for _, user := range users {
		_, err := twitterApi.PostDMToScreenName(msg, user)
		if err != nil {
			return err
		}
	}
	if allowSelf {
		self, err := getTwitterSelf()
		if err != nil {
			return err
		}
		_, err = twitterApi.PostDMToScreenName(msg, self)
		if err != nil {
			return err
		}
	}
	return nil
}

func twitterInteract() error {
	if config.Interaction == nil {
		return nil
	}
	dms, err := twitterApi.GetDirectMessages(nil)
	if err != nil {
		return err
	}
	senderToDM := make(map[string]anaconda.DirectMessage)
	for _, dm := range dms {
		sender := dm.SenderScreenName
		allowed, err := twitterCheckUser(sender, false, config.Interaction.Users)
		if err != nil {
			return err
		}
		if allowed {
			_, exists := senderToDM[sender]
			if !exists {
				senderToDM[sender] = dm
			}
		}
	}
	for sender, dm := range senderToDM {
		err := twitterResponse(sender, dm)
		if err != nil {
			return err
		}
	}
	return nil
}

func twitterResponse(sender string, dm anaconda.DirectMessage) error {
	latest, exists := cache.LatestDM[sender]
	if !exists || latest != dm.Id {
		res, err := twitterApi.PostDMToScreenName(dm.Text, sender)
		if err != nil {
			return err
		}
		cache.LatestDM[sender] = res.Id
	}
	return nil
}
