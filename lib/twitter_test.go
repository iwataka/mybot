package mybot

import (
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/deep"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/mocks"
	"github.com/stretchr/testify/require"
)

func TestTwitterPostProcessorEach(t *testing.T) {
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
	require.NoError(t, err)
	tweet := anaconda.Tweet{}
	tweet.IdStr = "000"
	pp := TwitterPostProcessorEach{action, cache}

	err = pp.Process(tweet, true)
	require.NoError(t, err)
	ac := cache.GetTweetAction(tweet.Id)
	require.Nil(t, deep.Equal(ac, action))

	action2 := data.Action{
		Twitter: data.NewTwitterAction(),
		Slack:   data.NewSlackAction(),
	}
	action2.Twitter.Favorite = true
	pp2 := TwitterPostProcessorEach{action2, cache}

	err = pp2.Process(tweet, true)
	require.NoError(t, err)
	ac2 := cache.GetTweetAction(tweet.Id)
	require.True(t, ac2.Twitter.Favorite)
}

func Test_CheckTwitterError(t *testing.T) {
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
	require.False(t, CheckTwitterError(err), msg)
}

func TestTwitterAPI_NotifyToAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config, err := NewFileConfig("testdata/config.template.toml")
	require.NoError(t, err)

	twitterAPIMock := mocks.NewMockTwitterAPI(ctrl)
	twitterAPIMock.EXPECT().PostDMToScreenName(gomock.Any(), gomock.Any()).Return(anaconda.DirectMessage{}, nil)
	twitterAPIMock.EXPECT().GetSelf(gomock.Any()).Return(anaconda.User{Name: "user"}, nil)
	twitterAPIMock.EXPECT().PostDMToScreenName(gomock.Any(), gomock.Any()).Return(anaconda.DirectMessage{}, nil)
	twitterAPI := NewTwitterAPI(twitterAPIMock, config, nil)

	slackAPIMock := mocks.NewMockSlackAPI(ctrl)
	slackAPIMock.EXPECT().PostMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return("", "", nil)
	slackAPI := &SlackAPI{api: slackAPIMock, config: config}

	tweet := &anaconda.Tweet{
		Coordinates: &anaconda.Coordinates{Type: "Point"},
		Place:       anaconda.Place{Country: "japan"},
	}
	require.NoError(t, twitterAPI.NotifyToAll(slackAPI, tweet))
}
