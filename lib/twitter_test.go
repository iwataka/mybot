package mybot

import (
	"fmt"
	"testing"

	"github.com/iwataka/anaconda"
	"github.com/iwataka/deep"
	"github.com/iwataka/mybot/data"
)

func TestPostProcessorEach(t *testing.T) {
	action := data.Action{
		Twitter: data.TwitterAction{
			Collections: []string{"foo"},
		},
		Slack: data.SlackAction{
			Channels:  []string{"bar"},
			Reactions: []string{},
		},
	}
	action.Twitter.Retweet = true
	cache, err := data.NewFileCache("")
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

	action2 := data.Action{
		Twitter: data.NewTwitterAction(),
		Slack:   data.NewSlackAction(),
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
	err131 := anaconda.TwitterError{Code: 131}
	testCheckTwitterError(t, err131)
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
