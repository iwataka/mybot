package main

import (
	"fmt"
	"html"
	"net/url"
	"strings"

	"github.com/iwataka/anaconda"
)

// NOTE: This must be fixed because multiple applications having different
// values cause infinite number of messages.
const msgPrefix = "<bot message>\n"

type TwitterAPI struct {
	api    *anaconda.TwitterApi
	self   *anaconda.User
	cache  *MybotCache
	config *MybotConfig
}

type TwitterAuth struct {
	ConsumerKey       string `toml:"consumer_key"`
	ConsumerSecret    string `toml:"consumer_secret"`
	AccessToken       string `toml:"access_token"`
	AccessTokenSecret string `toml:"access_token_secret"`
}

type TwitterAction struct {
	Retweet     bool     `toml:"retweet"`
	Favorite    bool     `toml:"favorite"`
	Follow      bool     `toml:"follow"`
	Collections []string `toml:"collections"`
}

func (a *TwitterAction) add(action *TwitterAction) {
	a.Retweet = a.Retweet || action.Retweet
	a.Favorite = a.Favorite || action.Favorite
	a.Follow = a.Follow || action.Follow
	cols := a.Collections
	for _, col := range a.Collections {
		exists := false
		for _, c := range action.Collections {
			if col == c {
				exists = true
			}
		}
		if !exists {
			cols = append(cols, col)
		}
	}
	a.Collections = cols
}

func (a *TwitterAction) sub(action *TwitterAction) {
	a.Retweet = a.Retweet && !action.Retweet
	a.Favorite = a.Favorite && !action.Favorite
	a.Follow = a.Follow && !action.Follow
	cols := []string{}
	for _, col := range a.Collections {
		exists := false
		for _, c := range action.Collections {
			if col == c {
				exists = true
			}
		}
		if !exists {
			cols = append(cols, col)
		}
	}
	a.Collections = cols
}

func NewTwitterAPI(a *TwitterAuth, c *MybotCache, cfg *MybotConfig) *TwitterAPI {
	anaconda.SetConsumerKey(a.ConsumerKey)
	anaconda.SetConsumerSecret(a.ConsumerSecret)
	api := anaconda.NewTwitterApi(a.AccessToken, a.AccessTokenSecret)
	return &TwitterAPI{api, nil, c, cfg}
}

func (a *TwitterAPI) PostDMToScreenName(msg, name string) (anaconda.DirectMessage, error) {
	return a.api.PostDMToScreenName(msgPrefix+msg, name)
}

func (a *TwitterAPI) GetCollectionListByUserId(userId int64, v url.Values) (anaconda.CollectionListResult, error) {
	return a.api.GetCollectionListByUserId(userId, v)
}

func (a *TwitterAPI) PostTweet(msg string, v url.Values) (anaconda.Tweet, error) {
	return a.api.PostTweet(msg, v)
}

func (a *TwitterAPI) GetFriendsList(v url.Values) (anaconda.UserCursor, error) {
	return a.api.GetFriendsList(v)
}

// GetSelfCache returns the user of this client
func (a *TwitterAPI) GetSelf() (anaconda.User, error) {
	if a.self == nil {
		self, err := a.api.GetSelf(nil)
		if err != nil {
			return anaconda.User{}, err
		}
		a.self = &self
	}
	return *a.self, nil
}

