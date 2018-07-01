package mybot

import (
	"container/list"
	"errors"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/mocks"
)

func TestSlackConvertFromTweet(t *testing.T) {
	tweet := anaconda.Tweet{
		IdStr: "1",
		User: anaconda.User{
			IdStr: "1",
		},
	}
	text, params := convertFromTweetToSlackMsg(tweet)
	if text != TwitterStatusURL(tweet) {
		t.Fatal("Text is invalid")
	}
	if !params.UnfurlLinks || !params.UnfurlMedia {
		t.Fatal("Should unfurl all kinds of things")
	}
}

func TestNewSlackAPI(t *testing.T) {
	api := NewSlackAPIWithAuth("", nil, nil)
	if api.Enabled() {
		t.Fatalf("%v is expected to be disabled but not", api)
	}
}

func TestSlackAPIDequeueMsg(t *testing.T) {
	api := NewSlackAPIWithAuth("", nil, nil)
	msg := api.dequeueMsg("channel")
	if msg != nil {
		t.Fatalf("%s expected but %s found", nil, msg)
	}
}

func TestSlackAPIEnqueueMsg(t *testing.T) {
	api := NewSlackAPIWithAuth("", nil, nil)
	ch := "channel"
	msg := &SlackMsg{"text", nil}
	api.enqueueMsg(ch, msg.text, msg.params)
	m := api.dequeueMsg(ch)
	if !reflect.DeepEqual(msg, m) {
		t.Fatalf("%s expected but %s found", msg, m)
	}
	if api.dequeueMsg(ch) != nil {
		t.Fatal("dequeueMsg not working properly")
	}
}

func TestSlackAPIPostMessage(t *testing.T) {
	testSlackAPIPostMessage(t, true)
}

func TestSlackAPIPostMessageWithoutQueue(t *testing.T) {
	testSlackAPIPostMessage(t, false)
}

func testSlackAPIPostMessage(t *testing.T, queue bool) {
	ctrl := gomock.NewController(t)
	slackAPIMock := mocks.NewMockSlackAPI(ctrl)
	slackAPI := SlackAPI{api: slackAPIMock, msgQueue: make(map[string]*list.List)}

	slackAPIMock.EXPECT().PostMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return("", "", errors.New("channel_not_found"))
	slackAPIMock.EXPECT().CreateChannel(gomock.Any()).Return(nil, errors.New("user_is_bot"))
	slackAPIMock.EXPECT().PostMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return("", "", nil)

	ch := "channel"
	text := "text"
	var msg *SlackMsg
	if queue {
		msg = &SlackMsg{text, nil}
	}
	slackAPI.PostMessage(ch, text, nil, queue)
	m := slackAPI.dequeueMsg(ch)

	if !reflect.DeepEqual(msg, m) {
		t.Fatalf("%s expected but %s found", msg, m)
	}
}

func TestSlackAPISendMsgQueues(t *testing.T) {
	ctrl := gomock.NewController(t)
	slackAPIMock := mocks.NewMockSlackAPI(ctrl)
	slackAPI := SlackAPI{api: slackAPIMock, msgQueue: make(map[string]*list.List)}

	ch := "channel"

	err := slackAPI.sendMsgQueues(ch)
	if err != nil {
		t.Fatal(err)
	}

	slackAPIMock.EXPECT().PostMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return("", "", nil)

	text := "text"
	slackAPI.enqueueMsg(ch, text, nil)

	err = slackAPI.sendMsgQueues(ch)
	if err != nil {
		t.Fatal(err)
	}
	err = slackAPI.sendMsgQueues(ch)
	if err != nil {
		t.Fatal(err)
	}
}
