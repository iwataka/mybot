package data_test

import (
	"github.com/iwataka/deep"
	"github.com/iwataka/mybot/data"
	"github.com/stretchr/testify/require"

	"testing"
)

var testAction1 = data.Action{
	Twitter: data.TwitterAction{
		Collections: []string{"foo"},
	},
	Slack: data.SlackAction{
		Channels:  []string{"bar"},
		Reactions: []string{},
	},
}

var testAction2 = data.Action{
	Twitter: data.TwitterAction{
		Collections: []string{"twitter"},
	},
	Slack: data.SlackAction{
		Channels: []string{"slack"},
	},
}

var testAction3 = data.Action{
	Twitter: data.TwitterAction{
		Collections: []string{"facebook"},
	},
	Slack: data.SlackAction{
		Channels: []string{"mattermost"},
	},
}

func init() {
	testAction1.Twitter.Retweet = true
	testAction2.Twitter.Retweet = true
	testAction3.Twitter.Favorite = true
}

func TestAction_Add_Empty(t *testing.T) {
	result1 := data.NewAction().Add(testAction1)
	require.Nil(t, deep.Equal(result1, testAction1))
	result2 := testAction1.Add(data.NewAction())
	require.Nil(t, deep.Equal(result2, testAction1))
}

func TestAction_Add(t *testing.T) {
	result := testAction2.Add(testAction3)
	require.True(t, result.Twitter.Retweet)
	require.True(t, result.Twitter.Favorite)
	require.Len(t, result.Twitter.Collections, 2)
	require.Len(t, result.Slack.Channels, 2)
}

func TestAction_Sub_Empty(t *testing.T) {
	result1 := data.NewAction().Sub(testAction1)
	require.Nil(t, deep.Equal(result1, data.NewAction()))
	result2 := testAction1.Sub(data.NewAction())
	require.Nil(t, deep.Equal(result2, testAction1))
}

func TestAction_Sub(t *testing.T) {
	result := testAction2.Sub(testAction3)
	deep.IgnoreDifferenceBetweenEmptySliceAndNil = true
	require.Nil(t, deep.Equal(result, testAction2))
	deep.IgnoreDifferenceBetweenEmptySliceAndNil = false
}

func Test_NewAction_ReturnsEmpty(t *testing.T) {
	require.True(t, data.NewAction().IsEmpty())
}
