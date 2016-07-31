package main

import (
	"github.com/nlopes/slack"
)

var slackApi *slack.Client

func slackInitChannels(channels []string) error {
	cs, err := slackApi.GetChannels(false)
	if err != nil {
		return err
	}
	for _, channel := range channels {
		exists := false
		for _, c := range cs {
			if channel == c.Name {
				exists = true
				break
			}
		}
		if !exists {
			_, err := slackApi.CreateChannel(channel)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
