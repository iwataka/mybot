package core

import (
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/utils"
	"github.com/iwataka/slack"

	"container/list"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SlackAPI struct {
	api      models.SlackAPI
	config   Config
	cache    data.Cache
	msgQueue map[string]*concurrentQueue
	user     *slack.AuthTestResponse
}

type concurrentQueue struct {
	*list.List
	sync.Mutex
}

func newConcurrentQueue() *concurrentQueue {
	return &concurrentQueue{
		List: list.New(),
	}
}

func NewSlackAPIWithAuth(token string, config Config, cache data.Cache) *SlackAPI {
	var api models.SlackAPI
	if token != "" {
		api = slack.New(token)
	}
	return NewSlackAPI(api, config, cache)
}

func NewSlackAPI(api models.SlackAPI, config Config, cache data.Cache) *SlackAPI {
	return &SlackAPI{api, config, cache, make(map[string]*concurrentQueue), nil}
}

func (a *SlackAPI) Enabled() bool {
	return a.api != nil
}

func (a *SlackAPI) PostTweet(channel string, tweet anaconda.Tweet) error {
	text, params := convertFromTweetToSlackMsg(tweet)
	return a.PostMessage(channel, text, &params, false)
}

type SlackMsg struct {
	text   string
	params *slack.PostMessageParameters
}

func (a *SlackAPI) enqueueMsg(ch, text string, params *slack.PostMessageParameters) {
	if a.msgQueue[ch] == nil {
		a.msgQueue[ch] = newConcurrentQueue()
	}

	q := a.msgQueue[ch]
	q.Lock()
	defer q.Unlock()
	q.PushBack(&SlackMsg{text, params})
}

func (a *SlackAPI) dequeueMsg(ch string) *SlackMsg {
	q := a.msgQueue[ch]
	if q != nil {
		q.Lock()
		defer q.Unlock()
		front := q.Front()
		if front != nil {
			return q.Remove(front).(*SlackMsg)
		}
	}
	return nil
}

// TODO: Prevent infinite message loop
func (a *SlackAPI) PostMessage(channel, text string, params *slack.PostMessageParameters, channelIsOpen bool) error {
	var ps slack.PostMessageParameters
	if params != nil {
		ps = *params
	}
	_, _, err := a.api.PostMessage(channel, text, ps)
	if err != nil {
		if err.Error() == "channel_not_found" {
			// TODO: Prevent from creating multiple channels with the same name
			if channelIsOpen {
				_, err = a.api.CreateChannel(channel)
			} else {
				_, err = a.api.CreateGroup(channel)
			}
			if err != nil {
				if err.Error() == "user_is_bot" {
					err := a.notifyCreateChannel(channel)
					if err == nil {
						a.enqueueMsg(channel, text, params)
					}
					return utils.WithStack(err)
				} else {
					return utils.WithStack(err)
				}
			}
			_, _, err = a.api.PostMessage(channel, text, ps)
			if err != nil {
				return utils.WithStack(err)
			}
		} else {
			return utils.WithStack(err)
		}
	}
	return nil
}

func convertFromTweetToSlackMsg(t anaconda.Tweet) (string, slack.PostMessageParameters) {
	text := TwitterStatusURL(t)
	params := slack.PostMessageParameters{}
	params.IconURL = t.User.ProfileImageURL
	params.Username = fmt.Sprintf("%s@%s", t.User.Name, t.User.ScreenName)
	params.UnfurlLinks = true
	params.UnfurlMedia = true
	params.AsUser = false
	return text, params
}

func (a *SlackAPI) notifyCreateChannel(ch string) error {
	params := slack.PostMessageParameters{}
	msg := fmt.Sprintf("Create %s channel and invite me to it", ch)
	_, _, err := a.api.PostMessage("general", msg, params)
	return utils.WithStack(err)
}

func (a *SlackAPI) sendMsgQueues(ch string) error {
	q := a.msgQueue[ch]
	if q == nil {
		return nil
	}

	q.Lock()
	defer q.Unlock()
	for e := q.Front(); e != nil; e = e.Next() {
		m := e.Value.(*SlackMsg)
		err := a.PostMessage(ch, m.text, m.params, false)
		if err != nil {
			return utils.WithStack(err)
		}
	}
	a.msgQueue[ch] = newConcurrentQueue()
	return nil
}

func (a *SlackAPI) processMsgEvent(
	ch string,
	ev *slack.MessageEvent,
	vis VisionMatcher,
	lang LanguageMatcher,
	twitterAPI *TwitterAPI,
) ([]data.Action, error) {
	processedActions := []data.Action{}
	msgs := a.config.GetSlackMessages()
	for _, msg := range msgs {
		if !utils.CheckStringContained(msg.Channels, ch) {
			continue
		}
		match, err := msg.Filter.CheckSlackMsg(ev, vis, lang, a.cache)
		if err != nil {
			return nil, utils.WithStack(err)
		}
		if match {
			err := a.processMsgEventWithAction(ch, ev, msg.Action, twitterAPI)
			if err != nil {
				return nil, utils.WithStack(err)
			}
			processedActions = append(processedActions, msg.Action)
		}
	}
	return processedActions, nil
}

func (a *SlackAPI) processMsgEventWithAction(
	ch string,
	ev *slack.MessageEvent,
	action data.Action,
	twitterAPI *TwitterAPI,
) error {
	item := slack.NewRefToMessage(ev.Channel, ev.Timestamp)
	if action.Slack.Pin {
		err := a.api.AddPin(ev.Channel, item)
		if CheckSlackError(err) {
			return utils.WithStack(err)
		}
	}
	if action.Slack.Star {
		err := a.api.AddStar(ev.Channel, item)
		if CheckSlackError(err) {
			return utils.WithStack(err)
		}
	}
	for _, r := range action.Slack.Reactions {
		err := a.api.AddReaction(r, item)
		if CheckSlackError(err) {
			return utils.WithStack(err)
		}
	}
	for _, c := range action.Slack.Channels {
		if ch == c {
			continue
		}
		params := slack.PostMessageParameters{
			Attachments: ev.Attachments,
		}
		err := a.PostMessage(c, ev.Text, &params, false)
		if CheckSlackError(err) {
			return utils.WithStack(err)
		}
	}

	if action.Twitter.Tweet {
		_, err := twitterAPI.PostSlackMsg(ev.Text, ev.Attachments)
		if CheckTwitterError(err) {
			return utils.WithStack(err)
		}
	}
	return nil
}

func (a *SlackAPI) AuthTest() (*slack.AuthTestResponse, error) {
	if a.user != nil {
		return a.user, nil
	}
	if a.Enabled() {
		user, err := a.api.AuthTest()
		if err == nil {
			a.user = user
		}
		return user, utils.WithStack(err)
	}
	return nil, errors.New("Slack API is not available")
}

func (a *SlackAPI) Listen(vis VisionMatcher, lang LanguageMatcher, twitterAPI *TwitterAPI) *SlackListener {
	return &SlackListener{a, vis, lang, twitterAPI}
}

type SlackListener struct {
	api        *SlackAPI
	vis        VisionMatcher
	lang       LanguageMatcher
	twitterAPI *TwitterAPI
}

func (l *SlackListener) Start(ctx context.Context, outChan chan<- interface{}) error {
	rtm := l.api.api.NewRTM()
	go rtm.ManageConnection()
	defer func() { _ = rtm.Disconnect() }()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ChannelJoinedEvent:
				outChan <- NewReceivedEvent(SlackEventType, "channel joined", ev)
				err := l.api.sendMsgQueues(ev.Channel.Name)
				if err != nil {
					return utils.WithStack(err)
				}
			case *slack.GroupJoinedEvent:
				outChan <- NewReceivedEvent(SlackEventType, "group joined", ev)
				err := l.api.sendMsgQueues(ev.Channel.Name)
				if err != nil {
					return utils.WithStack(err)
				}
			case *slack.MessageEvent:
				t, err := parseSlackTimestamp(ev.Timestamp)
				if err != nil {
					return utils.WithStack(err)
				}
				if time.Since(*t)-time.Minute > 0 {
					continue
				}
				ch, err := getChannelNameByID(l.api.api, ev.Channel)
				if err != nil {
					return utils.WithStack(err)
				}
				if len(ch) > 0 {
					outChan <- NewReceivedEvent(SlackEventType, "message", ev)
					processedActions, err := l.api.processMsgEvent(ch, ev, l.vis, l.lang, l.twitterAPI)
					if err != nil {
						return utils.WithStack(err)
					}
					for _, a := range processedActions {
						outChan <- NewActionEvent(a, ev)
					}
				}
			case *slack.RTMError:
				return utils.WithStack(ev)
			case *slack.ConnectionErrorEvent:
				outChan <- NewReceivedEvent(SlackEventType, "connection error", ev)
				// Continue because ConnectionErrorEvent is treated as recoverable
				continue
			case *slack.InvalidAuthEvent:
				return fmt.Errorf(NewReceivedEvent("Slack", "invalid auth", ev).String())
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func getChannelNameByID(api models.SlackAPI, id string) (string, error) {
	chs, err := api.GetChannels(true)
	if err != nil {
		return "", utils.WithStack(err)
	}
	for _, c := range chs {
		if c.ID == id {
			return c.Name, nil
		}
	}
	grps, err := api.GetGroups(true)
	if err != nil {
		return "", utils.WithStack(err)
	}
	for _, g := range grps {
		if g.ID == id {
			return g.Name, nil
		}
	}
	return "", nil
}

func CheckSlackError(err error) bool {
	if err == nil {
		return false
	}

	switch err.Error() {
	case "already_reacted", "already_pinned", "already_starred", "internal_error":
		return false
	}
	return true
}

func parseSlackTimestamp(ts string) (*time.Time, error) {
	splittedTimestamp := strings.Split(ts, ".")
	sec, err := strconv.ParseInt(splittedTimestamp[0], 10, 64)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	nsec, err := strconv.ParseInt(splittedTimestamp[1], 10, 64)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	t := time.Unix(sec, nsec)
	return &t, nil
}
