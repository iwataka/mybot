package mybot

import (
	"errors"
	"fmt"
	"html"
	"log"
	"net/url"
	"strings"

	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/models"
)

// NOTE: This must be fixed because multiple applications having different
// values cause infinite number of messages.
const msgPrefix = "<bot message>\n"

// TwitterAction can indicate for various actions for Twitter's tweets.
type TwitterAction struct {
	models.TwitterActionProperties
	Collections []string `json:"collections" toml:"collections"`
}

func NewTwitterAction() *TwitterAction {
	return &TwitterAction{
		Collections: []string{},
	}
}

func (a *TwitterAction) Add(action *TwitterAction) {
	// If action is nil, you have nothing to do
	if action == nil {
		return
	}

	a.Retweet = a.Retweet || action.Retweet
	a.Favorite = a.Favorite || action.Favorite
	a.Follow = a.Follow || action.Follow

	m := make(map[string]bool)
	for _, c := range a.Collections {
		m[c] = true
	}
	for _, c := range action.Collections {
		m[c] = true
	}
	cols := []string{}
	for c, exists := range m {
		if exists {
			cols = append(cols, c)
		}
	}
	a.Collections = cols
}

func (a *TwitterAction) Sub(action *TwitterAction) {
	// If action is nil, you have nothing to do
	if action == nil {
		return
	}

	a.Retweet = a.Retweet && !action.Retweet
	a.Favorite = a.Favorite && !action.Favorite
	a.Follow = a.Follow && !action.Follow

	m := make(map[string]bool)
	for _, c := range a.Collections {
		m[c] = true
	}
	for _, c := range action.Collections {
		m[c] = false
	}
	cols := []string{}
	for c, exists := range m {
		if exists {
			cols = append(cols, c)
		}
	}
	a.Collections = cols
}

func (a *TwitterAction) IsEmpty() bool {
	return !a.Retweet &&
		!a.Favorite &&
		!a.Follow &&
		len(a.Collections) == 0
}

// TwitterAPI is a wrapper of anaconda.TwitterApi.
type TwitterAPI struct {
	api    *anaconda.TwitterApi
	self   *anaconda.User
	cache  Cache
	config Config
	debug  bool
}

// NewTwitterAPI takes a user's authentication, cache and configuration and
// returns TwitterAPI instance for that user
func NewTwitterAPI(auth *OAuthCredentials, c Cache, cfg Config) *TwitterAPI {
	api := anaconda.NewTwitterApi(auth.AccessToken, auth.AccessTokenSecret)
	return &TwitterAPI{api, nil, c, cfg, true}
}

