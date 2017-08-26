package main

import (
	"log"
	"time"

	"github.com/iwataka/anaconda"
	mybot "github.com/iwataka/mybot/lib"
	worker "github.com/iwataka/mybot/worker"
)

func manageWorkerWithStart(key int, userID string, w worker.RoutineWorker) {
	data := userSpecificDataMap[userID]
	ch, exists := data.workerChans[key]
	if !exists {
		ch = make(chan int)
		data.workerChans[key] = ch
	}
	go worker.ManageWorker(ch, data.statuses[key], w)
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
	stream     *anaconda.Stream
}

func newTwitterDMWorker(twitterAPI *mybot.TwitterAPI) *twitterDMWorker {
	return &twitterDMWorker{twitterAPI, nil}
}

func (w *twitterDMWorker) Start() {
	if !twitterAPIIsAvailable(w.twitterAPI) {
		return
	}

	r := w.twitterAPI.DefaultDirectMessageReceiver
	listener, err := w.twitterAPI.ListenMyself(nil, r)
	if err != nil {
		log.Print(err)
		return
	}
	w.stream = listener.Stream
	if err := listener.Listen(); err != nil {
		log.Print(err)
		return
	}
	log.Print("Failed to keep connection")
}

func (w *twitterDMWorker) Stop() {
	if w.stream != nil {
		w.stream.Stop()
	}
}

type twitterUserWorker struct {
	twitterAPI  *mybot.TwitterAPI
	slackAPI    *mybot.SlackAPI
	visionAPI   *mybot.VisionAPI
	languageAPI *mybot.LanguageAPI
	cache       mybot.Cache
	stream      *anaconda.Stream
}

func newTwitterUserWorker(
	twitterAPI *mybot.TwitterAPI,
	slackAPI *mybot.SlackAPI,
	visionAPI *mybot.VisionAPI,
	languageAPI *mybot.LanguageAPI,
	cache mybot.Cache,
) *twitterUserWorker {
	return &twitterUserWorker{twitterAPI, slackAPI, visionAPI, languageAPI, cache, nil}
}

func (w *twitterUserWorker) Start() {
	if !twitterAPIIsAvailable(w.twitterAPI) {
		return
	}

	listener, err := w.twitterAPI.ListenUsers(nil)
	if err != nil {
		log.Print(err)
		return
	}
	w.stream = listener.Stream
	if err := listener.Listen(w.visionAPI, w.languageAPI, w.slackAPI, w.cache); err != nil {
		log.Print(err)
		return
	}
	log.Print("Failed to keep connection")
}

func (w *twitterUserWorker) Stop() {
	if w.stream != nil {
		w.stream.Stop()
	}
}

type twitterPeriodicWorker struct {
	twitterAPI  *mybot.TwitterAPI
	slackAPI    *mybot.SlackAPI
	visionAPI   *mybot.VisionAPI
	languageAPI *mybot.LanguageAPI
	cache       mybot.Cache
	config      mybot.Config
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
) *twitterPeriodicWorker {
	statusInitValue := false
	return &twitterPeriodicWorker{twitterAPI, slackAPI, visionAPI, languageAPI, cache, config, nil, &statusInitValue}
}

func (w *twitterPeriodicWorker) Start() {
	if !twitterAPIIsAvailable(w.twitterAPI) {
		return
	}

	*w.status = true
	for *w.status {
		if err := runTwitterWithStream(w.twitterAPI, w.slackAPI, w.visionAPI, w.languageAPI, w.config); err != nil {
			log.Print(err)
			return
		}
		if err := w.cache.Save(); err != nil {
			log.Print(err)
			return
		}
		d, err := time.ParseDuration(w.config.GetTwitterDuration())
		if err != nil {
			log.Print(err)
			return
		}
		time.Sleep(d)
	}
}

func (w *twitterPeriodicWorker) Stop() {
	*w.status = false
}

type slackWorker struct {
	slackAPI    *mybot.SlackAPI
	twitterAPI  *mybot.TwitterAPI
	visionAPI   *mybot.VisionAPI
	languageAPI *mybot.LanguageAPI
	listener    *mybot.SlackListener
}

func newSlackWorker(
	slackAPI *mybot.SlackAPI,
	twitterAPI *mybot.TwitterAPI,
	visionAPI *mybot.VisionAPI,
	languageAPI *mybot.LanguageAPI,
) *slackWorker {
	return &slackWorker{slackAPI, twitterAPI, visionAPI, languageAPI, nil}
}

func (w *slackWorker) Start() {
	if w.slackAPI == nil || !w.slackAPI.Enabled() {
		return
	}

	w.listener = w.slackAPI.Listen()
	if err := w.listener.Start(w.visionAPI, w.languageAPI, w.twitterAPI); err != nil {
		log.Print(err)
		return
	}
	log.Print("Failed to keep connection")
}

func (w *slackWorker) Stop() {
	if w.listener != nil {
		w.listener.Stop()
	}
}
