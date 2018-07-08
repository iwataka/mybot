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

// WorkerSignal is a signal sent to Worker.
// Worker should behave as per the content of it and respond.
type WorkerSignal struct {
	signal    int
	timestamp time.Time
}

// NewWorkerSignal returns a new WorkerSignal with a specified signal type.
func NewWorkerSignal(signal int) *WorkerSignal {
	return &WorkerSignal{signal, time.Now()}
}

// String returns a text indicating a type of this WorkerSignal.
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

// WorkerStatus is a type indicating Worker status
type WorkerStatus int

// These constants indicate status type of Worker
const (
	StatusActive WorkerStatus = iota
	// StatusInactive means worker is inactive, unable to respond to
	// WorkerSignal.
	StatusInactive WorkerStatus = iota
	StatusStarted  WorkerStatus = iota
	StatusStopped  WorkerStatus = iota
	// StatusKilled means worker was finished forcefully.
	StatusKilled WorkerStatus = iota
	// StatusFinished means worker was finished successfully.
	StatusFinished WorkerStatus = iota
)

// String returns a text to indicating a type of this WorkerStatus.
func (s WorkerStatus) String() string {
	switch s {
	case StatusActive:
		return "Active"
	case StatusFinished:
		return "Finished"
	case StatusKilled:
		return "Killed"
	case StatusInactive:
		return "Inactive"
	case StatusStarted:
		return "Started"
	case StatusStopped:
		return "Stopped"
	default:
		return ""
	}
}

// Worker is worker which has its own operation and provides APIs to start/stop
// it.
type Worker interface {
	Start() error
	Stop() error
	// Name returns a name of this Worker, to distinguish this from others.
	Name() string
}

// ActivateWorker activates worker, which means worker gets ready to receive
// WorkerSignal to inChan. When worker receives WorkerSignal and changes its
// status, then return corresponded WorkerStatus or error via outChan.
func ActivateWorker(worker Worker, timeout time.Duration) (inChan chan *WorkerSignal, outChan chan interface{}) {
	inChan = make(chan *WorkerSignal)
	outChan = make(chan interface{})
	go activateWorker(inChan, outChan, worker, timeout)
	return
}

// ActivateWorkerWithoutOutChan is almost same as ActivateWorker but doesn't
// use outChan.
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
				err := worker.Stop()
				if err != nil {
					sendNonBlockingly(outChan, err, timeout)
				}
				select {
				case <-ch:
				case <-time.After(timeout):
					sendNonBlockingly(outChan, StatusInactive, timeout)
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
				outChan <- StatusActive
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
