package core

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/oauth"
	"github.com/iwataka/mybot/utils"
	"github.com/iwataka/slack"
)

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
) ([]anaconda.Tweet, []data.Action, error) {
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
		return nil, nil, utils.WithStack(err)
	}
	var pp TwitterPostProcessor
	if c.ShouldRepeat() {
		pp = &TwitterPostProcessorEach{action, a.cache}
	} else {
		pp = &TwitterPostProcessorTop{action, name, a.cache}
	}
	processedTweets, processedActions, err := a.processTweets(tweets, c, vision, lang, slack, action, pp)
	if err != nil {
		return nil, nil, utils.WithStack(err)
	}
	return processedTweets, processedActions, nil
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
) ([]anaconda.Tweet, []data.Action, error) {
	res, err := a.GetSearch(query, v)
	if err != nil {
		return nil, nil, utils.WithStack(err)
	}
	pp := &TwitterPostProcessorEach{action, a.cache}
	processedTweets, processedActions, err := a.processTweets(res.Statuses, c, vision, lang, slack, action, pp)
	if err != nil {
		return nil, nil, utils.WithStack(err)
	}
	return processedTweets, processedActions, utils.WithStack(err)
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
) ([]anaconda.Tweet, []data.Action, error) {
	processedTweets := []anaconda.Tweet{}
	processedActions := []data.Action{}
	// From the oldest to the newest
	for i := len(tweets) - 1; i >= 0; i-- {
		t := tweets[i]
		match, err := c.CheckTweet(t, v, l, a.cache)
		if err != nil {
			return nil, nil, utils.WithStack(err)
		}
		if match {
			done := a.cache.GetTweetAction(t.Id)
			undone := action.Sub(done)
			err = a.processTweet(t, undone, slack)
			if err != nil {
				return nil, nil, utils.WithStack(err)
			}
			processedTweets = append(processedTweets, t)
			processedActions = append(processedActions, undone)
		}
		err = pp.Process(t, match)
		if err != nil {
			return nil, nil, utils.WithStack(err)
		}
	}
	return processedTweets, processedActions, nil
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
	}
	if action.Twitter.Favorite && !t.Favorited {
		id := t.Id
		_, err := a.api.Favorite(id)
		if CheckTwitterError(err) {
			return utils.WithStack(err)
		}
	}
	for _, col := range action.Twitter.Collections {
		err := a.collectTweet(t, col)
		if CheckTwitterError(err) {
			return utils.WithStack(err)
		}
	}

	if slack.Enabled() {
		for _, ch := range action.Slack.Channels {
			err := slack.PostTweet(ch, t)
			if CheckSlackError(err) {
				return utils.WithStack(err)
			}
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
	if err != nil {
		return utils.WithStack(err)
	}
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
	stream *anaconda.Stream
	api    *TwitterAPI
	vis    VisionMatcher
	lang   LanguageMatcher
	slack  *SlackAPI
	cache  data.Cache
}

// ListenUsers listens timelines of the friends
func (a *TwitterAPI) ListenUsers(
	v url.Values,
	vis VisionMatcher,
	lang LanguageMatcher,
	slack *SlackAPI,
	cache data.Cache,
) (*TwitterUserListener, error) {
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
		return &TwitterUserListener{stream, a, vis, lang, slack, cache}, nil
	}
}

func (l *TwitterUserListener) Listen(ctx context.Context, outChan chan<- interface{}) error {
	for {
		select {
		case msg := <-l.stream.C:
			switch m := msg.(type) {
			case anaconda.Tweet:
				name := m.User.ScreenName
				timelines := l.api.config.GetTwitterTimelinesByScreenName(name)
				if len(timelines) != 0 {
					outChan <- NewReceivedEvent(TwitterEventType, "tweet", m)
				}

				for _, timeline := range timelines {
					if !checkTweetByTimelineConfig(m, timeline) {
						continue
					}
					match, err := timeline.Filter.CheckTweet(m, l.vis, l.lang, l.cache)
					if err != nil {
						return utils.WithStack(err)
					}
					if !match {
						continue
					}
					done := l.api.cache.GetTweetAction(m.Id)
					undone := timeline.Action.Sub(done)
					if err := l.api.processTweet(m, undone, l.slack); err != nil {
						return utils.WithStack(err)
					}
					outChan <- NewActionEvent(undone, m)
					l.api.cache.SetLatestTweetID(name, m.Id)
				}
				err := l.api.cache.Save()
				if err != nil {
					return utils.WithStack(err)
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (l *TwitterUserListener) Stop() {
	l.stream.Stop()
}

func checkTweetByTimelineConfig(t anaconda.Tweet, c TimelineConfig) bool {
	if c.ExcludeReplies && t.InReplyToScreenName != "" {
		return false
	}
	if !c.IncludeRts && t.RetweetedStatus != nil {
		return false
	}
	return true
}

type TwitterDMListener struct {
	stream *anaconda.Stream
	api    *TwitterAPI
}

// ListenMyself listens to the authenticated user by Twitter's User Streaming
// API and reacts with direct messages.
func (a *TwitterAPI) ListenMyself(v url.Values) (*TwitterDMListener, error) {
	ok, err := a.VerifyCredentials()
	if err != nil {
		return nil, utils.WithStack(err)
	} else if !ok {
		return nil, errors.New("Twitter Account Verification failed")
	}
	stream := a.api.UserStream(v)
	return &TwitterDMListener{stream, a}, nil
}

func (l *TwitterDMListener) Listen(ctx context.Context, outChan chan<- interface{}) error {
	for {
		select {
		case msg := <-l.stream.C:
			switch c := msg.(type) {
			case anaconda.DirectMessage:
				outChan <- NewReceivedEvent(TwitterEventType, "DM", c)
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
		case <-ctx.Done():
			return nil
		}
	}
}

func (l *TwitterDMListener) Stop() {
	l.stream.Stop()
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
		// https://developer.twitter.com/ja/docs/basics/response-codes
		// 130: Over capacity
		// 131: Internal error
		// 139: You have already favorited this status.
		// 187: The status text has already been Tweeted by the authenticated account.
		// 327: You have already retweeted this tweet.
		switch twitterErr.Code {
		case 130, 131, 139, 187, 327:
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
