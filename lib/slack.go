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

func (a *SlackAction) Add(action *SlackAction) {
	// If action is nil, you have nothing to do
	if action == nil {
		return
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
	a.Channels = chans
}

func (a *SlackAction) Sub(action *SlackAction) {
	// If action is nil, you have nothing to do
	if action == nil {
		return
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
	a.Channels = chans
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
	statusURL := TwitterStatusURL(tweet)
	text := fmt.Sprintf("%s\n%s created at %s", statusURL, tweet.User.Name, tweet.CreatedAt)

	att := slack.Attachment{}
	att.AuthorName = tweet.User.Name
	att.AuthorSubname = tweet.User.ScreenName
	att.Text = tweet.Text
	att.AuthorIcon = tweet.User.ProfileImageURL
	att.AuthorLink = tweet.User.URL

	params := slack.PostMessageParameters{}
	params.Attachments = []slack.Attachment{att}

	for _, m := range tweet.Entities.Media {
		a := slack.Attachment{}
		a.ImageURL = m.Media_url
		params.Attachments = append(params.Attachments, a)
	}

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

func (a *SlackAPI) notifyCreateChannel(ch string) error {
	params := slack.PostMessageParameters{}
	msg := fmt.Sprintf("Create #%s and invite me to it", ch)
	_, _, err := a.api.PostMessage("general", msg, params)
	return err
}