func (a *TwitterAPI) CheckUser(user string, allowSelf bool, users []string) (bool, error) {
	if allowSelf {
		self, err := a.GetSelf()
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

func (a *TwitterAPI) DoForAccount(name string, v url.Values, c TweetChecker, action *TwitterAction) ([]anaconda.Tweet, error) {
	latestID, exists := a.cache.LatestTweetID[name]
	v.Set("screen_name", name)
	if exists {
		v.Set("since_id", fmt.Sprintf("%d", latestID))
	}
	tweets, err := a.api.GetUserTimeline(v)
	if err != nil {
		return nil, err
	}
	var post postProcessor
	if c.shouldRepeat() {
		post = a.postProcessEach(action)
	} else {
		post = a.postProcess(name, a.cache.LatestTweetID)
	}
	result, err := a.doForTweets(tweets, c, action, post)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *TwitterAPI) DoForFavorites(name string, v url.Values, c TweetChecker, action *TwitterAction) ([]anaconda.Tweet, error) {
	latestID, exists := a.cache.LatestFavoriteID[name]
	v.Set("screen_name", name)
	if exists {
		v.Set("since_id", fmt.Sprintf("%d", latestID))
	}
	tweets, err := a.api.GetFavorites(v)
	if err != nil {
		return nil, err
	}
	var post postProcessor
	if c.shouldRepeat() {
		post = a.postProcessEach(action)
	} else {
		post = a.postProcess(name, a.cache.LatestFavoriteID)
	}
	result, err := a.doForTweets(tweets, c, action, post)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *TwitterAPI) DoForSearch(query string, v url.Values, c TweetChecker, action *TwitterAction) ([]anaconda.Tweet, error) {
	res, err := a.api.GetSearch(query, v)
	if err != nil {
		return nil, err
	}
	result, err := a.doForTweets(res.Statuses, c, action, a.postProcessEach(action))
	if err != nil {
		return nil, err
	}
	return result, err
}

type postProcessor func(anaconda.Tweet, bool) error

func (a *TwitterAPI) postProcess(name string, m map[string]int64) postProcessor {
	return func(t anaconda.Tweet, match bool) error {
		id, exists := m[name]
		if (exists && t.Id > id) || !exists {
			m[name] = t.Id
		}
		return nil
	}
}

func (a *TwitterAPI) postProcessEach(action *TwitterAction) postProcessor {
	return func(t anaconda.Tweet, match bool) error {
		if match {
			ac, exists := a.cache.Tweet2Action[t.IdStr]
			if exists {
				ac.add(action)
			} else {
				a.cache.Tweet2Action[t.IdStr] = action
			}
		}
		return nil
	}
}

func (a *TwitterAPI) doForTweets(tweets []anaconda.Tweet, c TweetChecker, action *TwitterAction, post postProcessor) ([]anaconda.Tweet, error) {
	result := []anaconda.Tweet{}
	// From the oldest to the newest
	for i := len(tweets) - 1; i >= 0; i-- {
		t := tweets[i]
		match, err := c.check(t)
		if err != nil {
			return nil, err
		}
		if match {
			done := a.cache.Tweet2Action[t.IdStr]
			err := a.processTweet(t, action, done)
			if err != nil {
				return nil, err
			}
			result = append(result, t)
		}
		err = post(t, match)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (a *TwitterAPI) processTweet(t anaconda.Tweet, action *TwitterAction, done *TwitterAction) error {
	ac := *action
	if done != nil {
		ac.sub(done)
	}
	if ac.Retweet && !t.Retweeted {
		_, err := a.api.Retweet(t.Id, false)
		if err != nil {
			e, ok := err.(*anaconda.ApiError)
			if ok {
				// Already retweeted
				if e.StatusCode != 403 {
					return e
				}
			} else {
				return err
			}
		}
	}
	if ac.Favorite && !t.Favorited {
		_, err := a.api.Favorite(t.Id)
		if err != nil {
			e, ok := err.(*anaconda.ApiError)
			if ok {
				// Already favorited
				if e.StatusCode != 403 {
					return e
				}
			} else {
				return err
			}
		}
	}
	if ac.Follow {
		_, err := a.api.FollowUser(t.User.ScreenName)
		if err != nil {
			e, ok := err.(*anaconda.ApiError)
			if ok {
				// He/She is already friend
				if e.StatusCode != 403 {
					return e
				}
			} else {
				return err
			}
		}
	}
	for _, col := range ac.Collections {
		err := a.collectTweet(t, col)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *TwitterAPI) collectTweet(tweet anaconda.Tweet, collection string) error {
	self, err := a.GetSelf()
	if err != nil {
		return err
	}
	list, err := a.GetCollectionListByUserId(self.Id, nil)
	exists := false
	var id string
	for i, t := range list.Objects.Timelines {
		if collection == t.Name {
			exists = true
			id = i
			break
		}
	}
	if !exists {
		col, err := a.api.CreateCollection(collection, nil)
		if err != nil {
			return err
		}
		id = col.Response.TimelineId
	}
	_, err = a.api.AddEntryToCollection(id, tweet.Id, nil)
	if err != nil {
		return err
	}
	return nil
}

// NotifyToAll sends metadata about the specified tweet to the all.
func (a *TwitterAPI) NotifyToAll(t *anaconda.Tweet) error {
	n := a.config.Twitter.Notification
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
		self, err := a.GetSelf()
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

func (a *TwitterAPI) Listen(v url.Values, receiver DirectMessageReceiver) error {
	if v == nil {
		v = url.Values{}
	}
	v.Set("with", "user")
	stream := a.api.UserStream(v)
	for {
		switch c := (<-stream.C).(type) {
		case anaconda.DirectMessage:
			if a.cache.LatestDMID < c.Id {
				a.cache.LatestDMID = c.Id
			}
			err := a.responseForDirectMessage(c, receiver)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Response is replaced with Listen
func (a *TwitterAPI) Response(receiver DirectMessageReceiver) error {
	latestID := a.cache.LatestDMID
	v := url.Values{}
	if latestID != 0 {
		v.Set("since_id", fmt.Sprintf("%d", latestID))
	}
	count := a.config.Interaction.Count
	if count != nil {
		v.Set("count", fmt.Sprintf("%d", *count))
	}
	dms, err := a.api.GetDirectMessages(v)
	if err != nil {
		return err
	}
	first := latestID == 0
	for _, dm := range dms {
		if dm.Id > latestID {
			latestID = dm.Id
		}
		if !first {
			err := a.responseForDirectMessage(dm, receiver)
			if err != nil {
				return err
			}
		}
	}
	a.cache.LatestDMID = latestID
	return nil
}

func (a *TwitterAPI) responseForDirectMessage(dm anaconda.DirectMessage, receiver DirectMessageReceiver) error {
	allowSelf := a.config.Interaction.AllowSelf
	users := a.config.Interaction.Users
	if strings.HasPrefix(html.UnescapeString(dm.Text), msgPrefix) {
		return nil
	}
	sender := dm.Sender.ScreenName
	allowed, err := a.CheckUser(sender, allowSelf, users)
	if err != nil {
		return err
	}
	if allowed {
		text, err := receiver(dm)
		if err != nil {
			return err
		}
		if text != "" {
			_, err := a.PostDMToScreenName(text, sender)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// TweetChecker function checks if the specified tweet is acceptable, which means it
// should be retweeted.
type TweetChecker interface {
	check(t anaconda.Tweet) (bool, error)
	shouldRepeat() bool
}

// DirectMessageReceiver function receives the specified direct message and
// does something according to the received message.
// This returns a text and it is a reply for the above message's sender.
// Returning an empty string means this function does nothing.
type DirectMessageReceiver func(anaconda.DirectMessage) (string, error)

func (a *TwitterAPI) DefaultDirectMessageReceiver(m anaconda.DirectMessage) (string, error) {
	text := html.UnescapeString(m.Text)
	if text == "collection" || text == "cols" {
		self, err := a.GetSelf()
		if err != nil {
			return "", err
		}
		res, err := a.GetCollectionListByUserId(self.Id, nil)
		if err != nil {
			return "", err
		}
		timelines := res.Objects.Timelines
		lines := []string{}
		for _, col := range timelines {
			line := fmt.Sprintf("%s: %s", col.Name, col.CollectionUrl)
			lines = append(lines, line)
		}
		return strings.Join(lines, "\n"), nil
	} else if text == "configuration" || text == "config" || text == "conf" {
		cfg := new(MybotConfig)
		*cfg = *config
		cfg.Authentication = nil
		bytes, err := cfg.TomlText(strings.Repeat(" ", 4))
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	} else {
		return fmt.Sprintf("Unknow command: %s", text), nil
	}
}
