package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/core"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/runner"
	"github.com/iwataka/mybot/utils"
	"github.com/iwataka/mybot/worker"
)

type workerMessageLogger struct {
	workerID string
}

func (h workerMessageLogger) HandleError(err error) {
	log.Println(err)
}

func (h workerMessageLogger) HandleWorkerStatus(s worker.WorkerStatus) {
	fmt.Printf("Worker %s: %s\n", s, h.workerID)
}

func activateWorkerAndStart(
	key int,
	workerMgrs map[int]*worker.WorkerManager,
	w models.Worker,
	whandler worker.WorkerManagerOutHandler,
	bufSize int,
	layers ...worker.WorkerChannelLayer,
) {
	if wm, exists := workerMgrs[key]; exists {
		wm.Close()
	}
	// Worker manager process
	wm := worker.NewWorkerManager(w, bufSize, layers...)
	workerMgrs[key] = wm

	go wm.HandleOutput(whandler)
	wm.Send(worker.StartSignal)
}

// TODO: implement this to userSpecificDataMap struct
func restartWorkers(userID string) {
	data := userSpecificDataMap[userID]
	for _, ch := range data.workerMgrs {
		ch.Send(worker.RestartSignal)
	}
}

type twitterDMWorker struct {
	twitterAPI *core.TwitterAPI
	id         string
	listener   *core.TwitterDMListener
}

func newTwitterDMWorker(twitterAPI *core.TwitterAPI, id string) *twitterDMWorker {
	return &twitterDMWorker{twitterAPI, id, nil}
}

func (w *twitterDMWorker) Start(ctx context.Context) error {
	if err := runner.TwitterAPIIsAvailable(w.twitterAPI); err != nil {
		return utils.WithStack(err)
	}

	var err error
	w.listener, err = w.twitterAPI.ListenMyself(nil)
	if err != nil {
		return utils.WithStack(err)
	}
	if err := w.listener.Listen(ctx); err != nil {
		return utils.WithStack(err)
	}
	return nil
}

func (w *twitterDMWorker) Name() string {
	return fmt.Sprintf("%s Twitter DM Worker", w.id)
}

type twitterUserWorker struct {
	twitterAPI  *core.TwitterAPI
	slackAPI    *core.SlackAPI
	visionAPI   core.VisionMatcher
	languageAPI core.LanguageMatcher
	cache       data.Cache
	id          string
	listener    *core.TwitterUserListener
}

func newTwitterUserWorker(
	twitterAPI *core.TwitterAPI,
	slackAPI *core.SlackAPI,
	visionAPI core.VisionMatcher,
	languageAPI core.LanguageMatcher,
	cache data.Cache,
	id string,
) *twitterUserWorker {
	return &twitterUserWorker{twitterAPI, slackAPI, visionAPI, languageAPI, cache, id, nil}
}

func (w *twitterUserWorker) Start(ctx context.Context) error {
	if err := runner.TwitterAPIIsAvailable(w.twitterAPI); err != nil {
		return utils.WithStack(err)
	}

	var err error
	w.listener, err = w.twitterAPI.ListenUsers(nil)
	if err != nil {
		return utils.WithStack(err)
	}
	if err := w.listener.Listen(ctx, w.visionAPI, w.languageAPI, w.slackAPI, w.cache); err != nil {
		return utils.WithStack(err)
	}
	return nil
}

func (w *twitterUserWorker) Name() string {
	return fmt.Sprintf("%s Twitter User Worker", w.id)
}

type twitterPeriodicWorker struct {
	runner runner.BatchRunner
	cache  utils.Savable
	config core.Config
	id     string
	stream *anaconda.Stream
}

func newTwitterPeriodicWorker(
	runner runner.BatchRunner,
	cache utils.Savable,
	config core.Config,
	id string,
) *twitterPeriodicWorker {
	return &twitterPeriodicWorker{runner, cache, config, id, nil}
}

func (w *twitterPeriodicWorker) Start(ctx context.Context) error {
	if err := w.runner.IsAvailable(); err != nil {
		return utils.WithStack(err)
	}

	d, err := time.ParseDuration(w.config.GetPollingDuration())
	if err != nil {
		return utils.WithStack(err)
	}
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			if err := w.runner.Run(); err != nil {
				return utils.WithStack(err)
			}
			if err := w.cache.Save(); err != nil {
				return utils.WithStack(err)
			}
		case <-ctx.Done():
			return utils.NewStreamInterruptedError()
		}
	}
}

func (w *twitterPeriodicWorker) Name() string {
	return fmt.Sprintf("%s Twitter Periodic Worker", w.id)
}

type slackWorker struct {
	slackAPI    *core.SlackAPI
	twitterAPI  *core.TwitterAPI
	visionAPI   core.VisionMatcher
	languageAPI core.LanguageMatcher
	id          string
	listener    *core.SlackListener
}

func newSlackWorker(
	slackAPI *core.SlackAPI,
	twitterAPI *core.TwitterAPI,
	visionAPI core.VisionMatcher,
	languageAPI core.LanguageMatcher,
	id string,
) *slackWorker {
	return &slackWorker{slackAPI, twitterAPI, visionAPI, languageAPI, id, nil}
}

func (w *slackWorker) Start(ctx context.Context) error {
	if w.slackAPI == nil {
		return fmt.Errorf("Slack API is nil")
	}
	if !w.slackAPI.Enabled() {
		return fmt.Errorf("Slack API is disabled")
	}

	w.listener = w.slackAPI.Listen()
	if err := w.listener.Start(ctx, w.visionAPI, w.languageAPI, w.twitterAPI); err != nil {
		return utils.WithStack(err)
	}
	return nil
}

func (w *slackWorker) Name() string {
	return fmt.Sprintf("%s Slack Worker", w.id)
}
