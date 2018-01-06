package mybot

import (
	"fmt"
	"testing"

	"github.com/iwataka/anaconda"
	"github.com/iwataka/deep"
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

	result1 := a1.Add(a2)
	if result1.Retweet != true {
		t.Fatalf("%v expected but %v found", false, result1.Retweet)
	}
	if result1.Favorite != true {
		t.Fatalf("%v expected but %v found", true, result1.Favorite)
	}
	if len(result1.Collections) != 3 {
		t.Fatalf("%d expected but %d found", 3, len(result1.Collections))
	}

	result2 := a1.Sub(a2)
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
	action := Action{
		Twitter: TwitterAction{
			Collections: []string{"foo"},
		},
		Slack: SlackAction{
			Channels:  []string{"bar"},
			Reactions: []string{},
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
	ac := cache.GetTweetAction(tweet.Id)
	if diff := deep.Equal(ac, action); diff != nil {
		t.Fatal(diff)
	}

	action2 := Action{
		Twitter: NewTwitterAction(),
		Slack:   NewSlackAction(),
	}
	action2.Twitter.Favorite = true
	pp2 := TwitterPostProcessorEach{action2, cache}

	err = pp2.Process(tweet, true)
	if err != nil {
		t.Fatal(err)
	}
	ac2 := cache.GetTweetAction(tweet.Id)
	if !ac2.Twitter.Favorite {
		t.Fatalf("%v is not cached properly", action2)
	}
}

func TestCheckTwitterError(t *testing.T) {
	err130 := anaconda.TwitterError{Code: 130}
	testCheckTwitterError(t, err130)
	err187 := anaconda.TwitterError{Code: 187}
	testCheckTwitterError(t, err187)
	err327 := anaconda.TwitterError{Code: 327}
	testCheckTwitterError(t, err327)

	apiErr := anaconda.ApiError{}
	apiErr.StatusCode = 400
	res := anaconda.TwitterErrorResponse{}
	res.Errors = []anaconda.TwitterError{err187, err327}
	apiErr.Decoded = res
	testCheckTwitterError(t, apiErr)

	apiServerErr := anaconda.ApiError{StatusCode: 503}
	testCheckTwitterError(t, apiServerErr)
}

func testCheckTwitterError(t *testing.T, err error) {
	var msg string
	switch e := err.(type) {
	case anaconda.TwitterError:
		msg = fmt.Sprintf("Error code %d should be ignored", e.Code)
	case anaconda.ApiError:
		msg = fmt.Sprintf("API Error %d should be ignored", e.StatusCode)
	}
	if CheckTwitterError(err) {
		t.Fatal(msg)
	}
}
