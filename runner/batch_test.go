package runner_test

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/core"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/mocks"
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/runner"
	"github.com/stretchr/testify/require"

	"errors"
	"testing"
)

func Test_TwitterAPIIsAvailable(t *testing.T) {
	ctrl := gomock.NewController(t)
	var twitterAPIMock *mocks.MockTwitterAPI
	var twitterAPI *core.TwitterAPI

	require.Error(t, runner.TwitterAPIIsAvailable(nil))

	twitterAPI = generateVerifiedTwitterAPI(t)
	require.NoError(t, runner.TwitterAPIIsAvailable(twitterAPI))

	require.Error(t, runner.TwitterAPIIsAvailable(&core.TwitterAPI{}))

	twitterAPIMock = mocks.NewMockTwitterAPI(ctrl)
	twitterAPIMock.EXPECT().VerifyCredentials().Return(false, nil)
	twitterAPI = core.NewTwitterAPI(twitterAPIMock, nil, nil)
	require.Error(t, runner.TwitterAPIIsAvailable(twitterAPI))

	twitterAPIMock = mocks.NewMockTwitterAPI(ctrl)
	twitterAPIMock.EXPECT().VerifyCredentials().Return(false, errors.New(""))
	twitterAPI = core.NewTwitterAPI(twitterAPIMock, nil, nil)
	require.Error(t, runner.TwitterAPIIsAvailable(twitterAPI))
}

func TestBatchRunnerUsedWithStream_IsAvailable(t *testing.T) {
	twitterAPI := generateVerifiedTwitterAPI(t)
	r := runner.NewBatchRunnerUsedWithStream(twitterAPI, nil, nil, nil, nil)
	require.NoError(t, r.IsAvailable())
}

func TestBatchRunnerUsedWithStream_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	r := generateBatchRunnerUsedWithStream(t, ctrl)
	processedTweets, processedActions, err := r.Run()
	require.Len(t, processedTweets, 2)
	require.Len(t, processedActions, 2)
	require.NoError(t, err)
}

func generateBatchRunnerUsedWithStream(t *testing.T, ctrl *gomock.Controller) *runner.BatchRunnerUsedWithStream {
	twitterAPIMock := mocks.NewMockTwitterAPI(ctrl)
	slackAPIMock := mocks.NewMockSlackAPI(ctrl)
	registerProcessSearch(twitterAPIMock, slackAPIMock)
	return generateBaseRunner(t, twitterAPIMock, slackAPIMock)
}

func generateBaseRunner(t *testing.T, twitterAPIMock models.TwitterAPI, slackAPIMock models.SlackAPI) *runner.BatchRunnerUsedWithStream {
	cache, err := data.NewFileCache("")
	require.NoError(t, err)
	config, err := core.NewFileConfig("../core/testdata/config.yaml")
	require.NoError(t, err)
	twitterAPI := core.NewTwitterAPI(twitterAPIMock, config, cache)
	slackAPI := core.NewSlackAPI(slackAPIMock, config, cache)
	var visionAPI core.VisionMatcher
	var languageAPI core.LanguageMatcher
	return runner.NewBatchRunnerUsedWithStream(twitterAPI, slackAPI, visionAPI, languageAPI, config)
}

func registerProcessSearch(twitterAPIMock *mocks.MockTwitterAPI, slackAPIMock *mocks.MockSlackAPI) {
	fooSearchRes := anaconda.SearchResponse{
		Statuses: []anaconda.Tweet{
			{
				Text:         "foo",
				RetweetCount: 100,
				Id:           0,
			},
		},
	}
	barSearchRes := anaconda.SearchResponse{
		Statuses: []anaconda.Tweet{
			{
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

		twitterAPIMock.EXPECT().GetFavorites(gomock.Any()).Return([]anaconda.Tweet{}, nil),
	)
}

func generateVerifiedTwitterAPI(t *testing.T) *core.TwitterAPI {
	ctrl := gomock.NewController(t)
	twitterAPIMock := mocks.NewMockTwitterAPI(ctrl)
	twitterAPIMock.EXPECT().VerifyCredentials().Return(true, nil)
	return core.NewTwitterAPI(twitterAPIMock, nil, nil)
}
