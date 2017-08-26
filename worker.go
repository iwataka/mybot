package main

import (
	"fmt"
	"log"
	"time"

	"github.com/iwataka/anaconda"
	mybot "github.com/iwataka/mybot/lib"
	worker "github.com/iwataka/mybot/worker"
)

func manageWorkerWithStart(key int, workerChans map[int]chan int, statuses map[int]*bool, w worker.RoutineWorker) {
	ch, exists := workerChans[key]
	if !exists {
		ch = make(chan int)
		workerChans[key] = ch
	}
	outChan := make(chan interface{})
	go worker.ManageWorker(ch, outChan, w)
	go func() {
		for {
			select {
			case msg := <-outChan:
				switch m := msg.(type) {
				case bool:
					if m {
						log.Printf("Start process: %s", w.Name())
						*statuses[key] = true
					} else {
						log.Printf("Stop process: %s", w.Name())
						*statuses[key] = false
					}
				case error:
					log.Printf("Error: %s (%s)", m.Error(), w.Name())
				case string:
					log.Printf("Message: %s (%s)", m, w.Name())
				}
			}
		}
	}()
	ch <- worker.StartSignal
}

func reloadWorkers(userID string) {
	data := userSpecificDataMap[userID]
	for _, ch := range data.workerChans {
		ch <- worker.RestartSignal
	}
}

type twitterDMWorker struct {
	twitterAPI *mybot.TwitterAPI
	id         string
	stream     *anaconda.Stream
}

func newTwitterDMWorker(twitterAPI *mybot.TwitterAPI, id string) *twitterDMWorker {
	return &twitterDMWorker{twitterAPI, id, nil}
}

func (w *twitterDMWorker) Start() error {
	if !twitterAPIIsAvailable(w.twitterAPI) {
		return nil
	}

	r := w.twitterAPI.DefaultDirectMessageReceiver
	listener, err := w.twitterAPI.ListenMyself(nil, r)
	if err != nil {
		return err
	}
	w.stream = listener.Stream
	if err := listener.Listen(); err != nil {
		return err
	}
	return nil
}

func (w *twitterDMWorker) Stop() {
	if w.stream != nil {
		w.stream.Stop()
	}
}

func (w *twitterDMWorker) Name() string {
	return fmt.Sprintf("%s Twitter DM Worker", w.id)
}

type twitterUserWorker struct {
	twitterAPI  *mybot.TwitterAPI
	slackAPI    *mybot.SlackAPI
	visionAPI   *mybot.VisionAPI
	languageAPI *mybot.LanguageAPI
	cache       mybot.Cache
	id          string
	stream      *anaconda.Stream
}

func newTwitterUserWorker(
	twitterAPI *mybot.TwitterAPI,
	slackAPI *mybot.SlackAPI,
	visionAPI *mybot.VisionAPI,
	languageAPI *mybot.LanguageAPI,
	cache mybot.Cache,
	id string,
) *twitterUserWorker {
	return &twitterUserWorker{twitterAPI, slackAPI, visionAPI, languageAPI, cache, id, nil}
}

func (w *twitterUserWorker) Start() error {
	if !twitterAPIIsAvailable(w.twitterAPI) {
		return nil
	}

	listener, err := w.twitterAPI.ListenUsers(nil)
	if err != nil {
		return err
	}
	w.stream = listener.Stream
	if err := listener.Listen(w.visionAPI, w.languageAPI, w.slackAPI, w.cache); err != nil {
		return err
	}
	return nil
}

func (w *twitterUserWorker) Stop() {
	if w.stream != nil {
		w.stream.Stop()
	}
}

func (w *twitterUserWorker) Name() string {
	return fmt.Sprintf("%s Twitter User Worker", w.id)
}

type twitterPeriodicWorker struct {
	twitterAPI  *mybot.TwitterAPI
	slackAPI    *mybot.SlackAPI
	visionAPI   *mybot.VisionAPI
	languageAPI *mybot.LanguageAPI
	cache       mybot.Cache
	config      mybot.Config
	id          string
	stream      *anaconda.Stream
	status      *bool
}

func newTwitterPeriodicWorker(
	twitterAPI *mybot.TwitterAPI,
	slackAPI *mybot.SlackAPI,
	visionAPI *mybot.VisionAPI,
	languageAPI *mybot.LanguageAPI,
	cache mybot.Cache,
	config mybot.Config,
	id string,
) *twitterPeriodicWorker {
	statusInitValue := false
	return &twitterPeriodicWorker{twitterAPI, slackAPI, visionAPI, languageAPI, cache, config, id, nil, &statusInitValue}
}

func (w *twitterPeriodicWorker) Start() error {
	if !twitterAPIIsAvailable(w.twitterAPI) {
		return nil
	}

	*w.status = true
	for *w.status {
		if err := runTwitterWithStream(w.twitterAPI, w.slackAPI, w.visionAPI, w.languageAPI, w.config); err != nil {
			return err
		}
		if err := w.cache.Save(); err != nil {
			return err
		}
		d, err := time.ParseDuration(w.config.GetTwitterDuration())
		if err != nil {
			return err
		}
		time.Sleep(d)
	}
	return nil
}

func (w *twitterPeriodicWorker) Stop() {
	*w.status = false
}

func (w *twitterPeriodicWorker) Name() string {
	return fmt.Sprintf("%s Twitter Periodic Worker", w.id)
}

type slackWorker struct {
	slackAPI    *mybot.SlackAPI
	twitterAPI  *mybot.TwitterAPI
	visionAPI   *mybot.VisionAPI
	languageAPI *mybot.LanguageAPI
	id          string
	listener    *mybot.SlackListener
}

func newSlackWorker(
	slackAPI *mybot.SlackAPI,
	twitterAPI *mybot.TwitterAPI,
	visionAPI *mybot.VisionAPI,
	languageAPI *mybot.LanguageAPI,
	id string,
) *slackWorker {
	return &slackWorker{slackAPI, twitterAPI, visionAPI, languageAPI, id, nil}
}

func (w *slackWorker) Start() error {
	if w.slackAPI == nil || !w.slackAPI.Enabled() {
		return nil
	}

	w.listener = w.slackAPI.Listen()
	if err := w.listener.Start(w.visionAPI, w.languageAPI, w.twitterAPI); err != nil {
		return err
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
