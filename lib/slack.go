package mybot

import (
	"fmt"

	"github.com/iwataka/anaconda"
	"github.com/nlopes/slack"
)

type SlackAction struct {
	Channels []string `json:"channels" toml:"channels"`
}

func NewSlackAction() *SlackAction {
	return &SlackAction{
		Channels: []string{},
	}
}

func (a *SlackAction) Add(action *SlackAction) *SlackAction {
	result := *a

	// If action is nil, you have nothing to do
	if action == nil {
		return &result
	}

	m := make(map[string]bool)
	for _, c := range a.Channels {
		m[c] = true
	}
	for _, c := range action.Channels {
		m[c] = true
	}
	chans := []string{}
	for c, exists := range m {
		if exists {
			chans = append(chans, c)
		}
	}
	result.Channels = chans

	return &result
}

func (a *SlackAction) Sub(action *SlackAction) *SlackAction {
	result := *a

	// If action is nil, you have nothing to do
	if action == nil {
		return &result
	}

	m := make(map[string]bool)
	for _, c := range a.Channels {
		m[c] = true
	}
	for _, c := range action.Channels {
		m[c] = false
	}
	chans := []string{}
	for c, exists := range m {
		if exists {
			chans = append(chans, c)
		}
	}
	result.Channels = chans

	return &result
}

func (a *SlackAction) IsEmpty() bool {
	return len(a.Channels) == 0
}

type SlackAPI struct {
	api *slack.Client
}

func NewSlackAPI(token string) *SlackAPI {
	var api *slack.Client
	if token != "" {
		api = slack.New(token)
	}
	return &SlackAPI{api}
}

func (a *SlackAPI) Enabled() bool {
	return a.api != nil
}

func (a *SlackAPI) PostTweet(channel string, tweet anaconda.Tweet) error {
	text, params := a.convertFromTweet(tweet)

	_, _, err := a.api.PostMessage(channel, text, params)
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
			err = a.PostTweet(channel, tweet)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (a *SlackAPI) convertFromTweet(t anaconda.Tweet) (string, slack.PostMessageParameters) {
	text := TwitterStatusURL(t)
	params := slack.PostMessageParameters{}
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
