package main

import (
	"fmt"
	"html"
	"log"
	"net/url"
	"strings"

	"github.com/iwataka/anaconda"
)

// NOTE: This must be fixed because multiple applications having different
// values cause infinite number of messages.
const msgPrefix = "<bot message>\n"

// TwitterAPI is a wrapper of anaconda.TwitterApi.
type TwitterAPI struct {
	api    *anaconda.TwitterApi
	self   *anaconda.User
	cache  *MybotCache
	config *MybotConfig
	debug  bool
}

// TwitterAuth contains values required for Twitter's user authentication.
type TwitterAuth struct {
	ConsumerKey       string `toml:"consumer_key"`
	ConsumerSecret    string `toml:"consumer_secret"`
	AccessToken       string `toml:"access_token"`
	AccessTokenSecret string `toml:"access_token_secret"`
}

// TwitterAction can indicate for various actions for Twitter's tweets.
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

// NewTwitterAPI takes a user's authentication, cache and configuration and
// returns TwitterAPI instance for that user
func NewTwitterAPI(a *TwitterAuth, c *MybotCache, cfg *MybotConfig) *TwitterAPI {
	anaconda.SetConsumerKey(a.ConsumerKey)
	anaconda.SetConsumerSecret(a.ConsumerSecret)
	api := anaconda.NewTwitterApi(a.AccessToken, a.AccessTokenSecret)
	return &TwitterAPI{api, nil, c, cfg, false}
}

// SetDebug enables/disables TwitterAPI's debug mode.
func (a *TwitterAPI) SetDebug(val bool) {
	a.debug = val
}

// PostDMToScreenName wraps anaconda.TwitterApi#PostDMToScreenName and has
// almost same function as the wrapped one, but posts messages with the
// specified prefix.
func (a *TwitterAPI) PostDMToScreenName(msg, name string) (anaconda.DirectMessage, error) {
	return a.api.PostDMToScreenName(msgPrefix+msg, name)
}

// GetCollectionListByUserId is just a wrapper of anaconda.TwitterApi#GetCollectionListByUserId
func (a *TwitterAPI) GetCollectionListByUserId(userId int64, v url.Values) (anaconda.CollectionListResult, error) {
	return a.api.GetCollectionListByUserId(userId, v)
}

// PostTweet is just a wrapper of anaconda.TwitterApi#PostTweet
func (a *TwitterAPI) PostTweet(msg string, v url.Values) (anaconda.Tweet, error) {
	return a.api.PostTweet(msg, v)
}

// GetFriendsList is just a wrapper of anaconda.TwitterApi#GetFriendsList
func (a *TwitterAPI) GetFriendsList(v url.Values) (anaconda.UserCursor, error) {
	return a.api.GetFriendsList(v)
}

// GetSelf gets the authenticated user's information and stores it as a cache,
// then returns it.
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

// CheckUser cheks if user is matched for the given allowSelf and users
// arguments.
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

// DoForAccount gets tweets from the specified user's timeline and do action
// for tweets filtered by c.
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

// DoForFavorites gets tweets from the specified user's favorite list and do
// action for tweets filtered by c.
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

// DoForSearch gets tweets from search result by the specified query and do
// action for tweets filtered by c.
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
		var id int64
		if t.RetweetedStatus == nil {
			id = t.Id
		} else {
			id = t.RetweetedStatus.Id
		}
		_, err := a.api.Retweet(id, false)
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
		var screenName string
		if t.RetweetedStatus == nil {
			screenName = t.User.ScreenName
		} else {
			screenName = t.RetweetedStatus.User.ScreenName
		}
		_, err := a.api.FollowUser(screenName)
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

// NotifyToAll sends metadata about the specified tweet, such as place, to the
// all users specified in the configuration.
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

// PostDMToAll posts the specified message to the all users specified in the
// configuration.
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

