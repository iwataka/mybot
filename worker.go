package main

import (
	"log"
	"time"

	"github.com/iwataka/anaconda"
	mybot "github.com/iwataka/mybot/lib"
)

type routineWorker interface {
	start()
	stop()
}

func manageWorkerWithStart(key int, userID string, worker routineWorker) {
	data := userSpecificDataMap[userID]
	ch, exists := data.workerChans[key]
	if !exists {
		ch = make(chan int)
		data.workerChans[key] = ch
	}
	go manageWorker(ch, data.statuses[key], worker)
	ch <- startSignal
}

func manageWorker(ch chan int, status *bool, worker routineWorker) {
	for {
		select {
		case signal := <-ch:
			switch signal {
			case startSignal:
				if !*status {
					go worker.start()
				}
			case restartSignal:
				if !*status {
					worker.stop()
				}
				go worker.start()
			case stopSignal:
				worker.stop()
			case killSignal:
				worker.stop()
				return
			}
		}
	}
}

func reloadWorkers(userID string) {
	data := userSpecificDataMap[userID]
	for _, ch := range data.workerChans {
		ch <- restartSignal
	}
}

type twitterDMWorker struct {
	twitterAPI *mybot.TwitterAPI
	stream     *anaconda.Stream
	status     *bool
}

func newTwitterDMWorker(twitterAPI *mybot.TwitterAPI, status *bool) *twitterDMWorker {
	return &twitterDMWorker{twitterAPI, nil, status}
}

func (w *twitterDMWorker) start() {
	if !twitterAPIIsAvailable(w.twitterAPI) {
		return
	}

	*w.status = true
	defer func() { *w.status = false }()

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

func (w *twitterDMWorker) stop() {
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
	status      *bool
}

func newTwitterUserWorker(
	twitterAPI *mybot.TwitterAPI,
	slackAPI *mybot.SlackAPI,
	visionAPI *mybot.VisionAPI,
	languageAPI *mybot.LanguageAPI,
	cache mybot.Cache,
	status *bool,
) *twitterUserWorker {
	return &twitterUserWorker{twitterAPI, slackAPI, visionAPI, languageAPI, cache, nil, status}
}

func (w *twitterUserWorker) start() {
	if !twitterAPIIsAvailable(w.twitterAPI) {
		return
	}

	*w.status = true
	defer func() { *w.status = false }()

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

func (w *twitterUserWorker) stop() {
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
	status *bool,
) *twitterPeriodicWorker {
	return &twitterPeriodicWorker{twitterAPI, slackAPI, visionAPI, languageAPI, cache, config, nil, status}
}

func (w *twitterPeriodicWorker) start() {
	if !twitterAPIIsAvailable(w.twitterAPI) {
		return
	}

	*w.status = true
	defer func() { *w.status = false }()

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

func (w *twitterPeriodicWorker) stop() {
	*w.status = false
}

type slackWorker struct {
	slackAPI    *mybot.SlackAPI
	twitterAPI  *mybot.TwitterAPI
	visionAPI   *mybot.VisionAPI
	languageAPI *mybot.LanguageAPI
	listener    *mybot.SlackListener
	status      *bool
}

func newSlackWorker(
	slackAPI *mybot.SlackAPI,
	twitterAPI *mybot.TwitterAPI,
	visionAPI *mybot.VisionAPI,
	languageAPI *mybot.LanguageAPI,
	status *bool,
) *slackWorker {
	return &slackWorker{slackAPI, twitterAPI, visionAPI, languageAPI, nil, status}
}

func (w *slackWorker) start() {
	if w.slackAPI == nil || !w.slackAPI.Enabled() {
		return
	}

	*w.status = true
	defer func() { *w.status = false }()

	w.listener = w.slackAPI.Listen()
	if err := w.listener.Start(w.visionAPI, w.languageAPI, w.twitterAPI); err != nil {
		log.Print(err)
		return
	}
	log.Print("Failed to keep connection")
}

func (w *slackWorker) stop() {
	if w.listener != nil {
		w.listener.Stop()
	}
}
