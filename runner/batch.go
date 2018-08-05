package runner

import (
	"github.com/iwataka/mybot/lib"
	"github.com/iwataka/mybot/utils"

	"fmt"
	"net/url"
)

// BatchRunner wraps a batch process and provides a feature to run it.
type BatchRunner interface {
	Run() error
	// IsAvailable returns true if this runner is available.
	// You should check the availability by calling this and if this
	// returns false, you can't call Run.
	IsAvailable() error
}

// BatchRunnerUsedWithStream implements mybot batch processing and is intended
// to be used with stream processing.
type BatchRunnerUsedWithStream struct {
	twitterAPI  *mybot.TwitterAPI
	slackAPI    *mybot.SlackAPI
	visionAPI   mybot.VisionMatcher
	languageAPI mybot.LanguageMatcher
	config      mybot.Config
}

// NewBatchRunnerUsedWithStream returns new BatchRunnerUsedWithStream with
// specified arguments.
func NewBatchRunnerUsedWithStream(
	twitterAPI *mybot.TwitterAPI,
	slackAPI *mybot.SlackAPI,
	visionAPI mybot.VisionMatcher,
	languageAPI mybot.LanguageMatcher,
	config mybot.Config,
) *BatchRunnerUsedWithStream {
	return &BatchRunnerUsedWithStream{twitterAPI, slackAPI, visionAPI, languageAPI, config}
}

// Run processes Twitter search/favorite API result and then makes notifications
// of it based on r.config.
func (r BatchRunnerUsedWithStream) Run() error {
	for _, a := range r.config.GetTwitterSearches() {
		v := url.Values{}
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		if len(a.ResultType) != 0 {
			v.Set("result_type", a.ResultType)
		}
		for _, query := range a.Queries {
			_, err := r.twitterAPI.ProcessSearch(
				query,
				v,
				a.Filter,
				r.visionAPI,
				r.languageAPI,
				r.slackAPI,
				a.Action,
			)
			if err != nil {
				return utils.WithStack(err)
			}
		}
	}

	for _, a := range r.config.GetTwitterFavorites() {
		v := url.Values{}
		if a.Count != nil {
			v.Set("count", fmt.Sprintf("%d", *a.Count))
		}
		for _, name := range a.ScreenNames {
			_, err := r.twitterAPI.ProcessFavorites(
				name,
				v,
				a.Filter,
				r.visionAPI,
				r.languageAPI,
				r.slackAPI,
				a.Action,
			)
			if err != nil {
				return utils.WithStack(err)
			}
		}
	}

	return nil
}

// IsAvailable returns true if Twitter API is available because it is data
// fetcher and all other API depends on it. It is the responsibility of
// Twitter API to check other APIs are available.
func (r BatchRunnerUsedWithStream) IsAvailable() error {
	return TwitterAPIIsAvailable(r.twitterAPI)
}

// BatchRunnerUsedWithoutStream implements mybot batch processing and be
// intended to be used without stream processing (that means for command-line
// usage).
type BatchRunnerUsedWithoutStream struct {
	baseRunner *BatchRunnerUsedWithStream
}

// NewBatchRunnerUsedWithoutStream returns a new BatchRunnerUsedWithoutStream
// based on a specified baseRunner.
func NewBatchRunnerUsedWithoutStream(baseRunner *BatchRunnerUsedWithStream) *BatchRunnerUsedWithoutStream {
	return &BatchRunnerUsedWithoutStream{baseRunner}
}

// Run firstly calls r.baseRunner.Run().
// Then processes Twitter timeline API result and makes notifications of it
// based on r.baseRunner.config.
//
// TODO: Implement slack stream processing (or abolish this method and
// command-line usage)
func (r BatchRunnerUsedWithoutStream) Run() error {
	err := r.baseRunner.Run()
	if err != nil {
		return utils.WithStack(err)
	}

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
			_, err := r.baseRunner.twitterAPI.ProcessTimeline(
				name,
				v,
				a.Filter,
				r.baseRunner.visionAPI,
				r.baseRunner.languageAPI,
				r.baseRunner.slackAPI,
				a.Action,
			)
			if err != nil {
				return utils.WithStack(err)
			}
		}
	}

	return nil
}

// IsAvailable returns true if Twitter API is available because it is data
// fetcher and all other API depends on it. It is the responsibility of
// Twitter API to check other APIs are available.
//
// TODO: Check Slack API is available.
func (r BatchRunnerUsedWithoutStream) IsAvailable() error {
	return TwitterAPIIsAvailable(r.baseRunner.twitterAPI)
}

// TwitterAPIIsAvailable returns nil if twitterAPI client is available to use,
// which means that twitterAPI's methods are callable and it is verified by a
// valid credential.
func TwitterAPIIsAvailable(twitterAPI *mybot.TwitterAPI) error {
	if twitterAPI == nil {
		return fmt.Errorf("Twitter API is nil")
	}
	success, err := twitterAPI.VerifyCredentials()
	if err != nil {
		return utils.WithStack(err)
	}
	if !success {
		return fmt.Errorf("Twitter API credential verification failed")
	}
	return nil
}
