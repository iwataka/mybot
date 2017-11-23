package mybot

import (
	"fmt"
	"net/url"

	"github.com/iwataka/anaconda"
)

type BatchRunner interface {
	Run() error
	Verify() error
}

type BatchRunnerWithStream struct {
	twitterAPI  *TwitterAPI
	slackAPI    *SlackAPI
	visionAPI   VisionMatcher
	languageAPI LanguageMatcher
	config      Config
}

func NewBatchRunnerWithStream(
	twitterAPI *TwitterAPI,
	slackAPI *SlackAPI,
	visionAPI VisionMatcher,
	languageAPI LanguageMatcher,
	config Config,
) *BatchRunnerWithStream {
	return &BatchRunnerWithStream{twitterAPI, slackAPI, visionAPI, languageAPI, config}
}

func (r BatchRunnerWithStream) Run() error {
	tweets := []anaconda.Tweet{}
	for _, a := range r.config.GetTwitterSearches() {
		v := url.Values{}
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		if len(a.ResultType) != 0 {
			v.Set("result_type", a.ResultType)
		}
		for _, query := range a.Queries {
			ts, err := r.twitterAPI.ProcessSearch(
				query,
				v,
				a.Filter,
				r.visionAPI,
				r.languageAPI,
				r.slackAPI,
				a.Action,
			)
			if err != nil {
				return err
			}
			tweets = append(tweets, ts...)
		}
	}
	for _, a := range r.config.GetTwitterFavorites() {
		v := url.Values{}
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		for _, name := range a.ScreenNames {
			ts, err := r.twitterAPI.ProcessFavorites(
				name,
				v,
				a.Filter,
				r.visionAPI,
				r.languageAPI,
				r.slackAPI,
				a.Action,
			)
			tweets = append(tweets, ts...)
			if err != nil {
				return err
			}
		}
	}
	for _, t := range tweets {
		err := r.twitterAPI.NotifyToAll(&t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r BatchRunnerWithStream) Verify() error {
	return TwitterAPIIsAvailable(r.twitterAPI)
}

type BatchRunnerWithoutStream struct {
	baseRunner *BatchRunnerWithStream
}

func NewBatchRunnerWithoutStream(baseRunner *BatchRunnerWithStream) *BatchRunnerWithoutStream {
	return &BatchRunnerWithoutStream{baseRunner}
}

func (r BatchRunnerWithoutStream) Run() error {
	err := r.baseRunner.Run()
	if err != nil {
		return err
	}
	tweets := []anaconda.Tweet{}
	for _, a := range r.baseRunner.config.GetTwitterTimelines() {
		v := url.Values{}
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		if a.ExcludeReplies != nil {
			v.Set("exclude_replies", fmt.Sprintf("%v", *a.ExcludeReplies))
		}
		if a.IncludeRts != nil {
			v.Set("include_rts", fmt.Sprintf("%v", *a.IncludeRts))
		}
		for _, name := range a.ScreenNames {
			ts, err := r.baseRunner.twitterAPI.ProcessTimeline(
				name,
				v,
				a.Filter,
				r.baseRunner.visionAPI,
				r.baseRunner.languageAPI,
				r.baseRunner.slackAPI,
				a.Action,
			)
			tweets = append(tweets, ts...)
			if err != nil {
				return err
			}
		}
	}
	for _, t := range tweets {
		err := r.baseRunner.twitterAPI.NotifyToAll(&t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r BatchRunnerWithoutStream) Verify() error {
	return TwitterAPIIsAvailable(r.baseRunner.twitterAPI)
}

func TwitterAPIIsAvailable(twitterAPI *TwitterAPI) error {
	if twitterAPI == nil {
		return fmt.Errorf("Twitter API is nil")
	}
	success, err := twitterAPI.VerifyCredentials()
	if !success {
		return fmt.Errorf("Twitter API credential verification failed")
	}
	if err != nil {
		return err
	}
	return nil
}
