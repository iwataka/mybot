package mybot

import (
	"errors"
	"fmt"
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

// msgPrefix is placed at the head of bot messages and indicates that messages
// are sent by bot. This prevents infinite-loop bot-to-bot communication.
const msgPrefix = "<bot message>\n"

// TwitterAPI is a wrapper of anaconda.TwitterApi.
type TwitterAPI struct {
	api    models.TwitterAPI
	config Config
	cache  data.Cache
	self   *anaconda.User
}

// NewTwitterAPIWithAuth takes a user's authentication, cache and configuration and
// returns TwitterAPI instance for that user
func NewTwitterAPIWithAuth(auth oauth.OAuthCreds, config Config, cache data.Cache) *TwitterAPI {
	at, ats := auth.GetCreds()
	api := anaconda.NewTwitterApi(at, ats)
	return NewTwitterAPI(api, config, cache)
}

func NewTwitterAPI(api models.TwitterAPI, config Config, cache data.Cache) *TwitterAPI {
	return &TwitterAPI{api, config, cache, nil}
}

func (a *TwitterAPI) BaseAPI() models.TwitterAPI {
	return a.api
}

func (a *TwitterAPI) VerifyCredentials() (bool, error) {
	if a.api == nil {
		return false, fmt.Errorf("Twitter API is not available")
	} else {
		return a.api.VerifyCredentials()
	}
}

func (a *TwitterAPI) PostSlackMsg(text string, atts []slack.Attachment) (anaconda.Tweet, error) {
	return a.api.PostTweet(text, nil)
}

// GetSelf gets the authenticated user's information and stores it as a cache,
// then returns it.
func (a *TwitterAPI) GetSelf() (anaconda.User, error) {
	if a.self == nil {
		self, err := a.api.GetSelf(nil)
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
	latestID := a.cache.GetLatestTweetID(name)
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
		return nil, utils.WithStack(err)
	}
	var pp TwitterPostProcessor
	if c.ShouldRepeat() {
		pp = &TwitterPostProcessorEach{action, a.cache}
	} else {
		pp = &TwitterPostProcessorTop{action, name, a.cache}
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
	latestID := a.cache.GetLatestFavoriteID(name)
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
		return nil, utils.WithStack(err)
	}
	var pp TwitterPostProcessor
	if c.ShouldRepeat() {
		pp = &TwitterPostProcessorEach{action, a.cache}
	} else {
		pp = &TwitterPostProcessorTop{action, name, a.cache}
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
	pp := &TwitterPostProcessorEach{action, a.cache}
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
		match, err := c.CheckTweet(t, v, l, a.cache)
		if err != nil {
			return nil, utils.WithStack(err)
		}
		if match {
			done := a.cache.GetTweetAction(t.Id)
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
		_, err := a.api.Retweet(id, false)
		if CheckTwitterError(err) {
			return utils.WithStack(err)
		}
		fmt.Printf("Retweet the tweet[%d]\n", id)
	}
	if action.Twitter.Favorite && !t.Favorited {
		id := t.Id
		_, err := a.api.Favorite(id)
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
	list, err := a.api.GetCollectionListByUserId(self.Id, nil)
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
			return utils.WithStack(err)
		}
		id = col.Response.TimelineId
	}
	_, err = a.api.AddEntryToCollection(id, tweet.Id, nil)
	if err != nil {
		return utils.WithStack(err)
	}
	return nil
}

// NotifyToAll sends metadata about the specified tweet, such as place, to the
// all users specified in the configuration.
func (a *TwitterAPI) NotifyToAll(slackAPI *SlackAPI, t *anaconda.Tweet) error {
	if t.HasCoordinates() {
		n := a.config.GetTwitterNotification()
		msg := fmt.Sprintf("ID: %s\nCountry: %s\nCreatedAt: %s", t.IdStr, t.Place.Country, t.CreatedAt)
		_, err := n.Place.Notify(a, slackAPI, msg)
		if err != nil {
			return utils.WithStack(err)
		}
	}
	return nil
}

func (a *TwitterAPI) GetSearch(query string, url url.Values) (anaconda.SearchResponse, error) {
	return a.api.GetSearch(query, url)
}

func (a *TwitterAPI) GetUserSearch(searchTerm string, v url.Values) ([]anaconda.User, error) {
	return a.api.GetUserSearch(searchTerm, v)
}

func (a *TwitterAPI) GetFavorites(vals url.Values) ([]anaconda.Tweet, error) {
	return a.api.GetFavorites(vals)
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
				ts := l.api.config.GetTwitterTimelines()
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
						done := l.api.cache.GetTweetAction(c.Id)
						err = l.api.processTweet(c, timeline.Action.Sub(done), slack)
						if err != nil {
							return utils.WithStack(err)
						}
						l.api.cache.SetLatestTweetID(name, c.Id)
					}
				}
				err := l.api.cache.Save()
				if err != nil {
					return utils.WithStack(err)
				}
			}
		case <-l.innerChan:
			return utils.NewStreamInterruptedError()
		}
	}
}

func (l *TwitterUserListener) Stop() error {
	l.stream.Stop()
	select {
	case l.innerChan <- true:
		return nil
	case <-time.After(l.timeout):
		return fmt.Errorf("Failed to stop twitter DM listener")
	}
}

// ListenUsers listens timelines of the friends
func (a *TwitterAPI) ListenUsers(v url.Values, timeout time.Duration) (*TwitterUserListener, error) {
	if v == nil {
		v = url.Values{}
	}
	names := a.config.GetTwitterScreenNames()
	usernames := strings.Join(names, ",")
	if len(usernames) == 0 {
		return nil, errors.New("No user specified")
	} else {
		users, err := a.api.GetUsersLookup(usernames, nil)
		if err != nil {
			return nil, utils.WithStack(err)
		}
		userids := []string{}
		for _, u := range users {
			userids = append(userids, u.IdStr)
		}
		v.Set("follow", strings.Join(userids, ","))
		stream := a.api.PublicStreamFilter(v)
		return &TwitterUserListener{stream, a, make(chan bool), timeout}, nil
	}
}

type TwitterDMListener struct {
	stream    *anaconda.Stream
	api       *TwitterAPI
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

				// TODO: Handle direct messages in the same way as the other sources
				id := l.api.cache.GetLatestDMID()
				if id < c.Id {
					l.api.cache.SetLatestDMID(c.Id)
				}
				err := l.api.cache.Save()
				if err != nil {
					return utils.WithStack(err)
				}
			}
		case <-l.innerChan:
			return utils.NewStreamInterruptedError()
		}
	}
}

func (l *TwitterDMListener) Stop() error {
	l.stream.Stop()
	select {
	case l.innerChan <- true:
		return nil
	case <-time.After(l.timeout):
		return fmt.Errorf("Failed to stop twitter DM listener")
	}
}

// ListenMyself listens to the authenticated user by Twitter's User Streaming
// API and reacts with direct messages.
func (a *TwitterAPI) ListenMyself(v url.Values, timeout time.Duration) (*TwitterDMListener, error) {
	ok, err := a.VerifyCredentials()
	if err != nil {
		return nil, utils.WithStack(err)
	} else if !ok {
		return nil, errors.New("Twitter Account Verification failed")
	}
	stream := a.api.UserStream(v)
	return &TwitterDMListener{stream, a, make(chan bool), timeout}, nil
}

// TweetChecker function checks if the specified tweet is acceptable, which means it
// should be retweeted.
type TweetChecker interface {
	CheckTweet(t anaconda.Tweet, v VisionMatcher, l LanguageMatcher, c data.Cache) (bool, error)
	ShouldRepeat() bool
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
