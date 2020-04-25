package runner

import (
	"github.com/iwataka/mybot/core"
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
	twitterAPI  *core.TwitterAPI
	slackAPI    *core.SlackAPI
	visionAPI   core.VisionMatcher
	languageAPI core.LanguageMatcher
	config      core.Config
}

// NewBatchRunnerUsedWithStream returns new BatchRunnerUsedWithStream with
// specified arguments.
func NewBatchRunnerUsedWithStream(
	twitterAPI *core.TwitterAPI,
	slackAPI *core.SlackAPI,
	visionAPI core.VisionMatcher,
	languageAPI core.LanguageMatcher,
	config core.Config,
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

// TwitterAPIIsAvailable returns nil if twitterAPI client is available to use,
// which means that twitterAPI's methods are callable and it is verified by a
// valid credential.
func TwitterAPIIsAvailable(twitterAPI *core.TwitterAPI) error {
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
