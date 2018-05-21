package worker

import (
	"time"
)

const (
	StartSignal = iota
	RestartSignal
	StopSignal
	KillSignal
	PingSignal
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

const (
	StatusAlive = iota
	StatusStarted
	StatusStopped
	StatusKilled
	StatusFinished
	StatusRepliedNothing
)

type WorkerStatus int

type Worker interface {
	Start() error
	Stop()
	Name() string
}

func ActivateWorker(worker Worker, timeout time.Duration) (inChan chan *WorkerSignal, outChan chan interface{}) {
	inChan = make(chan *WorkerSignal)
	outChan = make(chan interface{})
	go activateWorker(inChan, outChan, worker, timeout)
	return
}

func ActivateWorkerWithoutOutChan(worker Worker, timeout time.Duration) (inChan chan *WorkerSignal) {
	inChan = make(chan *WorkerSignal)
	go activateWorker(inChan, nil, worker, timeout)
	return
}

func activateWorker(inChan chan *WorkerSignal, outChan chan interface{}, worker Worker, timeout time.Duration) {
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
				case <-time.After(timeout):
					sendNonBlockingly(outChan, StatusRepliedNothing, timeout)
				}
				status = false
			}
			if signal == KillSignal {
				sendNonBlockingly(outChan, StatusKilled, timeout)
				return
			}
		}
		// Start worker
		if signal == StartSignal || signal == RestartSignal {
			if !status {
				go startWorkerAndNotify(worker, outChan, ch, timeout)
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

	sendNonBlockingly(outChan, StatusFinished, timeout)
}

func startWorkerAndNotify(w Worker, outChan chan interface{}, ch chan bool, timeout time.Duration) {
	defer func() {
		if outChan != nil {
			sendNonBlockingly(outChan, StatusStopped, timeout)
		}
		ch <- true
	}()
	if outChan != nil {
		sendNonBlockingly(outChan, StatusStarted, timeout)
	}
	err := w.Start()
	if err != nil {
		sendNonBlockingly(outChan, err, timeout)
	}
}

func sendNonBlockingly(ch chan interface{}, data interface{}, timeout time.Duration) {
	select {
	case ch <- data:
	case <-time.After(timeout):
		return
	}
}
