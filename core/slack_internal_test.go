package core

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/mocks"
	"github.com/stretchr/testify/require"
)

func Test_convertFromTweetToSlackMsg(t *testing.T) {
	tweet := anaconda.Tweet{
		IdStr: "1",
		User: anaconda.User{
			IdStr: "1",
		},
	}
	text, params := convertFromTweetToSlackMsg(tweet)

	require.Equal(t, TwitterStatusURL(tweet), text)
	require.True(t, params.UnfurlLinks)
	require.True(t, params.UnfurlMedia)
}

func Test_NewSlackAPIWithAuth(t *testing.T) {
	api := NewSlackAPIWithAuth("", nil, nil)
	require.False(t, api.Enabled())
}

func TestSlackAPI_dequeueMsg(t *testing.T) {
	api := NewSlackAPIWithAuth("", nil, nil)
	msg := api.dequeueMsg("channel")
	require.Nil(t, msg)
}

func TestSlackAPI_enqueueMsg(t *testing.T) {
	api := NewSlackAPIWithAuth("", nil, nil)
	ch := "channel"
	msg := &SlackMsg{"text", nil}
	api.enqueueMsg(ch, msg.text, msg.params)
	m := api.dequeueMsg(ch)

	require.Equal(t, msg, m)
	require.Nil(t, api.dequeueMsg(ch))
}

func TestSlackAPI_PostMessage(t *testing.T) {
	testSlackAPIPostMessage(t, true)
}

func TestSlackAPI_PostMessage_WithPrivateChannel(t *testing.T) {
	testSlackAPIPostMessage(t, false)
}

func testSlackAPIPostMessage(t *testing.T, channelIsOpen bool) {
	ctrl := gomock.NewController(t)
	slackAPIMock := mocks.NewMockSlackAPI(ctrl)
	slackAPI := SlackAPI{api: slackAPIMock, msgQueue: make(map[string]*concurrentQueue)}

	slackAPIMock.EXPECT().PostMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("channel_not_found"))
	if channelIsOpen {
		slackAPIMock.EXPECT().CreateChannel(gomock.Any()).Return(errors.New("user_is_bot"))
	} else {
		slackAPIMock.EXPECT().CreateGroup(gomock.Any()).Return(errors.New("user_is_bot"))
	}
	slackAPIMock.EXPECT().PostMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	ch := "channel"
	text := "text"
	msg := &SlackMsg{text, nil}
	err := slackAPI.PostMessage(ch, text, nil, channelIsOpen)
	require.NoError(t, err)
	m := slackAPI.dequeueMsg(ch)

	require.Equal(t, msg, m)
}

func TestSlackAPI_sendMsgQueues(t *testing.T) {
	ctrl := gomock.NewController(t)
	slackAPIMock := mocks.NewMockSlackAPI(ctrl)
	slackAPI := SlackAPI{api: slackAPIMock, msgQueue: make(map[string]*concurrentQueue)}
	ch := "channel"

	err := slackAPI.sendMsgQueues(ch)
	require.NoError(t, err)

	slackAPIMock.EXPECT().PostMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	text := "text"
	slackAPI.enqueueMsg(ch, text, nil)

	err = slackAPI.sendMsgQueues(ch)
	require.NoError(t, err)
	err = slackAPI.sendMsgQueues(ch)
	require.NoError(t, err)
}