// Listen listens to the authenticated user by Twitter's User Streaming API and
// reacts with direct messages.
func (a *TwitterAPI) Listen(v url.Values, receiver DirectMessageReceiver, file string) error {
	if v == nil {
		v = url.Values{}
	}
	v.Set("with", "followings")
	stream := a.api.UserStream(v)
	for {
		switch c := (<-stream.C).(type) {
		case anaconda.DirectMessage:
			if a.debug {
				log.Printf("[Receive direct message]\n%v\n", c)
			}
			if a.config.Interaction != nil {
				conf := a.config.Interaction
				match, err := a.CheckUser(c.SenderScreenName, conf.AllowSelf, conf.Users)
				if err != nil {
					return err
				}
				if match {
					if a.cache.LatestDMID < c.Id {
						a.cache.LatestDMID = c.Id
					}
					err := a.responseForDirectMessage(c, receiver)
					if err != nil {
						return err
					}
				}
				err = a.cache.Save(file)
				if err != nil {
					return err
				}
			}
		case anaconda.Tweet:
			if a.debug {
				log.Printf("[Receive tweet]\n%s\n", c)
			}
			name := c.User.ScreenName
			timelines := []TimelineConfig{}
			for _, t := range a.config.Twitter.Timelines {
				if t.ScreenName != nil && *t.ScreenName == name {
					timelines = append(timelines, t)
				} else if t.ScreenNames != nil {
					for _, n := range t.ScreenNames {
						if n == name {
							timelines = append(timelines, t)
							break
						}
					}
				}
			}
			for _, timeline := range timelines {
				if timeline.ExcludeReplies != nil && *timeline.ExcludeReplies && c.InReplyToScreenName != "" {
					continue
				}
				if timeline.IncludeRts != nil && !*timeline.IncludeRts && c.RetweetedStatus != nil {
					continue
				}
				match, err := timeline.Filter.check(c)
				if err != nil {
					return err
				}
				if match {
					done := a.cache.Tweet2Action[c.IdStr]
					err := a.processTweet(c, timeline.Action, done)
					if err != nil {
						return err
					}
					a.cache.LatestTweetID[name] = c.Id
					err = a.cache.Save(file)
					if err != nil {
						return err
					}
				}
			}
		case anaconda.EventTweet:
			if c.Event.Event == "favorite" {
				if a.debug {
					log.Printf("[Receive event]\n%v\n", c)
				}
				name := c.Source.ScreenName
				favorites := []FavoriteConfig{}
				for _, f := range a.config.Twitter.Favorites {
					if f.ScreenName != nil && *f.ScreenName == name {
						favorites = append(favorites, f)
					} else if f.ScreenNames != nil {
						for _, n := range f.ScreenNames {
							if n == name {
								favorites = append(favorites, f)
							}
						}
					}
				}
				for _, favorite := range favorites {
					tweet := *c.TargetObject
					match, err := favorite.Filter.check(tweet)
					if err != nil {
						return err
					}
					if match {
						done := a.cache.Tweet2Action[tweet.IdStr]
						err := a.processTweet(tweet, favorite.Action, done)
						if err != nil {
							return err
						}
						a.cache.LatestFavoriteID[name] = tweet.Id
						err = a.cache.Save(file)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

// Response gets direct messages sent to the authenticated user and react with
// them.
// This is currently DEPRECATED and replaced with Listen.
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

func (a *TwitterAPI) FollowAll() error {
	for _, t := range a.config.Twitter.Timelines {
		if t.ScreenName != nil {
			_, err := a.api.FollowUser(*t.ScreenName)
			if err != nil {
				return err
			}
		} else if t.ScreenNames != nil {
			for _, n := range t.ScreenNames {
				_, err := a.api.FollowUser(n)
				if err != nil {
					return err
				}
			}
		}
	}
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

// DefaultDirectMessageReceiver returns a reply from the specified direct
// message.
func (a *TwitterAPI) DefaultDirectMessageReceiver(m anaconda.DirectMessage) (string, error) {
	text := html.UnescapeString(m.Text)
	lowers := strings.ToLower(text)
	if lowers == "collection" || lowers == "cols" {
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
	} else if lowers == "configuration" || lowers == "config" || lowers == "conf" {
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
