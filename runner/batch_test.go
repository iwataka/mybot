package runner_test

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/lib"
	"github.com/iwataka/mybot/mocks"
	"github.com/iwataka/mybot/models"
	. "github.com/iwataka/mybot/runner"
	"github.com/stretchr/testify/require"

	"errors"
	"testing"
)

func Test_TwitterAPIIsAvailablt(t *testing.T) {
	ctrl := gomock.NewController(t)
	var twitterAPIMock *mocks.MockTwitterAPI
	var twitterAPI *mybot.TwitterAPI

	require.Error(t, TwitterAPIIsAvailable(nil))

	twitterAPI = generateVerifiedTwitterAPI(t)
	require.NoError(t, TwitterAPIIsAvailable(twitterAPI))

	require.Error(t, TwitterAPIIsAvailable(&mybot.TwitterAPI{}))

	twitterAPIMock = mocks.NewMockTwitterAPI(ctrl)
	twitterAPIMock.EXPECT().VerifyCredentials().Return(false, nil)
	twitterAPI = mybot.NewTwitterAPI(twitterAPIMock, nil, nil)
	require.Error(t, TwitterAPIIsAvailable(twitterAPI))

	twitterAPIMock = mocks.NewMockTwitterAPI(ctrl)
	twitterAPIMock.EXPECT().VerifyCredentials().Return(false, errors.New(""))
	twitterAPI = mybot.NewTwitterAPI(twitterAPIMock, nil, nil)
	require.Error(t, TwitterAPIIsAvailable(twitterAPI))
}

func TestBatchRunnerUsedWithStream_IsAvailable(t *testing.T) {
	twitterAPI := generateVerifiedTwitterAPI(t)
	r := NewBatchRunnerUsedWithStream(twitterAPI, nil, nil, nil, nil)
	require.NoError(t, r.IsAvailable())
}

func TestBatchRunnerUsedWithoutStream_IsAvailable(t *testing.T) {
	twitterAPI := generateVerifiedTwitterAPI(t)
	baseRunner := NewBatchRunnerUsedWithStream(twitterAPI, nil, nil, nil, nil)
	r := NewBatchRunnerUsedWithoutStream(baseRunner)
	require.NoError(t, r.IsAvailable())
}

func TestBatchRunnerUsedWithStream_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	r := generateBatchRunnerUsedWithStream(t, ctrl)
	require.NoError(t, r.Run())
}

func TestBatchRunnerUsedWithoutStream_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	r := generateBatchRunnerUsedWithoutStream(t, ctrl)
	require.NoError(t, r.Run())
}

func generateBatchRunnerUsedWithStream(t *testing.T, ctrl *gomock.Controller) *BatchRunnerUsedWithStream {
	twitterAPIMock := mocks.NewMockTwitterAPI(ctrl)
	slackAPIMock := mocks.NewMockSlackAPI(ctrl)
	registerProcessSearch(twitterAPIMock, slackAPIMock)
	return generateBaseRunner(t, twitterAPIMock, slackAPIMock)
}

func generateBatchRunnerUsedWithoutStream(t *testing.T, ctrl *gomock.Controller) *BatchRunnerUsedWithoutStream {
	twitterAPIMock := mocks.NewMockTwitterAPI(ctrl)
	slackAPIMock := mocks.NewMockSlackAPI(ctrl)
	registerProcessSearch(twitterAPIMock, slackAPIMock)
	registerProcessTimeline(twitterAPIMock, slackAPIMock)
	baseRunner := generateBaseRunner(t, twitterAPIMock, slackAPIMock)
	return NewBatchRunnerUsedWithoutStream(baseRunner)
}

