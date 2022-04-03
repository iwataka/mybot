package models

import (
	"github.com/iwataka/slack"
)

type SlackAPI interface {
	PostMessage(ch, msg string, params slack.PostMessageParameters) error
	CreateChannel(name string) error
	CreateGroup(group string) error
	NewRTM() *slack.RTM
	GetChannels(excludeArchived bool) ([]Channel, error)
	GetGroups(excludeArchived bool) ([]Group, error)
	AddPin(ch, timestamp string) error
	AddStar(ch, timestamp string) error
	AddReaction(ch, timestamp, name string) error
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

func (s *SlackAPIImpl) PostMessage(ch string, msg string, params slack.PostMessageParameters) error {
	_, _, err := s.api.PostMessage(ch, msg, params)
	return err
}

func (s *SlackAPIImpl) CreateChannel(name string) error {
	_, err := s.api.CreateChannel(name)
	return err
}

func (s *SlackAPIImpl) CreateGroup(group string) error {
	_, err := s.api.CreateGroup(group)
	return err
}

func (s *SlackAPIImpl) NewRTM() *slack.RTM {
	return s.api.NewRTM()
}

func (s *SlackAPIImpl) GetChannels(excludeArchived bool) ([]Channel, error) {
	channels, err := s.api.GetChannels(excludeArchived)
	if err != nil {
		return nil, err
	}

	chs := make([]Channel, len(channels))
	for i, ch := range channels {
		chs[i] = Channel{
			ID:   ch.ID,
			Name: ch.Name,
		}
	}
	return chs, nil
}

type Channel struct {
	ID   string `json:"id" toml:"id" bson:"id" yaml:"id"`
	Name string `json:"name" toml:"name" bson:"name" yaml:"name"`
}

func (s *SlackAPIImpl) GetGroups(excludeArchived bool) ([]Group, error) {
	groups, err := s.api.GetGroups(excludeArchived)
	if err != nil {
		return nil, err
	}

	grps := make([]Group, len(groups))
	for i, grp := range groups {
		grps[i] = Group{
			ID:   grp.ID,
			Name: grp.Name,
		}
	}
	return grps, nil
}

type Group struct {
	ID   string `json:"id" toml:"id" bson:"id" yaml:"id"`
	Name string `json:"name" toml:"name" bson:"name" yaml:"name"`
}

func (s *SlackAPIImpl) AddPin(ch, timestamp string) error {
	item := slack.NewRefToMessage(ch, timestamp)
	return s.api.AddPin(ch, item)
}

func (s *SlackAPIImpl) AddStar(ch, timestamp string) error {
	item := slack.NewRefToMessage(ch, timestamp)
	return s.api.AddStar(ch, item)
}

func (s *SlackAPIImpl) AddReaction(ch, timestamp, name string) error {
	item := slack.NewRefToMessage(ch, timestamp)
	return s.api.AddReaction(name, item)
}

func (s *SlackAPIImpl) AuthTest() (*slack.AuthTestResponse, error) {
	return s.api.AuthTest()
}