func (a *TwitterAPI) VerifyCredentials() (bool, error) {
	if a.api == nil {
		return false, nil
	} else {
		return a.api.VerifyCredentials()
	}
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

// ProcessTimeline gets tweets from the specified user's timeline and do action
// for tweets filtered by c.
func (a *TwitterAPI) ProcessTimeline(
	name string,
	v url.Values,
	c TweetChecker,
	vision VisionMatcher,
	lang LanguageMatcher,
	slack *SlackAPI,
	action *TweetAction,
) ([]anaconda.Tweet, error) {
	latestID, err := a.cache.GetLatestTweetID(name)
	if err != nil {
		return nil, err
	}
	v.Set("screen_name", name)
	if latestID > 0 {
		v.Set("since_id", fmt.Sprintf("%d", latestID))
	} else {
		// If the latest tweet ID doesn't exist, this fetches just the
		// latest tweet and store that ID.
		v.Set("count", "1")
	}
	tweets, err := a.api.GetUserTimeline(v)
	if err != nil {
		return nil, err
	}
	var pp TwitterPostProcessor
	if c.shouldRepeat() {
		pp = &TwitterPostProcessorEach{action, a.cache}
	} else {
		pp = &TwitterPostProcessorTop{name, a.cache}
	}
	result, err := a.processTweets(tweets, c, vision, lang, slack, action, pp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ProcessFavorites gets tweets from the specified user's favorite list and do
// action for tweets filtered by c.
func (a *TwitterAPI) ProcessFavorites(
	name string,
	v url.Values,
	c TweetChecker,
	vision VisionMatcher,
	lang LanguageMatcher,
	slack *SlackAPI,
	action *TweetAction,
) ([]anaconda.Tweet, error) {
	latestID, err := a.cache.GetLatestFavoriteID(name)
	if err != nil {
		return nil, err
	}
	v.Set("screen_name", name)
	if latestID > 0 {
		v.Set("since_id", fmt.Sprintf("%d", latestID))
	} else {
		// If the latest favorite ID doesn't exist, this fetches just
		// the latest tweet and store that ID.
		v.Set("count", "1")
	}
	tweets, err := a.api.GetFavorites(v)
	if err != nil {
		return nil, err
	}
	var pp TwitterPostProcessor
	if c.shouldRepeat() {
		pp = &TwitterPostProcessorEach{action, a.cache}
	} else {
		pp = &TwitterPostProcessorTop{name, a.cache}
	}
	result, err := a.processTweets(tweets, c, vision, lang, slack, action, pp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ProcessSearch gets tweets from search result by the specified query and do
// action for tweets filtered by c.
func (a *TwitterAPI) ProcessSearch(
	query string,
	v url.Values,
	c TweetChecker,
	vision VisionMatcher,
	lang LanguageMatcher,
	slack *SlackAPI,
	action *TweetAction,
) ([]anaconda.Tweet, error) {
	res, err := a.api.GetSearch(query, v)
	if err != nil {
		return nil, err
	}
	pp := &TwitterPostProcessorEach{action, a.cache}
	result, err := a.processTweets(res.Statuses, c, vision, lang, slack, action, pp)
	if err != nil {
		return nil, err
	}
	return result, err
}

type (
	TwitterPostProcessor interface {
		Process(anaconda.Tweet, bool) error
	}
	TwitterPostProcessorTop struct {
		screenName string
		cache      Cache
	}
	TwitterPostProcessorEach struct {
		action *TweetAction
		cache  Cache
	}
)

func (p *TwitterPostProcessorTop) Process(t anaconda.Tweet, match bool) error {
	id, err := p.cache.GetLatestTweetID(p.screenName)
	if err != nil {
		return err
	}
	if t.Id > id {
		err := p.cache.SetLatestTweetID(p.screenName, t.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *TwitterPostProcessorEach) Process(t anaconda.Tweet, match bool) error {
	if match {
		ac, err := p.cache.GetTweetAction(t.Id)
		if err != nil {
			return err
		}
		if ac != nil {
			ac.Add(p.action)
		} else {
			err := p.cache.SetTweetAction(t.Id, p.action)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *TwitterAPI) processTweets(
	tweets []anaconda.Tweet,
	c TweetChecker,
	v VisionMatcher,
	l LanguageMatcher,
	slack *SlackAPI,
	action *TweetAction,
	pp TwitterPostProcessor,
) ([]anaconda.Tweet, error) {
	result := []anaconda.Tweet{}
	// From the oldest to the newest
	for i := len(tweets) - 1; i >= 0; i-- {
		t := tweets[i]
		match, err := c.check(t, v, l, a.cache)
		if err != nil {
			return nil, err
		}
		if match {
			done, err := a.cache.GetTweetAction(t.Id)
			if err != nil {
				return nil, err
			}
			err = a.processTweet(t, action, done, slack)
			if err != nil {
				return nil, err
			}
			result = append(result, t)
		}
		err = pp.Process(t, match)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (a *TwitterAPI) processTweet(
	t anaconda.Tweet,
	action *TweetAction,
	done *TweetAction,
	slack *SlackAPI,
) error {
	ac := *action
	if done != nil {
		ac.Sub(done)
	}
	if ac.Twitter.Retweet && !t.Retweeted {
		var id int64
		if t.RetweetedStatus == nil {
			id = t.Id
		} else {
			id = t.RetweetedStatus.Id
		}
		_, err := a.api.Retweet(id, false)
		if CheckTwitterError(err) {
			return err
		}
	}
	if ac.Twitter.Favorite && !t.Favorited {
		_, err := a.api.Favorite(t.Id)
		if err != nil {
			return err
		}
	}
	if ac.Twitter.Follow {
		var screenName string
		if t.RetweetedStatus == nil {
			screenName = t.User.ScreenName
		} else {
			screenName = t.RetweetedStatus.User.ScreenName
		}
		_, err := a.api.FollowUser(screenName)
		if err != nil {
			return err
		}
	}
	for _, col := range ac.Twitter.Collections {
		err := a.collectTweet(t, col)
		if err != nil {
			return err
		}
	}

	if slack.Enabled() && ac.Slack != nil {
		for _, ch := range ac.Slack.Channels {
			err := slack.PostTweet(ch, t)
			if err != nil {
				return err
			}
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
	n, err := a.config.GetTwitterNotification()
	if err != nil {
		return err
	}
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

type TwitterUserListener struct {
	Stream *anaconda.Stream
	api    *TwitterAPI
}

func (l *TwitterUserListener) Listen(
	vis VisionMatcher,
	lang LanguageMatcher,
	slack *SlackAPI,
	cache Cache,
) error {
	for {
		switch c := (<-l.Stream.C).(type) {
		case anaconda.Tweet:
			if l.api.debug {
				log.Printf("Tweet by %s created at %s\n", c.User.Name, c.CreatedAt)
			}
			name := c.User.ScreenName
			timelines := []TimelineConfig{}
			ts, err := l.api.config.GetTwitterTimelines()
			if err != nil {
				return err
			}
			for _, t := range ts {
				for _, n := range t.ScreenNames {
					if n == name {
						timelines = append(timelines, t)
						break
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
				match, err := timeline.Filter.check(c, vis, lang, cache)
				if err != nil {
					return err
				}
				if match {
					done, err := l.api.cache.GetTweetAction(c.Id)
					if err != nil {
						return err
					}
					err = l.api.processTweet(c, timeline.Action, done, slack)
					if err != nil {
						return err
					}
					err = l.api.cache.SetLatestTweetID(name, c.Id)
					if err != nil {
						return err
					}
				}
			}
			err = l.api.cache.Save()
			if err != nil {
				return err
			}
		}
	}
}

// ListenUsers listens timelines of the friends
func (a *TwitterAPI) ListenUsers(v url.Values) (*TwitterUserListener, error) {
	if v == nil {
		v = url.Values{}
	}
	names, err := a.config.GetTwitterScreenNames()
	if err != nil {
		return nil, err
	}
	usernames := strings.Join(names, ",")
	if err != nil {
		return nil, err
	}
	if len(usernames) == 0 {
		return nil, errors.New("No user specified")
	} else {
		users, err := a.api.GetUsersLookup(usernames, nil)
		if err != nil {
			return nil, err
		}
		userids := []string{}
		for _, u := range users {
			userids = append(userids, u.IdStr)
		}
		v.Set("follow", strings.Join(userids, ","))
		stream := a.api.PublicStreamFilter(v)
		return &TwitterUserListener{stream, a}, nil
	}
}

type TwitterDMListener struct {
	Stream   *anaconda.Stream
	api      *TwitterAPI
	receiver DirectMessageReceiver
}

func (l *TwitterDMListener) Listen() error {
	for {
		switch c := (<-l.Stream.C).(type) {
		case anaconda.DirectMessage:
			if l.api.debug {
				log.Printf("Message by %s created at %s\n", c.Sender.Name, c.CreatedAt)
			}
			conf, err := l.api.config.GetInteraction()
			if err != nil {
				return err
			}
			if conf != nil {
				match, err := l.api.CheckUser(c.SenderScreenName, conf.AllowSelf, conf.Users)
				if err != nil {
					return err
				}
				if match {
					id, err := l.api.cache.GetLatestDMID()
					if err != nil {
						return err
					}
					if id < c.Id {
						err := l.api.cache.SetLatestDMID(c.Id)
						if err != nil {
							return err
						}
					}
					err = l.api.responseForDirectMessage(c, l.receiver)
					if err != nil {
						return err
					}
				}
			}
			err = l.api.cache.Save()
			if err != nil {
				return err
			}
		}
	}
}

// ListenMyself listens to the authenticated user by Twitter's User Streaming
// API and reacts with direct messages.
func (a *TwitterAPI) ListenMyself(v url.Values, receiver DirectMessageReceiver) (*TwitterDMListener, error) {
	ok, err := a.VerifyCredentials()
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, errors.New("Twitter Account Verification failed")
	}
	stream := a.api.UserStream(v)
	return &TwitterDMListener{stream, a, receiver}, nil
}

func (a *TwitterAPI) responseForDirectMessage(dm anaconda.DirectMessage, receiver DirectMessageReceiver) error {
	interaction, err := a.config.GetInteraction()
	if err != nil {
		return err
	}
	allowSelf := interaction.AllowSelf
	users := interaction.Users
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
	check(t anaconda.Tweet, v VisionMatcher, l LanguageMatcher, c Cache) (bool, error)
	shouldRepeat() bool
}

// DirectMessageReceiver function receives the specified direct message and
// does something according to the received message.
// This returns a text and it is a reply for the above message's sender.
// Returning an empty string means this function does nothing.
type DirectMessageReceiver func(anaconda.DirectMessage) (string, error)

type DirectMessageCommand struct {
	Name        string
	Description string
	Exec        func(*TwitterAPI, []string, []*DirectMessageCommand) (string, error)
}

var collectionsCommand = &DirectMessageCommand{
	Name:        "collections,cols",
	Description: "Shows a list of Twitter collections.",
	Exec: func(a *TwitterAPI, args []string, cmds []*DirectMessageCommand) (string, error) {
		if len(args) != 0 {
			return "This command can't accept any arguments", nil
		}

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
	},
}

var retweetCommand = &DirectMessageCommand{
	Name:        "retweet",
	Description: "Add configuration to retweet all tweet of the specified users",
	Exec: func(a *TwitterAPI, args []string, cmds []*DirectMessageCommand) (string, error) {
		timeline := NewTimelineConfig()
		timeline.ScreenNames = args
		timeline.Action.Twitter.Retweet = true
		err := a.config.AddTwitterTimeline(timeline)
		if err != nil {
			return "", err
		}
		err = a.config.Save()
		if err != nil {
			return "", err
		}
		return "Add configuration successfully", nil
	},
}

var favoriteCommand = &DirectMessageCommand{
	Name:        "favorite",
	Description: "Add configuration to favorite all favorites of the specified users",
	Exec: func(a *TwitterAPI, args []string, cmds []*DirectMessageCommand) (string, error) {
		favorite := NewFavoriteConfig()
		favorite.ScreenNames = args
		favorite.Action.Twitter.Favorite = true
		err := a.config.AddTwitterFavorite(favorite)
		if err != nil {
			return "", err
		}
		err = a.config.Save()
		if err != nil {
			return "", err
		}
		return "Add configuration successfully", nil
	},
}

var helpCommand = &DirectMessageCommand{
	Name:        "help,h",
	Description: "Shows the help text",
	Exec: func(a *TwitterAPI, args []string, cmds []*DirectMessageCommand) (string, error) {
		// If only slash is given, shows the help text.
		reply := "Use these commands with / at the head."
		for _, cmd := range cmds {
			reply += "\n"
			reply += fmt.Sprintf("  [%s]: %s", cmd.Name, cmd.Description)
		}
		return reply, nil
	},
}

var directMessageCommandList = []*DirectMessageCommand{
	collectionsCommand,
	helpCommand,
	retweetCommand,
	favoriteCommand,
}

// DefaultDirectMessageReceiver returns a reply from the specified direct
// message.
func (a *TwitterAPI) DefaultDirectMessageReceiver(m anaconda.DirectMessage) (string, error) {
	fields := strings.Fields(html.UnescapeString(m.Text))
	cmd := fields[0]
	args := []string{}
	if len(fields) > 1 {
		args = fields[1:]
	}
	// If the given command doesn't start with slash, ignore it.
	if !strings.HasPrefix(cmd, "/") {
		return "", nil
	} else if len(cmd) < 2 {
		return helpCommand.Exec(a, args, directMessageCommandList)
	}
	cmd = cmd[1:]

	for _, c := range directMessageCommandList {
		names := strings.Split(c.Name, ",")
		for _, name := range names {
			if cmd == name {
				return c.Exec(a, args, directMessageCommandList)
			}
		}
	}
	return "", nil
}

func CheckTwitterError(err error) bool {
	if err == nil {
		return false
	}

	switch twitterErr := err.(type) {
	case anaconda.TwitterError:
		// 327 means that duplicated message exists
		if twitterErr.Code == 327 {
			return false
		}
	case *anaconda.ApiError:
		for _, e := range twitterErr.Decoded.Errors {
			if CheckTwitterError(e) {
				return true
			}
		}
		return false
	}
	return true
}

func TwitterStatusURL(t anaconda.Tweet) string {
	srcFmt := "https://twitter.com/%s/status/%s"
	return fmt.Sprintf(srcFmt, t.User.IdStr, t.IdStr)
}