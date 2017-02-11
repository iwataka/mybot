package mybot

import (
	"reflect"
	"testing"

	"github.com/iwataka/anaconda"
)

func TestTwitterAction(t *testing.T) {
	a1 := TwitterAction{
		Retweet:     true,
		Favorite:    false,
		Follow:      true,
		Collections: []string{"col1", "col2"},
	}
	a2 := a1
	a3 := TwitterAction{
		Retweet:     true,
		Favorite:    true,
		Follow:      false,
		Collections: []string{"col1", "col3"},
	}

	a1.Add(&a3)
	if a1.Retweet != true {
		t.Fatalf("%v expected but %v found", false, a1.Retweet)
	}
	if a1.Favorite != true {
		t.Fatalf("%v expected but %v found", true, a1.Favorite)
	}
	if a1.Follow != true {
		t.Fatalf("%v expected but %v found", true, a1.Follow)
	}
	if len(a1.Collections) != 3 {
		t.Fatalf("%d expected but %d found", 3, len(a1.Collections))
	}

	a2.Sub(&a3)
	if a2.Retweet != false {
		t.Fatalf("%v expected but %v found", false, a2.Retweet)
	}
	if a2.Favorite != false {
		t.Fatalf("%v expected but %v found", false, a2.Favorite)
	}
	if a2.Follow != true {
		t.Fatalf("%v expected but %v found", true, a2.Follow)
	}
	if len(a2.Collections) != 1 {
		t.Fatalf("%d expected but %d found", 1, len(a2.Collections))
	}
}

func TestPostProcessorEach(t *testing.T) {
	action := &TweetAction{
		Twitter: &TwitterAction{
			Retweet:     true,
			Favorite:    false,
			Follow:      false,
			Collections: []string{"foo"},
		},
		Slack: &SlackAction{
			Channels: []string{"bar"},
		},
	}
	cache, err := NewFileCache("")
	if err != nil {
		t.Fatal(err)
	}
	tweet := anaconda.Tweet{}
	tweet.IdStr = "000"
	pp := TwitterPostProcessorEach{action, cache}

	err = pp.Process(tweet, true)
	if err != nil {
		t.Fatal(err)
	}
	ac, err := cache.GetTweetAction(tweet.Id)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ac, action) {
		t.Fatalf("%v is not cached properly", action)
	}

	action2 := &TweetAction{
		Twitter: &TwitterAction{
			Favorite: true,
		},
		Slack: NewSlackAction(),
	}
	pp2 := TwitterPostProcessorEach{action2, cache}

	err = pp2.Process(tweet, true)
	if err != nil {
		t.Fatal(err)
	}
	ac2, err := cache.GetTweetAction(tweet.Id)
	if err != nil {
		t.Fatal(err)
	}
	if !ac2.Twitter.Favorite {
		t.Fatalf("%v is not cached properly", action2)
	}
}
