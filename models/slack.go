package models

import (
	"github.com/iwataka/slack"
)

type SlackAPI interface {
	PostMessage(ch, msg string, params slack.PostMessageParameters) (string, string, error)
	CreateChannel(name string) (*slack.Channel, error)
	CreateGroup(group string) (*slack.Group, error)
	NewRTM() *slack.RTM
	GetChannels(excludeArchived bool) ([]slack.Channel, error)
	AddPin(ch string, item slack.ItemRef) error
	AddStar(ch string, item slack.ItemRef) error
	AddReaction(name string, item slack.ItemRef) error
	AuthTest() (*slack.AuthTestResponse, error)
}

type SlackActionProperties struct {
	Pin  bool `json:"pin" toml:"pin" bson:"pin"`
	Star bool `json:"star" toml:"star" bson:"star"`
}
