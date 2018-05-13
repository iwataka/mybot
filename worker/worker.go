package worker

import (
	"fmt"
	"time"
)

const (
	StartSignal = iota
	RestartSignal
	StopSignal
	KillSignal
	PingSignal
)

const (
	StatusAlive = iota
)

type WorkerSignal struct {
	signal    int
	timestamp time.Time
}

func NewWorkerSignal(signal int) *WorkerSignal {
	return &WorkerSignal{signal, time.Now()}
}

func (s WorkerSignal) String() string {
	switch s.signal {
	case StartSignal:
		return "Start"
	case RestartSignal:
		return "Restart"
	case StopSignal:
		return "Stop"
	case KillSignal:
		return "Kill"
	case PingSignal:
		return "Ping"
	default:
		return ""
	}
}

type Worker interface {
	Start() error
	Stop()
	Name() string
}

func ActivateWorker(inChan chan *WorkerSignal, outChan chan interface{}, worker Worker) {
	go activateWorker(inChan, outChan, worker)
}

func activateWorker(inChan chan *WorkerSignal, outChan chan interface{}, worker Worker) {
	ch := make(chan bool)
	status := false
	timestamp := time.Now()

	for workerSignal := range inChan {
		signal := workerSignal.signal
		t := workerSignal.timestamp
		// Stop worker
		if signal == RestartSignal || signal == StopSignal || signal == KillSignal {
			force := signal != RestartSignal
			if status && (force || timestamp.Before(t)) {
				worker.Stop()
				select {
				case <-ch:
				case <-time.After(time.Minute):
					msg := fmt.Sprintf("Faield to wait stopping worker %s (timeout: 1m)", worker.Name())
					sendNonBlockingly(outChan, msg)
				}
				status = false
			}
			if signal == KillSignal {
				sendNonBlockingly(outChan, fmt.Sprintf("Worker manager for %s killed", worker.Name()))
				return
			}
		}
		// Start worker
		if signal == StartSignal || signal == RestartSignal {
			if !status {
				go startWorkerAndNotify(worker, outChan, ch)
				status = true
				timestamp = t
			}
		}
		if signal == PingSignal {
			if outChan != nil {
				outChan <- StatusAlive
			}
		}
	}

	sendNonBlockingly(outChan, fmt.Sprintf("Worker manager for %s finished successfully", worker.Name()))
}

func startWorkerAndNotify(w Worker, outChan chan interface{}, ch chan bool) {
	defer func() {
		if outChan != nil {
			sendNonBlockingly(outChan, false)
		}
		ch <- true
	}()
	if outChan != nil {
		sendNonBlockingly(outChan, true)
	}
	err := w.Start()
	if err != nil {
		sendNonBlockingly(outChan, err)
	}
}

func sendNonBlockingly(ch chan interface{}, data interface{}) {
	select {
	case ch <- data:
	case <-time.After(time.Minute):
		return
	}
}
