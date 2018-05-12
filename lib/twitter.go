package mybot

import (
	"errors"
	"fmt"
	"html"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/oauth"
	"github.com/iwataka/mybot/utils"
	"github.com/iwataka/slack"
)

// NOTE: This must be fixed because multiple applications having different
// values cause infinite number of messages.
const msgPrefix = "<bot message>\n"

// TwitterAPI is a wrapper of anaconda.TwitterApi.
type TwitterAPI struct {
	API    models.TwitterAPI
	Cache  data.Cache
	Config Config
	self   *anaconda.User
}

// NewTwitterAPI takes a user's authentication, cache and configuration and
// returns TwitterAPI instance for that user
func NewTwitterAPI(auth oauth.OAuthCreds, c data.Cache, cfg Config) *TwitterAPI {
	at, ats := auth.GetCreds()
	api := anaconda.NewTwitterApi(at, ats)
	return &TwitterAPI{api, c, cfg, nil}
}

func (a *TwitterAPI) VerifyCredentials() (bool, error) {
	if a.API == nil {
		return false, fmt.Errorf("Twitter API is not available")
	} else {
		return a.API.VerifyCredentials()
	}
}

// PostDMToScreenName wraps anaconda.TwitterApi#PostDMToScreenName and has
// almost same function as the wrapped one, but posts messages with the
// specified prefix.
func (a *TwitterAPI) PostDMToScreenName(msg, name string) (anaconda.DirectMessage, error) {
	return a.API.PostDMToScreenName(msgPrefix+msg, name)
}

// GetCollectionListByUserId is just a wrapper of anaconda.TwitterApi#GetCollectionListByUserId
func (a *TwitterAPI) GetCollectionListByUserId(userId int64, v url.Values) (anaconda.CollectionListResult, error) {
	return a.API.GetCollectionListByUserId(userId, v)
}

// PostTweet is just a wrapper of anaconda.TwitterApi#PostTweet
func (a *TwitterAPI) PostTweet(msg string, v url.Values) (anaconda.Tweet, error) {
	return a.API.PostTweet(msg, v)
}

func (a *TwitterAPI) PostSlackMsg(text string, atts []slack.Attachment) (anaconda.Tweet, error) {
	return a.PostTweet(text, nil)
}

// GetFriendsList is just a wrapper of anaconda.TwitterApi#GetFriendsList
func (a *TwitterAPI) GetFriendsList(v url.Values) (anaconda.UserCursor, error) {
	return a.API.GetFriendsList(v)
}

