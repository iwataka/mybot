package mybot

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/models"
	"github.com/nlopes/slack"
)

type SlackAction struct {
	models.SlackActionProperties
	Reactions []string `json:"reactions" toml:"reactions"`
	Channels  []string `json:"channels" toml:"channels"`
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
	api    *slack.Client
	config Config
	cache  Cache
}

func NewSlackAPI(token string, config Config, cache Cache) *SlackAPI {
	var api *slack.Client
	if token != "" {
		api = slack.New(token)
	}
	return &SlackAPI{api, config, cache}
}

func (a *SlackAPI) Enabled() bool {
	return a.api != nil
}

func (a *SlackAPI) PostTweet(channel string, tweet anaconda.Tweet) error {
	text, params := convertFromTweetToSlackMsg(tweet)
	return a.PostMesage(channel, text, &params)
}

func (a *SlackAPI) PostMesage(channel, text string, params *slack.PostMessageParameters) error {
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
					return a.notifyCreateChannel(channel)
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

func (a *SlackAPI) Listen() *SlackListener {
	return &SlackListener{true, a}
}

type SlackListener struct {
	enabled bool
	api     *SlackAPI
}

func (l *SlackListener) Start(
	vis VisionMatcher,
	lang LanguageMatcher,
	twitterAPI *TwitterAPI,
) error {
	rtm := l.api.api.NewRTM()
	go rtm.ManageConnection()

	logFields := log.Fields{
		"type":   "slack",
		"action": "fetch",
	}

	for l.enabled {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				log.WithFields(
					logFields,
				).Infof("Message to %s by %s", ev.Channel, ev.Username)
				err := l.processMsgEvent(ev, vis, lang, twitterAPI)
				if err != nil {
					return err
				}
			case *slack.RTMError:
				log.WithFields(logFields).Infof("%T", ev)
				return ev
			case *slack.InvalidAuthEvent:
				log.WithFields(logFields).Infof("%T", ev)
				return fmt.Errorf("Invalid authentication")
			}
		}
	}
	return nil
}

func (l *SlackListener) Stop() {
	l.enabled = false
}

func (l *SlackListener) processMsgEvent(
	ev *slack.MessageEvent,
	vis VisionMatcher,
	lang LanguageMatcher,
	twitterAPI *TwitterAPI,
) error {
	ch := ev.Channel
	msgs, err := l.api.config.GetSlackMessages()
	if err != nil {
		return err
	}
	for _, msg := range msgs {
		if !StringsContains(msg.Channels, ch) {
			continue
		}
		match, err := msg.Filter.CheckSlackMsg(ev, vis, lang, l.api.cache)
		if err != nil {
			return err
		}
		if match {
			err := l.processMsgEventWithAction(ev, msg.Action, twitterAPI)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *SlackListener) processMsgEventWithAction(
	ev *slack.MessageEvent,
	action *Action,
	twitterAPI *TwitterAPI,
) error {
	logFields := log.Fields{
		"type":   "slack",
		"action": "process",
	}

	item := slack.NewRefToMessage(ev.Channel, ev.Timestamp)
	if action.Slack.Pin {
		err := l.api.api.AddPin(ev.Channel, item)
		if err != nil {
			return err
		}
		log.WithFields(logFields).Infoln("Pin the message")
	}
	if action.Slack.Star {
		err := l.api.api.AddStar(ev.Channel, item)
		if err != nil {
			return err
		}
		log.WithFields(logFields).Infoln("Star the message")
	}
	for _, r := range action.Slack.Reactions {
		err := l.api.api.AddReaction(r, item)
		if err != nil {
			return err
		}
		log.WithFields(logFields).Infoln("React to the message")
	}
	for _, c := range action.Slack.Channels {
		if ev.Channel == c {
			continue
		}
		params := slack.PostMessageParameters{
			Attachments: ev.Attachments,
		}
		err := l.api.PostMesage(c, ev.Text, &params)
		if err != nil {
			return err
		}
		log.WithFields(logFields).Infof("Send the message to %s", c)
	}

	if action.Twitter.Tweet {
		_, err := twitterAPI.PostSlackMsg(ev.Text, ev.Attachments)
		if err != nil {
			return err
		}
		log.WithFields(logFields).Infoln("Tweet the message")
	}
	return nil
}
