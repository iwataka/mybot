package main

import (
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

var twitterApi *anaconda.TwitterApi

var twitterTalkUsers = [0]string{}

func retweet(screenName string, trimUser bool, checker func(anaconda.Tweet) bool) error {
	v := url.Values{}
	v.Set("screen_name", screenName)
	tweets, err := twitterApi.GetUserTimeline(v)
	if err != nil {
		return err
	}
	latestId, exists := cache.LatestTweetId[screenName]
	finds := false
	updates := false
	for i := len(tweets) - 1; i >= 0; i-- {
		tweet := tweets[i]
		if checker(tweet) {
			if exists && latestId == tweet.Id {
				finds = true
			} else {
				updates = true
				cache.LatestTweetId[screenName] = tweet.Id
				if finds {
					_, err := twitterApi.Retweet(tweet.Id, trimUser)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	if !exists && updates {
		_, err := twitterApi.Retweet(cache.LatestTweetId[screenName], trimUser)
		if err != nil {
			return err
		}
	}
	return nil
}

func talk() error {
	dms, err := twitterApi.GetDirectMessages(nil)
	if err != nil {
		return err
	}
	self, err := twitterApi.GetSelf(nil)
	if err != nil {
		return err
	}
	userToDM := make(map[string]anaconda.DirectMessage)
	for _, dm := range dms {
		sender := dm.SenderScreenName
		allowed := false
		if sender == self.ScreenName {
			allowed = true
		} else {
			for _, u := range twitterTalkUsers {
				if sender == u {
					allowed = true
					break
				}
			}
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
