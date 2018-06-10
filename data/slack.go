package data

import (
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/utils"
)

type SlackAction struct {
	models.SlackActionProperties
	Reactions []string `json:"reactions,omitempty" toml:"reactions,omitempty" bson:"reactions,omitempty"`
	Channels  []string `json:"channels,omitempty" toml:"channels,omitempty" bson:"channels,omitempty"`
}

func NewSlackAction() SlackAction {
	return SlackAction{
		Channels:  []string{},
		Reactions: []string{},
	}
}

func (a SlackAction) Add(action SlackAction) SlackAction {
	return a.op(action, true)
}

func (a SlackAction) Sub(action SlackAction) SlackAction {
	return a.op(action, false)
}

func (a SlackAction) op(action SlackAction, add bool) SlackAction {
	result := a

	result.Pin = utils.CalcBools(a.Pin, action.Pin, add)
	result.Star = utils.CalcBools(a.Star, action.Star, add)
	result.Reactions = utils.CalcStringSlices(a.Reactions, action.Reactions, add)
	result.Channels = utils.CalcStringSlices(a.Channels, action.Channels, add)

	return result

}

func (a SlackAction) IsEmpty() bool {
	return !a.Pin &&
		!a.Star &&
		len(a.Channels) == 0 &&
		len(a.Reactions) == 0
}
