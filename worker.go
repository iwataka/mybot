package main

import (
	"fmt"
	"log"
	"time"

	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/lib"
	"github.com/iwataka/mybot/runner"
	"github.com/iwataka/mybot/utils"
	"github.com/iwataka/mybot/worker"
)

type WorkerMessageHandler interface {
	Handle(msg interface{}) error
}

type DefaultWorkerMessageHandler struct {
	config     mybot.Config
	twitterAPI *mybot.TwitterAPI
	slackAPI   *mybot.SlackAPI
	workerID   string
}

func (h DefaultWorkerMessageHandler) Handle(msg interface{}) error {
	logMessage := ""
	switch m := msg.(type) {
	case worker.WorkerStatus:
		logMessage = fmt.Sprintf("Worker %s: %s", m, h.workerID)
	case error:
		logMessage = m.Error()
	}
	sendsSomeone, err := h.config.GetLogNotification().Notify(h.twitterAPI, h.slackAPI, logMessage)
	if err != nil {
		return err
	}
	if !sendsSomeone {
		log.Println(logMessage)
	}
	return nil
}

// TODO: Test statuses are changed correctly.
func activateWorkerAndStart(
	key int,
	workerChans map[int]chan *worker.WorkerSignal,
	statuses map[int]bool,
	w worker.Worker,
	msgHandler WorkerMessageHandler,
) {
	if ch, exists := workerChans[key]; exists {
		close(ch)
	}
	// Worker manager process
	ch, outChan := worker.ActivateWorker(w, time.Minute)
	workerChans[key] = ch
	// Process handling logs from the above worker manager
	go func() {
		for msg := range outChan {
			switch m := msg.(type) {
			case worker.WorkerStatus:
				switch m {
				case worker.StatusStarted:
					statuses[key] = true
				case worker.StatusStopped:
					statuses[key] = false
				}
			}
			msgHandler.Handle(msg)
		}
	}()
	// Process sending ping to worker manager priodically
	go func() {
		ticker := time.NewTicker(time.Minute * 10)
		for range ticker.C {
			select {
			case ch <- worker.NewWorkerSignal(worker.PingSignal):
			case <-time.After(time.Minute):
				log.Printf("Failed to ping worker manager process (timeout: 1m)\n")
			}
		}
	}()
	ch <- worker.NewWorkerSignal(worker.StartSignal)
}

func reloadWorkers(userID string) {
	data := userSpecificDataMap[userID]
	for _, ch := range data.workerChans {
		select {
		case ch <- worker.NewWorkerSignal(worker.RestartSignal):
		case <-time.After(time.Minute):
			log.Println("Failed to reload worker (timeout: 1m)")
		}
	}
}

type twitterDMWorker struct {
	twitterAPI *mybot.TwitterAPI
	id         string
	listener   *mybot.TwitterDMListener
	timeout    time.Duration
}

func newTwitterDMWorker(twitterAPI *mybot.TwitterAPI, id string, timeout time.Duration) *twitterDMWorker {
	return &twitterDMWorker{twitterAPI, id, nil, timeout}
}

func (w *twitterDMWorker) Start() error {
	if err := runner.TwitterAPIIsAvailable(w.twitterAPI); err != nil {
		return utils.WithStack(err)
	}

	var err error
	r := w.twitterAPI.DefaultDirectMessageReceiver
	w.listener, err = w.twitterAPI.ListenMyself(nil, r, w.timeout)
	if err != nil {
		return utils.WithStack(err)
	}
	if err := w.listener.Listen(); err != nil {
		return utils.WithStack(err)
	}
	return nil
}

func (w *twitterDMWorker) Stop() error {
	if w.listener != nil {
		return w.listener.Stop()
	}
	return nil
}

func (w *twitterDMWorker) Name() string {
	return fmt.Sprintf("%s Twitter DM Worker", w.id)
}