// GetSelf gets the authenticated user's information and stores it as a cache,
// then returns it.
func (a *TwitterAPI) GetSelf() (anaconda.User, error) {
	if a.self == nil {
		self, err := a.API.GetSelf(nil)
		if err != nil {
			return anaconda.User{}, utils.WithStack(err)
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
			return false, utils.WithStack(err)
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
	action data.Action,
) ([]anaconda.Tweet, error) {
	latestID := a.Cache.GetLatestTweetID(name)
	v.Set("screen_name", name)
	if latestID > 0 {
		v.Set("since_id", fmt.Sprintf("%d", latestID))
	} else {
		// If the latest tweet ID doesn't exist, this fetches just the
		// latest tweet and store that ID.
		v.Set("count", "1")
	}
	tweets, err := a.API.GetUserTimeline(v)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	var pp TwitterPostProcessor
	if c.ShouldRepeat() {
		pp = &TwitterPostProcessorEach{action, a.Cache}
	} else {
		pp = &TwitterPostProcessorTop{action, name, a.Cache}
	}
	result, err := a.processTweets(tweets, c, vision, lang, slack, action, pp)
	if err != nil {
		return nil, utils.WithStack(err)
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
	action data.Action,
) ([]anaconda.Tweet, error) {
	latestID := a.Cache.GetLatestFavoriteID(name)
	v.Set("screen_name", name)
	if latestID > 0 {
		v.Set("since_id", fmt.Sprintf("%d", latestID))
	} else {
		// If the latest favorite ID doesn't exist, this fetches just
		// the latest tweet and store that ID.
		v.Set("count", "1")
	}
	tweets, err := a.API.GetFavorites(v)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	var pp TwitterPostProcessor
	if c.ShouldRepeat() {
		pp = &TwitterPostProcessorEach{action, a.Cache}
	} else {
		pp = &TwitterPostProcessorTop{action, name, a.Cache}
	}
	result, err := a.processTweets(tweets, c, vision, lang, slack, action, pp)
	if err != nil {
		return nil, utils.WithStack(err)
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
	action data.Action,
) ([]anaconda.Tweet, error) {
	res, err := a.GetSearch(query, v)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	pp := &TwitterPostProcessorEach{action, a.Cache}
	result, err := a.processTweets(res.Statuses, c, vision, lang, slack, action, pp)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	return result, utils.WithStack(err)
}

type (
	TwitterPostProcessor interface {
		Process(anaconda.Tweet, bool) error
	}
	TwitterPostProcessorTop struct {
		action     data.Action
		screenName string
		cache      data.Cache
	}
	TwitterPostProcessorEach struct {
		action data.Action
		cache  data.Cache
	}
)

func (p *TwitterPostProcessorTop) Process(t anaconda.Tweet, match bool) error {
	id := p.cache.GetLatestTweetID(p.screenName)
	if t.Id > id {
		p.cache.SetLatestTweetID(p.screenName, t.Id)
	}
	if match {
		ac := p.cache.GetTweetAction(t.Id)
		p.cache.SetTweetAction(t.Id, ac.Add(p.action))
	}
	return nil
}

func (p *TwitterPostProcessorEach) Process(t anaconda.Tweet, match bool) error {
	if match {
		ac := p.cache.GetTweetAction(t.Id)
		p.cache.SetTweetAction(t.Id, ac.Add(p.action))
	}
	return nil
}

func (a *TwitterAPI) processTweets(
	tweets []anaconda.Tweet,
	c TweetChecker,
	v VisionMatcher,
	l LanguageMatcher,
	slack *SlackAPI,
	action data.Action,
	pp TwitterPostProcessor,
) ([]anaconda.Tweet, error) {
	result := []anaconda.Tweet{}
	// From the oldest to the newest
	for i := len(tweets) - 1; i >= 0; i-- {
		t := tweets[i]
		match, err := c.CheckTweet(t, v, l, a.Cache)
		if err != nil {
			return nil, utils.WithStack(err)
		}
		if match {
			done := a.Cache.GetTweetAction(t.Id)
			err = a.processTweet(t, action.Sub(done), slack)
			if err != nil {
				return nil, utils.WithStack(err)
			}
			result = append(result, t)
		}
		err = pp.Process(t, match)
		if err != nil {
			return nil, utils.WithStack(err)
		}
	}
	return result, nil
}

func (a *TwitterAPI) processTweet(
	t anaconda.Tweet,
	action data.Action,
	slack *SlackAPI,
) error {
	if action.Twitter.Retweet && !t.Retweeted {
		var id int64
		if t.RetweetedStatus == nil {
			id = t.Id
		} else {
			id = t.RetweetedStatus.Id
		}
		_, err := a.API.Retweet(id, false)
		if CheckTwitterError(err) {
			return utils.WithStack(err)
		}
		fmt.Printf("Retweet the tweet[%d]\n", id)
	}
	if action.Twitter.Favorite && !t.Favorited {
		id := t.Id
		_, err := a.API.Favorite(id)
		if err != nil {
			return utils.WithStack(err)
		}
		fmt.Println("Favorite the tweet[%d]", id)
	}
	for _, col := range action.Twitter.Collections {
		err := a.collectTweet(t, col)
		if err != nil {
			return utils.WithStack(err)
		}
		fmt.Printf("Collect the tweet[%s] to %s\n", t.IdStr, col)
	}

	if slack.Enabled() {
		for _, ch := range action.Slack.Channels {
			err := slack.PostTweet(ch, t)
			if err != nil {
				return utils.WithStack(err)
			}
			fmt.Printf("Send the tweet[%s] to #%s\n", t.IdStr, ch)
		}
	}

	return nil
}

func (a *TwitterAPI) collectTweet(tweet anaconda.Tweet, collection string) error {
	self, err := a.GetSelf()
	if err != nil {
		return utils.WithStack(err)
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
		col, err := a.API.CreateCollection(collection, nil)
		if err != nil {
			return utils.WithStack(err)
		}
		id = col.Response.TimelineId
	}
	_, err = a.API.AddEntryToCollection(id, tweet.Id, nil)
	if err != nil {
		return utils.WithStack(err)
	}
	return nil
}

// NotifyToAll sends metadata about the specified tweet, such as place, to the
// all users specified in the configuration.
func (a *TwitterAPI) NotifyToAll(t *anaconda.Tweet) error {
	n := a.Config.GetTwitterNotification()
	if t.HasCoordinates() {
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
			return utils.WithStack(err)
		}
	}
	if allowSelf {
		self, err := a.GetSelf()
		if err != nil {
			return utils.WithStack(err)
		}
		_, err = a.PostDMToScreenName(msg, self.ScreenName)
		if err != nil {
			return utils.WithStack(err)
		}
	}
	return nil
}

func (a *TwitterAPI) GetSearch(query string, url url.Values) (anaconda.SearchResponse, error) {
	return a.API.GetSearch(query, url)
}

func (a *TwitterAPI) GetUserSearch(searchTerm string, v url.Values) ([]anaconda.User, error) {
	return a.API.GetUserSearch(searchTerm, v)
}

type TwitterUserListener struct {
	stream    *anaconda.Stream
	api       *TwitterAPI
	innerChan chan bool
	timeout   time.Duration
}

func (l *TwitterUserListener) Listen(
	vis VisionMatcher,
	lang LanguageMatcher,
	slack *SlackAPI,
	cache data.Cache,
) error {
	for {
		select {
		case msg := <-l.stream.C:
			switch c := msg.(type) {
			case anaconda.Tweet:
				name := c.User.ScreenName
				timelines := []TimelineConfig{}
				ts := l.api.Config.GetTwitterTimelines()
				for _, t := range ts {
					for _, n := range t.ScreenNames {
						if n == name {
							timelines = append(timelines, t)
							break
						}
					}
				}

				if len(timelines) != 0 {
					fmt.Printf("Tweet[%s] created by %s at %s\n", c.IdStr, name, c.CreatedAt)
				}

				for _, timeline := range timelines {
					if timeline.ExcludeReplies != nil && *timeline.ExcludeReplies && c.InReplyToScreenName != "" {
						continue
					}
					if timeline.IncludeRts != nil && !*timeline.IncludeRts && c.RetweetedStatus != nil {
						continue
					}
					match, err := timeline.Filter.CheckTweet(c, vis, lang, cache)
					if err != nil {
						return utils.WithStack(err)
					}
					if match {
						done := l.api.Cache.GetTweetAction(c.Id)
						err = l.api.processTweet(c, timeline.Action.Sub(done), slack)
						if err != nil {
							return utils.WithStack(err)
						}
						l.api.Cache.SetLatestTweetID(name, c.Id)
					}
				}
				err := l.api.Cache.Save()
				if err != nil {
					return utils.WithStack(err)
				}
			}
		case <-l.innerChan:
			return utils.NewStreamInterruptedError()
		}
	}
}

func (l *TwitterUserListener) Stop() {
	l.stream.Stop()
	select {
	case l.innerChan <- true:
	case <-time.After(l.timeout):
		log.Println("Failed to stop twitter DM listener")
	}
}

// ListenUsers listens timelines of the friends
func (a *TwitterAPI) ListenUsers(v url.Values, timeout time.Duration) (*TwitterUserListener, error) {
	if v == nil {
		v = url.Values{}
	}
	names := a.Config.GetTwitterScreenNames()
	usernames := strings.Join(names, ",")
	if len(usernames) == 0 {
		return nil, errors.New("No user specified")
	} else {
		users, err := a.API.GetUsersLookup(usernames, nil)
		if err != nil {
			return nil, utils.WithStack(err)
		}
		userids := []string{}
		for _, u := range users {
			userids = append(userids, u.IdStr)
		}
		v.Set("follow", strings.Join(userids, ","))
		stream := a.API.PublicStreamFilter(v)
		return &TwitterUserListener{stream, a, make(chan bool), timeout}, nil
	}
}

type TwitterDMListener struct {
	stream    *anaconda.Stream
	api       *TwitterAPI
	receiver  DirectMessageReceiver
	innerChan chan bool
	timeout   time.Duration
}

func (l *TwitterDMListener) Listen() error {
	for {
		select {
		case msg := <-l.stream.C:
			switch c := msg.(type) {
			case anaconda.DirectMessage:
				fmt.Printf("DM[%s] created by %s at %s\n", c.IdStr, c.Sender.ScreenName, c.CreatedAt)

				conf := l.api.Config.GetTwitterInteraction()
				match, err := l.api.CheckUser(c.SenderScreenName, conf.AllowSelf, conf.Users)
				if err != nil {
					return utils.WithStack(err)
				}
				if match {
					id := l.api.Cache.GetLatestDMID()
					if id < c.Id {
						l.api.Cache.SetLatestDMID(c.Id)
					}
					err = l.api.responseForDirectMessage(c, l.receiver)
					if err != nil {
						return utils.WithStack(err)
					}
				}
				err = l.api.Cache.Save()
				if err != nil {
					return utils.WithStack(err)
				}
			}
		case <-l.innerChan:
			return utils.NewStreamInterruptedError()
		}
	}
}

func (l *TwitterDMListener) Stop() {
	l.stream.Stop()
	select {
	case l.innerChan <- true:
	case <-time.After(l.timeout):
		log.Println("Failed to stop twitter DM listener")
	}
}

// ListenMyself listens to the authenticated user by Twitter's User Streaming
// API and reacts with direct messages.
func (a *TwitterAPI) ListenMyself(v url.Values, receiver DirectMessageReceiver, timeout time.Duration) (*TwitterDMListener, error) {
	ok, err := a.VerifyCredentials()
	if err != nil {
		return nil, utils.WithStack(err)
	} else if !ok {
		return nil, errors.New("Twitter Account Verification failed")
	}
	stream := a.API.UserStream(v)
	return &TwitterDMListener{stream, a, receiver, make(chan bool), timeout}, nil
}

func (a *TwitterAPI) responseForDirectMessage(dm anaconda.DirectMessage, receiver DirectMessageReceiver) error {
	interaction := a.Config.GetTwitterInteraction()
	allowSelf := interaction.AllowSelf
	users := interaction.Users
	if strings.HasPrefix(html.UnescapeString(dm.Text), msgPrefix) {
		return nil
	}
	sender := dm.Sender.ScreenName
	allowed, err := a.CheckUser(sender, allowSelf, users)
	if err != nil {
		return utils.WithStack(err)
	}
	if allowed {
		text, err := receiver(dm)
		if err != nil {
			return utils.WithStack(err)
		}
		if text != "" {
			_, err := a.PostDMToScreenName(text, sender)
			if err != nil {
				return utils.WithStack(err)
			}
		}
	}
	return nil
}

// TweetChecker function checks if the specified tweet is acceptable, which means it
// should be retweeted.
type TweetChecker interface {
	CheckTweet(t anaconda.Tweet, v VisionMatcher, l LanguageMatcher, c data.Cache) (bool, error)
	ShouldRepeat() bool
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
			return "", utils.WithStack(err)
		}
		res, err := a.GetCollectionListByUserId(self.Id, nil)
		if err != nil {
			return "", utils.WithStack(err)
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
		a.Config.AddTwitterTimeline(timeline)
		err := a.Config.Save()
		if err != nil {
			return "", utils.WithStack(err)
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
		a.Config.AddTwitterFavorite(favorite)
		err := a.Config.Save()
		if err != nil {
			return "", utils.WithStack(err)
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
	case *anaconda.TwitterError:
		// 130: Over capacity
		// 187: The status text has already been Tweeted by the authenticated account.
		// 327: You have already retweeted this tweet.
		switch twitterErr.Code {
		case 130, 131, 187, 327:
			return false
		}
	case anaconda.TwitterError:
		return CheckTwitterError(&twitterErr)
	case *anaconda.ApiError:
		code := twitterErr.StatusCode
		// Status code 5?? means server error
		if code >= 500 && code < 600 {
			return false
		}
		for _, e := range twitterErr.Decoded.Errors {
			if CheckTwitterError(e) {
				return true
			}
		}
		return false
	case anaconda.ApiError:
		return CheckTwitterError(&twitterErr)
	}

	return true
}

func TwitterStatusURL(t anaconda.Tweet) string {
	srcFmt := "https://twitter.com/%s/status/%s"
	return fmt.Sprintf(srcFmt, t.User.IdStr, t.IdStr)
}
