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
	workerID  string
	logger    *log.Logger
	errLogger *log.Logger
}

func (h workerMessageLogger) Handle(out interface{}) {
	switch out.(type) {
	case error:
		h.errLogger.Printf("Worker[%s]: %#v\n", h.workerID, out)
	case worker.WorkerStatus:
		h.logger.Printf("Worker %s: %s\n", out, h.workerID)
	case fmt.Stringer:
		h.logger.Printf("Worker[%s]: %s\n", h.workerID, out)
	}
}

func activateWorkerAndStart(
	w models.Worker,
	whandler worker.WorkerManagerOutHandler,
	bufSize int,
	layers ...worker.WorkerChannelLayer,
) *worker.WorkerManager {
	wm := worker.NewWorkerManager(w, bufSize, layers...)
	go wm.HandleOutput(whandler)
	wm.Send(worker.StartSignal)
	return wm
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

func (w *twitterDMWorker) Start(ctx context.Context, outChan chan<- interface{}) error {
	if err := runner.TwitterAPIIsAvailable(w.twitterAPI); err != nil {
		return utils.WithStack(err)
	}

	var err error
	w.listener, err = w.twitterAPI.ListenMyself(nil)
	if err != nil {
		return utils.WithStack(err)
	}
	if err := w.listener.Listen(ctx, outChan); err != nil {
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

func (w *twitterUserWorker) Start(ctx context.Context, outChan chan<- interface{}) error {
	if err := runner.TwitterAPIIsAvailable(w.twitterAPI); err != nil {
		return utils.WithStack(err)
	}

	var err error
	w.listener, err = w.twitterAPI.ListenUsers(nil, w.visionAPI, w.languageAPI, w.slackAPI, w.cache)
	if err != nil {
		return utils.WithStack(err)
	}
	if err := w.listener.Listen(ctx, outChan); err != nil {
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

func (w *twitterPeriodicWorker) Start(ctx context.Context, outChan chan<- interface{}) error {
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
			outChan <- core.NewReceivedEvent(core.TwitterEventType, "periodic", nil)
			processedTweets, processedActions, err := w.runner.Run()
			if err != nil {
				return utils.WithStack(err)
			}
			for i := 0; i < len(processedTweets); i++ {
				outChan <- core.NewActionEvent(processedActions[i], processedTweets[i])
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

func (w *slackWorker) Start(ctx context.Context, outChan chan<- interface{}) error {
	if w.slackAPI == nil {
		return fmt.Errorf("Slack API is nil")
	}
	if !w.slackAPI.Enabled() {
		return fmt.Errorf("Slack API is disabled")
	}

	w.listener = w.slackAPI.Listen(w.visionAPI, w.languageAPI, w.twitterAPI)
	if err := w.listener.Start(ctx, outChan); err != nil {
		return utils.WithStack(err)
	}
	return nil
}

func (w *slackWorker) Name() string {
	return fmt.Sprintf("%s Slack Worker", w.id)
}