type twitterUserWorker struct {
	twitterAPI  *mybot.TwitterAPI
	slackAPI    *mybot.SlackAPI
	visionAPI   mybot.VisionMatcher
	languageAPI mybot.LanguageMatcher
	cache       data.Cache
	id          string
	listener    *mybot.TwitterUserListener
	timeout     time.Duration
}

func newTwitterUserWorker(
	twitterAPI *mybot.TwitterAPI,
	slackAPI *mybot.SlackAPI,
	visionAPI mybot.VisionMatcher,
	languageAPI mybot.LanguageMatcher,
	cache data.Cache,
	id string,
	timeout time.Duration,
) *twitterUserWorker {
	return &twitterUserWorker{twitterAPI, slackAPI, visionAPI, languageAPI, cache, id, nil, timeout}
}

func (w *twitterUserWorker) Start() error {
	if err := runner.TwitterAPIIsAvailable(w.twitterAPI); err != nil {
		return utils.WithStack(err)
	}

	var err error
	w.listener, err = w.twitterAPI.ListenUsers(nil, w.timeout)
	if err != nil {
		return utils.WithStack(err)
	}
	if err := w.listener.Listen(w.visionAPI, w.languageAPI, w.slackAPI, w.cache); err != nil {
		return utils.WithStack(err)
	}
	return nil
}

func (w *twitterUserWorker) Stop() error {
	if w.listener != nil {
		return w.listener.Stop()
	}
	return nil
}

func (w *twitterUserWorker) Name() string {
	return fmt.Sprintf("%s Twitter User Worker", w.id)
}

type twitterPeriodicWorker struct {
	runner    runner.BatchRunner
	cache     utils.Savable
	config    mybot.Config
	timeout   time.Duration
	id        string
	stream    *anaconda.Stream
	innerChan chan bool
}

func newTwitterPeriodicWorker(
	runner runner.BatchRunner,
	cache utils.Savable,
	config mybot.Config,
	timeout time.Duration,
	id string,
) *twitterPeriodicWorker {
	return &twitterPeriodicWorker{runner, cache, config, timeout, id, nil, make(chan bool)}
}

func (w *twitterPeriodicWorker) Start() error {
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
		case <-w.innerChan:
			return utils.NewStreamInterruptedError()
		}
	}
	return nil
}

func (w *twitterPeriodicWorker) Stop() error {
	select {
	case w.innerChan <- true:
		return nil
	case <-time.After(w.timeout):
		return fmt.Errorf("Faield to stop worker %s\n", w.Name())
	}
}

func (w *twitterPeriodicWorker) Name() string {
	return fmt.Sprintf("%s Twitter Periodic Worker", w.id)
}

type slackWorker struct {
	slackAPI    *mybot.SlackAPI
	twitterAPI  *mybot.TwitterAPI
	visionAPI   mybot.VisionMatcher
	languageAPI mybot.LanguageMatcher
	id          string
	listener    *mybot.SlackListener
}

func newSlackWorker(
	slackAPI *mybot.SlackAPI,
	twitterAPI *mybot.TwitterAPI,
	visionAPI mybot.VisionMatcher,
	languageAPI mybot.LanguageMatcher,
	id string,
) *slackWorker {
	return &slackWorker{slackAPI, twitterAPI, visionAPI, languageAPI, id, nil}
}

func (w *slackWorker) Start() error {
	if w.slackAPI == nil {
		return fmt.Errorf("Slack API is nil")
	}
	if !w.slackAPI.Enabled() {
		return fmt.Errorf("Slack API is disabled")
	}

	w.listener = w.slackAPI.Listen()
	if err := w.listener.Start(w.visionAPI, w.languageAPI, w.twitterAPI); err != nil {
		return utils.WithStack(err)
	}
	return nil
}

func (w *slackWorker) Stop() error {
	if w.listener != nil {
		return w.listener.Stop()
	}
	return nil
}

func (w *slackWorker) Name() string {
	return fmt.Sprintf("%s Slack Worker", w.id)
}
