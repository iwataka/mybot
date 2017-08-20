package models

import (
	"github.com/jinzhu/gorm"
	"github.com/nlopes/slack"
)

type SlackAPI interface {
	PostMessage(ch, msg string, params slack.PostMessageParameters) (string, string, error)
	CreateChannel(name string) (*slack.Channel, error)
	NewRTM() *slack.RTM
	GetChannels(excludeArchived bool) ([]slack.Channel, error)
	AddPin(ch string, item slack.ItemRef) error
	AddStar(ch string, item slack.ItemRef) error
	AddReaction(name string, item slack.ItemRef) error
}

type SlackAction struct {
	gorm.Model
	SlackActionProperties
	Reactions []SlackReaction
	Channels  []SlackChannel
}

type SlackChannel struct {
	gorm.Model
	SlackActionID uint
	Name          string
}

func (a *SlackAction) GetChannels() []string {
	result := []string{}
	for _, c := range a.Channels {
		result = append(result, c.Name)
	}
	return result
}

func (a *SlackAction) SetChannels(chs []string) {
	a.Channels = []SlackChannel{}
	for _, ch := range chs {
		c := SlackChannel{
			SlackActionID: a.ID,
			Name:          ch,
		}
		a.Channels = append(a.Channels, c)
	}
}

type SlackReaction struct {
	gorm.Model
	SlackActionID uint
	Text          string
}

func (a *SlackAction) GetReactions() []string {
	result := []string{}
	for _, c := range a.Reactions {
		result = append(result, c.Text)
	}
	return result
}

func (a *SlackAction) SetReactions(rs []string) {
	a.Reactions = []SlackReaction{}
	for _, text := range rs {
		r := SlackReaction{
			SlackActionID: a.ID,
			Text:          text,
		}
		a.Reactions = append(a.Reactions, r)
	}
}

type SlackActionProperties struct {
	Pin  bool `json:"pin" toml:"pin" bson:"pin"`
	Star bool `json:"star" toml:"star" bson:"star"`
}
