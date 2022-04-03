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
	GetGroups(excludeArchived bool) ([]slack.Group, error)
	AddPin(ch string, item slack.ItemRef) error
	AddStar(ch string, item slack.ItemRef) error
	AddReaction(name string, item slack.ItemRef) error
	AuthTest() (*slack.AuthTestResponse, error)
}

type SlackActionProperties struct {
	Pin  bool `json:"pin" toml:"pin" bson:"pin" yaml:"pin"`
	Star bool `json:"star" toml:"star" bson:"star" yaml:"star"`
}

type SlackAPIImpl struct {
	api *slack.Client
}

func NewSlackAPI(token string) SlackAPI {
	return &SlackAPIImpl{
		api: slack.New(token),
	}
}

func (s *SlackAPIImpl) PostMessage(ch string, msg string, params slack.PostMessageParameters) (string, string, error) {
	return s.api.PostMessage(ch, msg, params)
}

func (s *SlackAPIImpl) CreateChannel(name string) (*slack.Channel, error) {
	return s.api.CreateChannel(name)
}

func (s *SlackAPIImpl) CreateGroup(group string) (*slack.Group, error) {
	return s.api.CreateGroup(group)
}

func (s *SlackAPIImpl) NewRTM() *slack.RTM {
	return s.api.NewRTM()
}

func (s *SlackAPIImpl) GetChannels(excludeArchived bool) ([]slack.Channel, error) {
	return s.api.GetChannels(excludeArchived)
}

func (s *SlackAPIImpl) GetGroups(excludeArchived bool) ([]slack.Group, error) {
	return s.api.GetGroups(excludeArchived)
}

func (s *SlackAPIImpl) AddPin(ch string, item slack.ItemRef) error {
	return s.api.AddPin(ch, item)
}

func (s *SlackAPIImpl) AddStar(ch string, item slack.ItemRef) error {
	return s.api.AddStar(ch, item)
}

func (s *SlackAPIImpl) AddReaction(name string, item slack.ItemRef) error {
	return s.api.AddReaction(name, item)
}

func (s *SlackAPIImpl) AuthTest() (*slack.AuthTestResponse, error) {
	return s.api.AuthTest()
}
