package mybot

import (
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/utils"
	"github.com/iwataka/slack"

	"container/list"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type SlackAPI struct {
	api      models.SlackAPI
	config   Config
	cache    data.Cache
	msgQueue map[string]*list.List
	user     *slack.AuthTestResponse
}

func NewSlackAPIWithAuth(token string, config Config, cache data.Cache) *SlackAPI {
	var api models.SlackAPI
	if token != "" {
		api = slack.New(token)
	}
	return NewSlackAPI(api, config, cache)
}

func NewSlackAPI(api models.SlackAPI, config Config, cache data.Cache) *SlackAPI {
	return &SlackAPI{api, config, cache, make(map[string]*list.List), nil}
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
		a.msgQueue[ch] = list.New()
	}
	a.msgQueue[ch].PushBack(&SlackMsg{text, params})
}

func (a *SlackAPI) dequeueMsg(ch string) *SlackMsg {
	q := a.msgQueue[ch]
	if q != nil {
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
	for e := q.Front(); e != nil; e = e.Next() {
		m := e.Value.(*SlackMsg)
		err := a.PostMessage(ch, m.text, m.params, false)
		if err != nil {
			return utils.WithStack(err)
		}
	}
	a.msgQueue[ch] = list.New()
	return nil
}

func (a *SlackAPI) processMsgEvent(
	ch string,
	ev *slack.MessageEvent,
	vis VisionMatcher,
	lang LanguageMatcher,
	twitterAPI *TwitterAPI,
) error {
	msgs := a.config.GetSlackMessages()
	for _, msg := range msgs {
		if !utils.CheckStringContained(msg.Channels, ch) {
			continue
		}
		match, err := msg.Filter.CheckSlackMsg(ev, vis, lang, a.cache)
		if err != nil {
			return utils.WithStack(err)
		}
		if match {
			err := a.processMsgEventWithAction(ch, ev, msg.Action, twitterAPI)
			if err != nil {
				return utils.WithStack(err)
			}
		}
	}
	return nil
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
		fmt.Println("Pin the message")
	}
	if action.Slack.Star {
		err := a.api.AddStar(ev.Channel, item)
		if CheckSlackError(err) {
			return utils.WithStack(err)
		}
		fmt.Println("Star the message")
	}
	for _, r := range action.Slack.Reactions {
		err := a.api.AddReaction(r, item)
		if CheckSlackError(err) {
			return utils.WithStack(err)
		}
		fmt.Println("React to the message")
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
		fmt.Printf("Send the message to %s\n", c)
	}

	if action.Twitter.Tweet {
		_, err := twitterAPI.PostSlackMsg(ev.Text, ev.Attachments)
		if CheckTwitterError(err) {
			return utils.WithStack(err)
		}
		fmt.Println("Tweet the message")
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

func (a *SlackAPI) Listen() *SlackListener {
	return &SlackListener{a, make(chan bool)}
}

type SlackListener struct {
	api       *SlackAPI
	innerChan chan bool
}

func (l *SlackListener) Start(
	vis VisionMatcher,
	lang LanguageMatcher,
	twitterAPI *TwitterAPI,
) error {
	rtm := l.api.api.NewRTM()
	go rtm.ManageConnection()
	defer func() { _ = rtm.Disconnect() }()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ChannelJoinedEvent:
				fmt.Printf("Joined to %s\n", ev.Channel.Name)
				err := l.api.sendMsgQueues(ev.Channel.Name)
				if err != nil {
					return utils.WithStack(err)
				}
			case *slack.GroupJoinedEvent:
				fmt.Printf("Joined to %s\n", ev.Channel.Name)
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
				ch := ""
				if ch == "" {
					chs, err := l.api.api.GetChannels(true)
					if err != nil {
						return utils.WithStack(err)
					}
					for _, c := range chs {
						if c.ID == ev.Channel {
							ch = c.Name
							break
						}
					}
				}
				if ch == "" {
					grps, err := l.api.api.GetGroups(true)
					if err != nil {
						return utils.WithStack(err)
					}
					for _, g := range grps {
						if g.ID == ev.Channel {
							ch = g.Name
							break
						}
					}
				}
				if ch != "" {
					fmt.Printf("Receive message sent to %s by %s\n", ch, ev.User)
					err = l.api.processMsgEvent(ch, ev, vis, lang, twitterAPI)
					if err != nil {
						return utils.WithStack(err)
					}
				}
			case *slack.RTMError:
				return utils.WithStack(ev)
			case *slack.ConnectionErrorEvent:
				log.Println(ev)
				// Continue because ConnectionErrorEvent is treated as recoverable
				continue
			case *slack.InvalidAuthEvent:
				return fmt.Errorf("Invalid slack authentication")
			}
		case <-l.innerChan:
			return utils.NewStreamInterruptedError()
		}
	}
}

func (l *SlackListener) Stop() error {
	select {
	case l.innerChan <- true:
		return nil
	case <-time.After(time.Minute):
		return fmt.Errorf("Faield to stop slack listener (timeout: 1m)")
	}
}

func CheckSlackError(err error) bool {
	if err == nil {
		return false
	}

	if err.Error() == "invalid_name" {
		log.Printf("%+v\n", err)
		return false
	}
	if err.Error() == "already_reacted" {
		log.Printf("%+v\n", err)
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
