package data_test

import (
	"github.com/iwataka/deep"
	. "github.com/iwataka/mybot/data"
	"github.com/stretchr/testify/assert"

	"testing"
)

var testAction1 = Action{
	Twitter: TwitterAction{
		Collections: []string{"foo"},
	},
	Slack: SlackAction{
		Channels:  []string{"bar"},
		Reactions: []string{},
	},
}

var testAction2 = Action{
	Twitter: TwitterAction{
		Collections: []string{"twitter"},
	},
	Slack: SlackAction{
		Channels: []string{"slack"},
	},
}

var testAction3 = Action{
	Twitter: TwitterAction{
		Collections: []string{"facebook"},
	},
	Slack: SlackAction{
		Channels: []string{"mattermost"},
	},
}

func init() {
	testAction1.Twitter.Retweet = true
	testAction2.Twitter.Retweet = true
	testAction3.Twitter.Favorite = true
}

func TestAction_AddEmpty(t *testing.T) {
	result1 := NewAction().Add(testAction1)
	assert.Nil(t, deep.Equal(result1, testAction1))
	result2 := testAction1.Add(NewAction())
	assert.Nil(t, deep.Equal(result2, testAction1))
}

func TestAction_Add(t *testing.T) {
	result := testAction2.Add(testAction3)
	assert.True(t, result.Twitter.Retweet)
	assert.True(t, result.Twitter.Favorite)
	assert.Len(t, result.Twitter.Collections, 2)
	assert.Len(t, result.Slack.Channels, 2)
}

func TestAction_SubEmpty(t *testing.T) {
	result1 := NewAction().Sub(testAction1)
	assert.Nil(t, deep.Equal(result1, NewAction()))
	result2 := testAction1.Sub(NewAction())
	assert.Nil(t, deep.Equal(result2, testAction1))
}

func TestAction_Sub(t *testing.T) {
	result := testAction2.Sub(testAction3)
	deep.IgnoreDifferenceBetweenEmptySliceAndNil = true
	assert.Nil(t, deep.Equal(result, testAction2))
	deep.IgnoreDifferenceBetweenEmptySliceAndNil = false
}

func TestNewAction_ReturnsEmpty(t *testing.T) {
	assert.True(t, NewAction().IsEmpty())
}
