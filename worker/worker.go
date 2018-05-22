/*
Package worker provides a way to manipulate concurrent processing.
This guarantees all start/restart/stop/kill operation for worker is always
thread-safe by using Go channel feature.
*/
package worker

import (
	"time"
)

// These constants indicate signal type sent to worker
const (
	StartSignal = iota
	RestartSignal
	StopSignal
	KillSignal
	PingSignal
)

// WorkerSignal is a signal sent to worker.
// Worker should behave as per the content of it and respond.
type WorkerSignal struct {
	signal    int
	timestamp time.Time
}

// NewWorkerSignal returns a new WorkerSignal with a specified signal type.
func NewWorkerSignal(signal int) *WorkerSignal {
	return &WorkerSignal{signal, time.Now()}
}

// String returns a text indicating a type of this worker signal.
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

func (s WorkerStatus) String() string {
	switch s {
	case StatusAlive:
		return "Alive"
	case StatusFinished:
		return "Finished"
	case StatusKilled:
		return "Killed"
	case StatusRepliedNothing:
		return "Replied Nothing"
	case StatusStarted:
		return "Started"
	case StatusStopped:
		return "Stopped"
	default:
		return ""
	}
}

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
					sendNonBlockingly(outChan, WorkerStatus(StatusRepliedNothing), timeout)
				}
				status = false
			}
			if signal == KillSignal {
				sendNonBlockingly(outChan, WorkerStatus(StatusKilled), timeout)
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
				outChan <- WorkerStatus(StatusAlive)
			}
		}
	}

	sendNonBlockingly(outChan, WorkerStatus(StatusFinished), timeout)
}

func startWorkerAndNotify(w Worker, outChan chan interface{}, ch chan bool, timeout time.Duration) {
	defer func() {
		if outChan != nil {
			sendNonBlockingly(outChan, WorkerStatus(StatusStopped), timeout)
		}
		ch <- true
	}()
	if outChan != nil {
		sendNonBlockingly(outChan, WorkerStatus(StatusStarted), timeout)
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
