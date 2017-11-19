package mybot

import (
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/models"
	"github.com/nlopes/slack"

	"container/list"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type SlackAction struct {
	models.SlackActionProperties
	Reactions []string `json:"reactions" toml:"reactions" bson:"reactions"`
	Channels  []string `json:"channels" toml:"channels" bson:"channels"`
}

func NewSlackAction() *SlackAction {
	return &SlackAction{
		Channels:  []string{},
		Reactions: []string{},
	}
}

func (a *SlackAction) Add(action *SlackAction) *SlackAction {
	return a.op(action, true)
}

func (a *SlackAction) Sub(action *SlackAction) *SlackAction {
	return a.op(action, false)
}

func (a *SlackAction) op(action *SlackAction, add bool) *SlackAction {
	result := *a

	// If action is nil, you have nothing to do
	if action == nil {
		return &result
	}

	result.Pin = BoolOp(a.Pin, action.Pin, add)
	result.Star = BoolOp(a.Star, action.Star, add)
	result.Reactions = StringsOp(a.Reactions, action.Reactions, add)
	result.Channels = StringsOp(a.Channels, action.Channels, add)

	return &result

}

func (a *SlackAction) IsEmpty() bool {
	return !a.Pin &&
		!a.Star &&
		len(a.Channels) == 0 &&
		len(a.Reactions) == 0

}

type SlackAPI struct {
	api      models.SlackAPI
	config   Config
	cache    Cache
	msgQueue map[string]*list.List
	user     *slack.AuthTestResponse
}

func NewSlackAPI(token string, config Config, cache Cache) *SlackAPI {
	var api models.SlackAPI
	if token != "" {
		api = slack.New(token)
	}
	return &SlackAPI{api, config, cache, make(map[string]*list.List), nil}
}

func (a *SlackAPI) Enabled() bool {
	return a.api != nil
}

func (a *SlackAPI) PostTweet(channel string, tweet anaconda.Tweet) error {
	text, params := convertFromTweetToSlackMsg(tweet)
	return a.PostMessage(channel, text, &params, true)
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

func (a *SlackAPI) PostMessage(channel, text string, params *slack.PostMessageParameters, queue bool) error {
	var ps slack.PostMessageParameters
	if params != nil {
		ps = *params
	}
	_, _, err := a.api.PostMessage(channel, text, ps)
	if err != nil {
		if err.Error() == "channel_not_found" {
			_, err := a.api.CreateChannel(channel)
			if err != nil {
				if err.Error() == "user_is_bot" {
					err := a.notifyCreateChannel(channel)
					if queue && err == nil {
						a.enqueueMsg(channel, text, params)
					}
					return err
				} else {
					return err
				}
			}
			_, _, err = a.api.PostMessage(channel, text, ps)
			if err != nil {
				return err
			}
		} else {
			return err
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
	msg := fmt.Sprintf("Create #%s and invite me to it", ch)
	_, _, err := a.api.PostMessage("general", msg, params)
	return err
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
			return err
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
		if !StringsContains(msg.Channels, ch) {
			continue
		}
		match, err := msg.Filter.CheckSlackMsg(ev, vis, lang, a.cache)
		if err != nil {
			return err
		}
		if match {
			err := a.processMsgEventWithAction(ch, ev, msg.Action, twitterAPI)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *SlackAPI) processMsgEventWithAction(
	ch string,
	ev *slack.MessageEvent,
	action *Action,
	twitterAPI *TwitterAPI,
) error {
	item := slack.NewRefToMessage(ev.Channel, ev.Timestamp)
	if action.Slack.Pin {
		err := a.api.AddPin(ev.Channel, item)
		if CheckSlackError(err) {
			return err
		}
		fmt.Println("Pin the message")
	}
	if action.Slack.Star {
		err := a.api.AddStar(ev.Channel, item)
		if CheckSlackError(err) {
			return err
		}
		fmt.Println("Star the message")
	}
	for _, r := range action.Slack.Reactions {
		err := a.api.AddReaction(r, item)
		if CheckSlackError(err) {
			return err
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
		err := a.PostMessage(c, ev.Text, &params, true)
		if CheckSlackError(err) {
			return err
		}
		fmt.Printf("Send the message to %s\n", c)
	}

	if action.Twitter.Tweet {
		_, err := twitterAPI.PostSlackMsg(ev.Text, ev.Attachments)
		if CheckTwitterError(err) {
			return err
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
		return user, err
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
	defer func() { rtm.Disconnect() }()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ChannelJoinedEvent:
				fmt.Printf("Joined to %s\n", ev.Channel.Name)
				err := l.api.sendMsgQueues(ev.Channel.Name)
				if err != nil {
					return err
				}
			case *slack.GroupJoinedEvent:
				fmt.Printf("Joined to %s\n", ev.Channel.Name)
				err := l.api.sendMsgQueues(ev.Channel.Name)
				if err != nil {
					return err
				}
			case *slack.MessageEvent:
				t, err := parseSlackTimestamp(ev.Timestamp)
				if err != nil {
					return err
				}
				if time.Now().Sub(*t)-time.Minute > 0 {
					continue
				}
				chs, err := l.api.api.GetChannels(true)
				if err != nil {
					return err
				}
				ch := ""
				for _, c := range chs {
					if c.ID == ev.Channel {
						ch = c.Name
						break
					}
				}
				if ch != "" {
					fmt.Printf("Receive message sent to %s by %s\n", ch, ev.User)
					err = l.api.processMsgEvent(ch, ev, vis, lang, twitterAPI)
					if err != nil {
						return err
					}
				}
			case *slack.RTMError:
				return ev
			case *slack.ConnectionErrorEvent:
				return ev
			case *slack.InvalidAuthEvent:
				return fmt.Errorf("Invalid slack authentication")
			}
		case <-l.innerChan:
			return NewInterruptedError()
		}
	}
	return nil
}

func (l *SlackListener) Stop() {
	l.innerChan <- true
}

func CheckSlackError(err error) bool {
	if err == nil {
		return false
	}

	if err.Error() == "invalid_name" {
		log.Print(err)
		return false
	}
	if err.Error() == "already_reacted" {
		log.Print(err)
		return false
	}
	return true
}

func parseSlackTimestamp(ts string) (*time.Time, error) {
	splittedTimestamp := strings.Split(ts, ".")
	sec, err := strconv.ParseInt(splittedTimestamp[0], 10, 64)
	if err != nil {
		return nil, err
	}
	nsec, err := strconv.ParseInt(splittedTimestamp[1], 10, 64)
	if err != nil {
		return nil, err
	}
	t := time.Unix(sec, nsec)
	return &t, nil
}
