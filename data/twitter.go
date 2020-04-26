package data

import (
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/utils"
)

// TwitterAction can indicate for various actions for Twitter's tweets.
type TwitterAction struct {
	models.TwitterActionProperties `yaml:",inline"`
	Collections                    []string `json:"collections" toml:"collections" bson:"collections" yaml:"collections"`
}

func NewTwitterAction() TwitterAction {
	return TwitterAction{
		Collections: []string{},
	}
}

func (a TwitterAction) Add(action TwitterAction) TwitterAction {
	return a.op(action, true)
}

func (a TwitterAction) Sub(action TwitterAction) TwitterAction {
	return a.op(action, false)
}

func (a TwitterAction) op(action TwitterAction, add bool) TwitterAction {
	result := a

	result.Tweet = utils.CalcBools(a.Tweet, action.Tweet, add)
	result.Retweet = utils.CalcBools(a.Retweet, action.Retweet, add)
	result.Favorite = utils.CalcBools(a.Favorite, action.Favorite, add)
	result.Collections = utils.CalcStringSlices(a.Collections, action.Collections, add)

	return result
}

func (a TwitterAction) IsEmpty() bool {
	return !a.Tweet &&
		!a.Retweet &&
		!a.Favorite &&
		len(a.Collections) == 0
}
