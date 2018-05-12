package runner_test

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/mybot/lib"
	"github.com/iwataka/mybot/mocks"
	. "github.com/iwataka/mybot/runner"

	"errors"
	"testing"
)

func TestTwitterAPIIsAvailable(t *testing.T) {
	var ctrl *gomock.Controller
	var twitterAPIMock *mocks.MockTwitterAPI
	var twitterAPI *mybot.TwitterAPI

	if err := TwitterAPIIsAvailable(nil); err == nil {
		t.Fatalf("Error expected but nothing occurred")
	}

	ctrl = gomock.NewController(t)
	twitterAPIMock = mocks.NewMockTwitterAPI(ctrl)
	twitterAPIMock.EXPECT().VerifyCredentials().Return(true, nil)
	twitterAPI = &mybot.TwitterAPI{API: twitterAPIMock}
	if err := TwitterAPIIsAvailable(twitterAPI); err != nil {
		t.Fatalf("Error not expected but occurred: %s", err.Error())
	}

	if err := TwitterAPIIsAvailable(&mybot.TwitterAPI{}); err == nil {
		t.Fatalf("Error expected but nothing occurred")
	}

	twitterAPIMock = mocks.NewMockTwitterAPI(ctrl)
	twitterAPIMock.EXPECT().VerifyCredentials().Return(false, nil)
	twitterAPI = &mybot.TwitterAPI{API: twitterAPIMock}
	if err := TwitterAPIIsAvailable(twitterAPI); err == nil {
		t.Fatalf("Error expected but nothing occurred")
	}

	twitterAPIMock = mocks.NewMockTwitterAPI(ctrl)
	twitterAPIMock.EXPECT().VerifyCredentials().Return(false, nil)
	twitterAPI = &mybot.TwitterAPI{API: twitterAPIMock}
	if err := TwitterAPIIsAvailable(twitterAPI); err == nil {
		t.Fatalf("Error expected but nothing occurred")
	}

	twitterAPIMock = mocks.NewMockTwitterAPI(ctrl)
	twitterAPIMock.EXPECT().VerifyCredentials().Return(false, errors.New(""))
	twitterAPI = &mybot.TwitterAPI{API: twitterAPIMock}
	if err := TwitterAPIIsAvailable(twitterAPI); err == nil {
		t.Fatalf("Error expected but nothing occurred")
	}
}
