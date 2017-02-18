package mybot

import (
	"reflect"
	"testing"

	"github.com/iwataka/anaconda"
)

func TestTwitterAction(t *testing.T) {
	a1 := TwitterAction{
		Collections: []string{"col1", "col2"},
	}
	a1.Retweet = true
	a2 := TwitterAction{
		Collections: []string{"col1", "col3"},
	}
	a2.Retweet = true
	a2.Favorite = true

	result1 := a1.Add(&a2)
	if result1.Retweet != true {
		t.Fatalf("%v expected but %v found", false, result1.Retweet)
	}
	if result1.Favorite != true {
		t.Fatalf("%v expected but %v found", true, result1.Favorite)
	}
	if len(result1.Collections) != 3 {
		t.Fatalf("%d expected but %d found", 3, len(result1.Collections))
	}

	result2 := a1.Sub(&a2)
	if result2.Retweet != false {
		t.Fatalf("%v expected but %v found", false, result2.Retweet)
	}
	if result2.Favorite != false {
		t.Fatalf("%v expected but %v found", false, result2.Favorite)
	}
	if len(result2.Collections) != 1 {
		t.Fatalf("%d expected but %d found", 1, len(result2.Collections))
	}
}

func TestPostProcessorEach(t *testing.T) {
	action := &Action{
		Twitter: &TwitterAction{
			Collections: []string{"foo"},
		},
		Slack: &SlackAction{
			Channels: []string{"bar"},
		},
	}
	action.Twitter.Retweet = true
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

	action2 := &Action{
		Twitter: NewTwitterAction(),
		Slack:   NewSlackAction(),
	}
	action2.Twitter.Favorite = true
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
