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

func TestSlackActionAdd(t *testing.T) {
	a1 := &SlackAction{
		Channels: []string{"foo", "bar"},
	}
	a1.Pin = true
	a1.Reactions = []string{":smile:"}
	a2 := &SlackAction{
		Channels: []string{"foo", "fizz"},
	}
	a2.Reactions = []string{":smile:", ":plane:"}
	result := a1.Add(a2)

	if !result.Pin {
		t.Fatalf("%v expected but %v found", true, result.Pin)
	}
	if result.Star {
		t.Fatalf("%v expected but %v found", false, result.Star)
	}
	if len(result.Reactions) != 2 {
		t.Fatalf("%v expected but %v found", 2, len(result.Reactions))
	}
	if len(result.Channels) != 3 {
		t.Fatalf("%v expected but %v found", 3, len(result.Channels))
	}
}

func TestSlackActionSub(t *testing.T) {
	a1 := &SlackAction{
		Channels: []string{"foo", "bar"},
	}
	a2 := &SlackAction{
		Channels: []string{"foo", "fizz"},
	}
	result := a1.Sub(a2)
	if len(result.Channels) != 1 {
		t.Fatalf("%v expected but %v found", 1, len(result.Channels))
	}
}

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
	api := NewSlackAPI("", nil, nil)
	if api.Enabled() {
		t.Fatalf("%v is expected to be disabled but not", api)
	}
}

func TestSlackActionIsEmpty(t *testing.T) {
	a := NewSlackAction()
	if !a.IsEmpty() {
		t.Fatalf("%v should be empty", a)
	}
	a.Channels = []string{"foo"}
	if a.IsEmpty() {
		t.Fatalf("%v should not be empty", a)
	}
}

func TestSlackAPIDequeueMsg(t *testing.T) {
	api := NewSlackAPI("", nil, nil)
	msg := api.dequeueMsg("channel")
	if msg != nil {
		t.Fatalf("%s expected but %s found", nil, msg)
	}
}

func TestSlackAPIEnqueueMsg(t *testing.T) {
	api := NewSlackAPI("", nil, nil)
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
