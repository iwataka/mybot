package main

import (
	"fmt"
	"log"
	"time"

	"github.com/iwataka/anaconda"
	mybot "github.com/iwataka/mybot/lib"
	worker "github.com/iwataka/mybot/worker"
)

// TODO: Make more simple and testable. I want to trace that workerChans is
// used and statuses is changed correctly.
func activateWorkerAndStart(key int, workerChans map[int]chan *worker.WorkerSignal, statuses map[int]*bool, w worker.Worker) {
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
				fmt.Printf("Worker %s: %s\n", m, w.Name())
				switch m {
				case worker.StatusStarted:
					*statuses[key] = true
				case worker.StatusStopped:
					*statuses[key] = false
				}
			case error:
				log.Printf("%+v\n", m)
			}
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
	if err := mybot.TwitterAPIIsAvailable(w.twitterAPI); err != nil {
		return mybot.WithStack(err)
	}

	var err error
	r := w.twitterAPI.DefaultDirectMessageReceiver
	w.listener, err = w.twitterAPI.ListenMyself(nil, r, w.timeout)
	if err != nil {
		return mybot.WithStack(err)
	}
	if err := w.listener.Listen(); err != nil {
		return mybot.WithStack(err)
	}
	return nil
}

func (w *twitterDMWorker) Stop() {
	if w.listener != nil {
		w.listener.Stop()
	}
}

func (w *twitterDMWorker) Name() string {
	return fmt.Sprintf("%s Twitter DM Worker", w.id)
}

type twitterUserWorker struct {
	twitterAPI  *mybot.TwitterAPI
	slackAPI    *mybot.SlackAPI
	visionAPI   mybot.VisionMatcher
	languageAPI mybot.LanguageMatcher
	cache       mybot.Cache
	id          string
	listener    *mybot.TwitterUserListener
	timeout     time.Duration
}

func newTwitterUserWorker(
	twitterAPI *mybot.TwitterAPI,
	slackAPI *mybot.SlackAPI,
	visionAPI mybot.VisionMatcher,
	languageAPI mybot.LanguageMatcher,
	cache mybot.Cache,
	id string,
	timeout time.Duration,
) *twitterUserWorker {
	return &twitterUserWorker{twitterAPI, slackAPI, visionAPI, languageAPI, cache, id, nil, timeout}
}

func (w *twitterUserWorker) Start() error {
	if err := mybot.TwitterAPIIsAvailable(w.twitterAPI); err != nil {
		return mybot.WithStack(err)
	}

	var err error
	w.listener, err = w.twitterAPI.ListenUsers(nil, w.timeout)
	if err != nil {
		return mybot.WithStack(err)
	}
	if err := w.listener.Listen(w.visionAPI, w.languageAPI, w.slackAPI, w.cache); err != nil {
		return mybot.WithStack(err)
	}
	return nil
}

func (w *twitterUserWorker) Stop() {
	if w.listener != nil {
		w.listener.Stop()
	}
}

func (w *twitterUserWorker) Name() string {
	return fmt.Sprintf("%s Twitter User Worker", w.id)
}

type twitterPeriodicWorker struct {
	runner    mybot.BatchRunner
	cache     mybot.Savable
	duration  string
	timeout   time.Duration
	id        string
	stream    *anaconda.Stream
	innerChan chan bool
}

func newTwitterPeriodicWorker(
	runner mybot.BatchRunner,
	cache mybot.Savable,
	duration string,
	timeout time.Duration,
	id string,
) *twitterPeriodicWorker {
	return &twitterPeriodicWorker{runner, cache, duration, timeout, id, nil, make(chan bool)}
}

func (w *twitterPeriodicWorker) Start() error {
	if err := w.runner.Verify(); err != nil {
		return mybot.WithStack(err)
	}

	d, err := time.ParseDuration(w.duration)
	if err != nil {
		return mybot.WithStack(err)
	}
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			if err := w.runner.Run(); err != nil {
				return mybot.WithStack(err)
			}
			if err := w.cache.Save(); err != nil {
				return mybot.WithStack(err)
			}
		case <-w.innerChan:
			return mybot.NewInterruptedError()
		}
	}
	return nil
}

func (w *twitterPeriodicWorker) Stop() {
	select {
	case w.innerChan <- true:
	case <-time.After(w.timeout):
		log.Printf("Faield to stop worker %s\n", w.Name())
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
		return mybot.Errorf("Slack API is nil")
	}
	if !w.slackAPI.Enabled() {
		return mybot.Errorf("Slack API is disabled")
	}

	w.listener = w.slackAPI.Listen()
	if err := w.listener.Start(w.visionAPI, w.languageAPI, w.twitterAPI); err != nil {
		return mybot.WithStack(err)
	}
	return nil
}

func (w *slackWorker) Stop() {
	if w.listener != nil {
		w.listener.Stop()
	}
}

func (w *slackWorker) Name() string {
	return fmt.Sprintf("%s Slack Worker", w.id)
}
