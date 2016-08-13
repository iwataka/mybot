package main

import (
	"fmt"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

type TwitterAPI struct {
	*anaconda.TwitterApi
	self  *anaconda.User
	cache *MybotCache
}

type TwitterAuth struct {
	ConsumerKey       string `yaml:"consumerKey"`
	ConsumerSecret    string `yaml:"consumerSecret"`
	AccessToken       string `yaml:"accessToken"`
	AccessTokenSecret string `yaml:"accessTokenSecret"`
}

func NewTwitterAPI(a *TwitterAuth, c *MybotCache) *TwitterAPI {
	anaconda.SetConsumerKey(a.ConsumerKey)
	anaconda.SetConsumerSecret(a.ConsumerSecret)
	api := anaconda.NewTwitterApi(a.AccessToken, a.AccessTokenSecret)
	return &TwitterAPI{api, nil, c}
}

// GetSelfCache returns the user of this client
func (a *TwitterAPI) GetSelfCache() (anaconda.User, error) {
	if a.self == nil {
		self, err := a.GetSelf(nil)
		if err != nil {
			return anaconda.User{}, err
		}
		a.self = &self
	}
	return *a.self, nil
}

func (a *TwitterAPI) CheckUser(user string, allowSelf bool, users []string) (bool, error) {
	if allowSelf {
		self, err := a.GetSelfCache()
		if err != nil {
			return false, err
		}
		if user == self.ScreenName {
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

func (a *TwitterAPI) RetweetWithChecker(name string, trimUser bool, cs ...TweetChecker) ([]anaconda.Tweet, error) {
	v := url.Values{}
	v.Set("screen_name", name)
	latestID, exists := a.cache.LatestTweetID[name]
	if exists {
		v.Set("since_id", fmt.Sprintf("%d", latestID))
	}
	tweets, err := a.GetUserTimeline(v)
	if err != nil {
		return nil, err
	}
	result := []anaconda.Tweet{}
	for i := len(tweets) - 1; i >= 0; i-- {
		t := tweets[i]
		match := true
		for _, c := range cs {
			m, err := c(t)
			if err != nil {
				return nil, err
			}
			if !m {
				match = false
				break
			}
		}
		if match {
			a.cache.LatestTweetID[name] = t.Id
			rt, err := a.Retweet(t.Id, trimUser)
			if err != nil {
				return nil, err
			}
			result = append(result, rt)
		}
	}
	return result, nil
}

// NotifyToAll sends metadata about the specified tweet to the all.
func (a *TwitterAPI) NotifyToAll(t *anaconda.Tweet, n *Notification) error {
	if n.Place != nil && t.HasCoordinates() {
		msg := fmt.Sprintf("ID: %s\nCountry: %s\nCreatedAt: %s", t.IdStr, t.Place.Country, t.CreatedAt)
		allowSelf := n.Place.AllowSelf
		users := n.Place.Users
		return a.PostDMToAll(msg, allowSelf, users)
	}
	return nil
}

// PostDMToAll posts the specified message to the all.
func (a *TwitterAPI) PostDMToAll(msg string, allowSelf bool, users []string) error {
	for _, user := range users {
		_, err := a.PostDMToScreenName(msg, user)
		if err != nil {
			return err
		}
	}
	if allowSelf {
		self, err := a.GetSelfCache()
		if err != nil {
			return err
		}
		_, err = a.PostDMToScreenName(msg, self.ScreenName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *TwitterAPI) Response(users []string, rs ...DirectMessageReceiver) error {
	dms, err := a.GetDirectMessages(nil)
	if err != nil {
		return err
	}
	senderToDM := make(map[string]anaconda.DirectMessage)
	for _, dm := range dms {
		sender := dm.SenderScreenName
		allowed, err := a.CheckUser(sender, false, users)
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
		latest, exists := a.cache.LatestDirectMessageID[sender]
		if !exists || latest != dm.Id {
			var text string
			for _, r := range rs {
				t, err := r(dm)
				if err != nil {
					return err
				}
				if t != "" {
					text = t
					break
				}
			}
			if text != "" {
				res, err := a.PostDMToScreenName(text, sender)
				if err != nil {
					return err
				}
				a.cache.LatestDirectMessageID[sender] = res.Id
			}
		}
	}
	return nil
}

// TweetChecker function checks if the specified tweet is acceptable, which means it
// should be retweeted.
type TweetChecker func(anaconda.Tweet) (bool, error)

// DirectMessageReceiver function receives the specified direct message and
// does something according to the received message.
// This returns a text and it is a reply for the above message's sender.
// Returning an empty string means this function does nothing.
type DirectMessageReceiver func(anaconda.DirectMessage) (string, error)

// DirectMessageEchoReceiver receives a direct message and does nothing, but
// returns the same text as the received one, so this is called `echo` receiver.
func DirectMessageEchoReceiver(m anaconda.DirectMessage) (string, error) {
	return m.Text, nil
}
