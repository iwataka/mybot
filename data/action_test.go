package data_test

import (
	"github.com/go-test/deep"
	. "github.com/iwataka/mybot/data"

	"testing"
)

func TestTweetAction_AddNil(t *testing.T) {
	a1 := Action{
		Twitter: NewTwitterAction(),
		Slack:   NewSlackAction(),
	}
	a2 := Action{
		Twitter: TwitterAction{
			Collections: []string{"foo"},
		},
		Slack: SlackAction{
			Channels:  []string{"bar"},
			Reactions: []string{},
		},
	}
	a2.Twitter.Retweet = true

	result1 := a1.Add(a2)
	if diff := deep.Equal(result1, a2); diff != nil {
		t.Fatal(diff)
	}

	result2 := a2.Add(a1)
	if diff := deep.Equal(result2, a2); diff != nil {
		t.Fatal(diff)
	}
}

func TestTweetAction_Add(t *testing.T) {
	a1 := Action{
		Twitter: TwitterAction{
			Collections: []string{"twitter"},
		},
		Slack: SlackAction{
			Channels: []string{"slack"},
		},
	}
	a1.Twitter.Retweet = true
	a2 := Action{
		Twitter: TwitterAction{
			Collections: []string{"facebook"},
		},
		Slack: SlackAction{
			Channels: []string{"mattermost"},
		},
	}
	a2.Twitter.Favorite = true
	result := a1.Add(a2)

	if !result.Twitter.Retweet {
		t.Fatalf("%v expected but %v found", true, result.Twitter.Retweet)
	}
	if !result.Twitter.Favorite {
		t.Fatalf("%v expected but %v found", true, result.Twitter.Favorite)
	}
	if len(result.Twitter.Collections) != 2 {
		t.Fatalf("%d expected but %d found", 2, len(result.Twitter.Collections))
	}
	if len(result.Slack.Channels) != 2 {
		t.Fatalf("%d expected but %d found", 2, len(result.Slack.Channels))
	}
}