func generateBaseRunner(t *testing.T, twitterAPIMock models.TwitterAPI, slackAPIMock models.SlackAPI) *BatchRunnerUsedWithStream {
	cache, err := data.NewFileCache("")
	require.NoError(t, err)
	config, err := mybot.NewFileConfig("../lib/testdata/config.template.toml")
	require.NoError(t, err)
	twitterAPI := mybot.NewTwitterAPI(twitterAPIMock, config, cache)
	slackAPI := mybot.NewSlackAPI(slackAPIMock, config, cache)
	var visionAPI mybot.VisionMatcher
	var languageAPI mybot.LanguageMatcher
	return NewBatchRunnerUsedWithStream(twitterAPI, slackAPI, visionAPI, languageAPI, config)
}

func registerProcessSearch(twitterAPIMock *mocks.MockTwitterAPI, slackAPIMock *mocks.MockSlackAPI) {
	fooSearchRes := anaconda.SearchResponse{
		Statuses: []anaconda.Tweet{
			anaconda.Tweet{
				Text:         "foo",
				RetweetCount: 100,
				Id:           0,
			},
		},
	}
	barSearchRes := anaconda.SearchResponse{
		Statuses: []anaconda.Tweet{
			anaconda.Tweet{
				Text:         "foo bar",
				RetweetCount: 100,
				Id:           1,
			},
		},
	}

	gomock.InOrder(
		twitterAPIMock.EXPECT().GetSearch(gomock.Any(), gomock.Any()).Return(fooSearchRes, nil),
		twitterAPIMock.EXPECT().Retweet(gomock.Any(), gomock.Any()).Return(fooSearchRes.Statuses[0], nil),
		twitterAPIMock.EXPECT().GetSelf(gomock.Any()).Return(anaconda.User{}, nil),
		twitterAPIMock.EXPECT().GetCollectionListByUserId(gomock.Any(), gomock.Any()).Return(anaconda.CollectionListResult{}, nil),
		twitterAPIMock.EXPECT().CreateCollection(gomock.Any(), gomock.Any()).Return(anaconda.CollectionShowResult{}, nil),
		twitterAPIMock.EXPECT().AddEntryToCollection(gomock.Any(), gomock.Any(), gomock.Any()).Return(anaconda.CollectionEntryAddResult{}, nil),
		slackAPIMock.EXPECT().PostMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return("", "", nil),

		twitterAPIMock.EXPECT().GetSearch(gomock.Any(), gomock.Any()).Return(barSearchRes, nil),
		twitterAPIMock.EXPECT().Retweet(gomock.Any(), gomock.Any()).Return(barSearchRes.Statuses[0], nil),
		twitterAPIMock.EXPECT().GetCollectionListByUserId(gomock.Any(), gomock.Any()).Return(anaconda.CollectionListResult{}, nil),
		twitterAPIMock.EXPECT().CreateCollection(gomock.Any(), gomock.Any()).Return(anaconda.CollectionShowResult{}, nil),
		twitterAPIMock.EXPECT().AddEntryToCollection(gomock.Any(), gomock.Any(), gomock.Any()).Return(anaconda.CollectionEntryAddResult{}, nil),
		slackAPIMock.EXPECT().PostMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return("", "", nil),
	)
}

func registerProcessTimeline(twitterAPIMock *mocks.MockTwitterAPI, slackAPIMock *mocks.MockSlackAPI) {
	tweet := anaconda.Tweet{
		Text: "foo",
		Id:   2,
	}

	gomock.InOrder(
		twitterAPIMock.EXPECT().GetUserTimeline(gomock.Any()).Return([]anaconda.Tweet{tweet}, nil),
		twitterAPIMock.EXPECT().GetUserTimeline(gomock.Any()).Return([]anaconda.Tweet{tweet}, nil),
		twitterAPIMock.EXPECT().Retweet(gomock.Any(), gomock.Any()).Return(tweet, nil),
	)
}

func generateVerifiedTwitterAPI(t *testing.T) *mybot.TwitterAPI {
	ctrl := gomock.NewController(t)
	twitterAPIMock := mocks.NewMockTwitterAPI(ctrl)
	twitterAPIMock.EXPECT().VerifyCredentials().Return(true, nil)
	return mybot.NewTwitterAPI(twitterAPIMock, nil, nil)
}
