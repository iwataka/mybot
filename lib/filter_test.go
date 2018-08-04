package mybot

import (
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/mocks"
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/utils"
	"github.com/iwataka/slack"
	"github.com/stretchr/testify/require"
)

const (
	invalidRegexpPattern = `\K`
)

func TestFilter_CheckTweet_PatternsMatched(t *testing.T) {
	tweet := anaconda.Tweet{
		Text: "foo is bar",
	}
	filter := &Filter{
		Patterns: []string{"foo"},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckTweet(tweet, nil, nil, cache)
	require.NoError(t, err)
	require.True(t, match)
}

func TestFilter_CheckTweet_PatternsInvalid(t *testing.T) {
	tweet := anaconda.Tweet{}
	filter := &Filter{
		Patterns: []string{invalidRegexpPattern},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	_, err := filter.CheckTweet(tweet, nil, nil, cache)
	require.Error(t, err)
}

func TestFilter_CheckTweet_PatternsUnmatched(t *testing.T) {
	tweet := anaconda.Tweet{
		Text: "fizz buzz",
	}
	filter := &Filter{
		Patterns: []string{"foo"},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckTweet(tweet, nil, nil, cache)
	require.NoError(t, err)
	require.False(t, match)
}

func TestFilter_CheckTweet_URLPatternsMatched(t *testing.T) {
	tweet := generateTweetWithURL()
	filter := &Filter{
		URLPatterns: []string{"http://example.com"},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckTweet(tweet, nil, nil, cache)
	require.NoError(t, err)
	require.True(t, match)
}

func TestFilter_CheckTweet_URLPatternsInvalid(t *testing.T) {
	tweet := generateTweetWithURL()
	filter := &Filter{
		URLPatterns: []string{invalidRegexpPattern},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	_, err := filter.CheckTweet(tweet, nil, nil, cache)
	require.Error(t, err)
}

func TestFilter_CheckTweet_URLPatternsUnmatched(t *testing.T) {
	tweet := generateTweetWithURL()
	filter := &Filter{
		URLPatterns: []string{"http://example2.com"},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckTweet(tweet, nil, nil, cache)
	require.NoError(t, err)
	require.False(t, match)
}

func generateTweetWithURL() anaconda.Tweet {
	return anaconda.Tweet{
		Entities: anaconda.Entities{
			Urls: []struct {
				Indices      []int  `json:"indices"`
				Url          string `json:"url"`
				Display_url  string `json:"display_url"`
				Expanded_url string `json:"expanded_url"`
			}{
				{
					Display_url: "http://example.com/foo",
				},
			},
		},
	}
}

func TestFilter_CheckTweet_NotHasMedia(t *testing.T) {
	tweet := anaconda.Tweet{}
	filter := &Filter{}
	filter.HasMedia = utils.TruePtr()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckTweet(tweet, nil, nil, cache)
	require.NoError(t, err)
	require.False(t, match)
}

func TestFilter_CheckTweet_NotRetweeted(t *testing.T) {
	tweet := anaconda.Tweet{}
	filter := &Filter{}
	filter.Retweeted = utils.TruePtr()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckTweet(tweet, nil, nil, cache)
	require.NoError(t, err)
	require.False(t, match)
}

func TestFilter_CheckTweet_FavoriteThresholdExceeded(t *testing.T) {
	tweet := anaconda.Tweet{
		FavoriteCount: 100,
	}
	filter := NewFilter()
	filter.FavoriteThreshold = utils.IntPtr(99)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckTweet(tweet, nil, nil, cache)
	require.NoError(t, err)
	require.True(t, match)
}

func TestFilter_CheckTweet_FavoriteThresholdNotExceeded(t *testing.T) {
	tweet := anaconda.Tweet{
		FavoriteCount: 100,
	}
	filter := NewFilter()
	filter.FavoriteThreshold = utils.IntPtr(101)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckTweet(tweet, nil, nil, cache)
	require.NoError(t, err)
	require.False(t, match)
}

func TestFilter_CheckTweet_RetweetedThresholdExceeded(t *testing.T) {
	tweet := anaconda.Tweet{
		RetweetCount: 100,
	}
	filter := NewFilter()
	filter.RetweetedThreshold = utils.IntPtr(99)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckTweet(tweet, nil, nil, cache)
	require.NoError(t, err)
	require.True(t, match)
}

func TestFilter_CheckTweet_RetweetedThresholdNotExceeded(t *testing.T) {
	tweet := anaconda.Tweet{
		RetweetCount: 100,
	}
	filter := NewFilter()
	filter.RetweetedThreshold = utils.IntPtr(101)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckTweet(tweet, nil, nil, cache)
	require.NoError(t, err)
	require.False(t, match)
}

func TestFilter_CheckTweet_LangNotMatched(t *testing.T) {
	tweet := anaconda.Tweet{
		Lang: "JP",
	}
	filter := &Filter{}
	filter.Lang = "EN"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckTweet(tweet, nil, nil, cache)
	require.NoError(t, err)
	require.False(t, match)
}

func TestFilter_CheckTweet_VisionMatched(t *testing.T) {
	tweet := generateVisionTweet()
	filter := generateVisionFilter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	v := generateVisionMatcher(ctrl, true, nil)
	cache := generateVisionCache(ctrl)

	match, err := filter.CheckTweet(tweet, v, nil, cache)
	require.NoError(t, err)
	require.True(t, match)
}

func TestFilter_CheckTweet_VisionUnmatched(t *testing.T) {
	tweet := generateVisionTweet()
	filter := generateVisionFilter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	v := generateVisionMatcher(ctrl, false, nil)
	cache := generateVisionCache(ctrl)

	match, err := filter.CheckTweet(tweet, v, nil, cache)
	require.NoError(t, err)
	require.False(t, match)
}

func TestFilter_CheckTweet_VisionError(t *testing.T) {
	tweet := generateVisionTweet()
	filter := generateVisionFilter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	v := generateVisionMatcher(ctrl, false, fmt.Errorf(""))
	cache := generateVisionCache(ctrl)

	_, err := filter.CheckTweet(tweet, v, nil, cache)
	require.Error(t, err)
}

func generateVisionTweet() anaconda.Tweet {
	return anaconda.Tweet{
		Entities: anaconda.Entities{
			Media: []anaconda.EntityMedia{
				{},
			},
		},
	}
}

func generateVisionFilter() *Filter {
	return &Filter{
		Vision: models.VisionCondition{
			Label: []string{"foo"},
		},
	}
}

func generateVisionMatcher(ctrl *gomock.Controller, match bool, err error) VisionMatcher {
	v := mocks.NewMockVisionMatcher(ctrl)
	v.EXPECT().Enabled().Return(true)
	v.EXPECT().
		MatchImages(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]string{""}, []bool{match}, err)
	return v
}

func generateVisionCache(ctrl *gomock.Controller) data.Cache {
	cache := mocks.NewMockCache(ctrl)
	cache.EXPECT().GetLatestImages(gomock.Any()).Return([]models.ImageCacheData{})
	return cache
}

func TestFilter_CheckTweet_LanguageMatched(t *testing.T) {
	tweet := anaconda.Tweet{Text: ""}
	filter := generateLanguageFilter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	l := generateLanguageMatcher(ctrl, true, nil)
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckTweet(tweet, nil, l, cache)
	require.NoError(t, err)
	require.True(t, match)
}

func TestFilter_CheckTweet_LanguageUnmatched(t *testing.T) {
	tweet := anaconda.Tweet{Text: ""}
	filter := generateLanguageFilter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	l := generateLanguageMatcher(ctrl, false, nil)
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckTweet(tweet, nil, l, cache)
	require.NoError(t, err)
	require.False(t, match)
}

func TestFilter_CheckTweet_LanguageError(t *testing.T) {
	tweet := anaconda.Tweet{Text: ""}
	filter := generateLanguageFilter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	l := generateLanguageMatcher(ctrl, false, fmt.Errorf(""))
	cache := mocks.NewMockCache(ctrl)

	_, err := filter.CheckTweet(tweet, nil, l, cache)
	require.Error(t, err)
}

func generateLanguageFilter() *Filter {
	return &Filter{
		Language: models.LanguageCondition{
			MaxSentiment: utils.Float64Ptr(0.5),
			MinSentiment: utils.Float64Ptr(0.0),
		},
	}
}

func generateLanguageMatcher(ctrl *gomock.Controller, match bool, err error) LanguageMatcher {
	l := mocks.NewMockLanguageMatcher(ctrl)
	l.EXPECT().Enabled().Return(true)
	l.EXPECT().
		MatchText(gomock.Any(), gomock.Any()).
		Return("", match, err)
	return l
}

func TestFilter_CheckSlackMsg_PatternsMatched(t *testing.T) {
	filter := &Filter{
		Patterns: []string{"foo"},
	}
	ev := &slack.MessageEvent{}
	ev.Text = "foo is bar"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckSlackMsg(ev, nil, nil, cache)
	require.NoError(t, err)
	require.True(t, match)
}

func TestFilter_CheckSlackMsg_PatternsInvalid(t *testing.T) {
	ev := &slack.MessageEvent{}
	filter := &Filter{
		Patterns: []string{invalidRegexpPattern},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	_, err := filter.CheckSlackMsg(ev, nil, nil, cache)
	require.Error(t, err)
}

func TestFilter_CheckSlackMsg_PatternsUnmatched(t *testing.T) {
	filter := &Filter{
		Patterns: []string{"foo"},
	}
	ev := &slack.MessageEvent{}
	ev.Text = "fizz buzz"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckSlackMsg(ev, nil, nil, cache)
	require.NoError(t, err)
	require.False(t, match)
}

func TestFilter_CheckSlackMsg_NotHasMedia(t *testing.T) {
	ev := &slack.MessageEvent{}
	ev.Attachments = []slack.Attachment{}
	filter := &Filter{}
	filter.HasMedia = utils.TruePtr()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckSlackMsg(ev, nil, nil, cache)
	require.NoError(t, err)
	require.False(t, match)
}

func TestFilter_CheckSlackMsg_VisionMatched(t *testing.T) {
	ev := generateVisionSlackMessageEvent()
	filter := generateVisionFilter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	v := generateVisionMatcher(ctrl, true, nil)
	cache := generateVisionCache(ctrl)

	match, err := filter.CheckSlackMsg(ev, v, nil, cache)
	require.NoError(t, err)
	require.True(t, match)
}

func TestFilter_CheckSlackMsg_VisionError(t *testing.T) {
	ev := generateVisionSlackMessageEvent()
	filter := generateVisionFilter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	v := generateVisionMatcher(ctrl, false, fmt.Errorf(""))
	cache := generateVisionCache(ctrl)

	_, err := filter.CheckSlackMsg(ev, v, nil, cache)
	require.Error(t, err)
}

func TestFilter_CheckSlackMsg_VisionUnmatched(t *testing.T) {
	ev := generateVisionSlackMessageEvent()
	filter := generateVisionFilter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	v := generateVisionMatcher(ctrl, false, nil)
	cache := generateVisionCache(ctrl)

	match, err := filter.CheckSlackMsg(ev, v, nil, cache)
	require.NoError(t, err)
	require.False(t, match)
}

func TestFilter_CheckSlackMsg_LanguageMatched(t *testing.T) {
	ev := &slack.MessageEvent{}
	ev.Text = ""
	filter := generateLanguageFilter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	l := generateLanguageMatcher(ctrl, true, nil)
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckSlackMsg(ev, nil, l, cache)
	require.NoError(t, err)
	require.True(t, match)
}

func TestFilter_CheckSlackMsg_LanguageUnmatched(t *testing.T) {
	ev := &slack.MessageEvent{}
	ev.Text = ""
	filter := generateLanguageFilter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	l := generateLanguageMatcher(ctrl, false, nil)
	cache := mocks.NewMockCache(ctrl)

	match, err := filter.CheckSlackMsg(ev, nil, l, cache)
	require.NoError(t, err)
	require.False(t, match)
}

func TestFilter_CheckSlackMsg_LanguageError(t *testing.T) {
	ev := &slack.MessageEvent{}
	ev.Text = ""
	filter := generateLanguageFilter()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	l := generateLanguageMatcher(ctrl, false, fmt.Errorf(""))
	cache := mocks.NewMockCache(ctrl)

	_, err := filter.CheckSlackMsg(ev, nil, l, cache)
	require.Error(t, err)
}

func generateVisionSlackMessageEvent() *slack.MessageEvent {
	ev := &slack.MessageEvent{}
	att := slack.Attachment{
		ImageURL: "url",
	}
	ev.Attachments = []slack.Attachment{att}
	return ev
}

func TestFilter_Validate(t *testing.T) {
	f := NewFilter()
	f.FavoriteThreshold = utils.IntPtr(100)
	f.Vision.Label = []string{"foo"}
	err := f.Validate()
	require.Error(t, err)
}

func TestFilter_ShouldRepeat(t *testing.T) {
	filter := &Filter{}

	require.False(t, filter.ShouldRepeat())

	filter.FavoriteThreshold = utils.IntPtr(100)
	filter.RetweetedThreshold = nil
	require.True(t, filter.ShouldRepeat())

	filter.FavoriteThreshold = nil
	filter.RetweetedThreshold = utils.IntPtr(100)
	require.True(t, filter.ShouldRepeat())
}
