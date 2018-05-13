/*
Package worker provides a way to manipulate concurrent processing.
This guarantees all start/restart/stop/kill operation for worker is always
thread-safe by using Go channel feature.
*/
package worker

import (
	"fmt"
	"log"
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

// These constants indicate status type of worker at that time.
const (
	// StatusAlive means enable to respond your signals.
	StatusAlive = iota
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

type Worker interface {
	Start() error
	Stop()
	Name() string
}

func ManageWorker(inChan chan *WorkerSignal, outChan chan interface{}, worker Worker) {
	innerChan := make(chan bool)
	innerStatus := false
	timestamp := time.Now()

	for workerSignal := range inChan {
		signal := workerSignal.signal
		t := workerSignal.timestamp
		// Stop worker
		if signal == RestartSignal || signal == StopSignal || signal == KillSignal {
			force := signal != RestartSignal
			if innerStatus && (force || timestamp.Before(t)) {
				worker.Stop()
				select {
				case <-innerChan:
				case <-time.After(time.Minute):
					log.Printf("Faield to wait stopping worker (timeout: 1m)\n")
				}
				innerStatus = false
			}
			if signal == KillSignal {
				nonBlockingOutput(outChan, fmt.Sprintf("Worker manager for %s killed", worker.Name()))
				return
			}
		}
		// Start worker
		if signal == StartSignal || signal == RestartSignal {
			if !innerStatus {
				go wrapWithStatusManagement(worker.Start, outChan, innerChan)
				innerStatus = true
				timestamp = t
			}
		}
		if signal == PingSignal {
			if outChan != nil {
				outChan <- StatusAlive
			}
		}
	}

	nonBlockingOutput(outChan, fmt.Sprintf("Worker manager for %s finished successfully", worker.Name()))
}

func wrapWithStatusManagement(f func() error, outChan chan interface{}, innerChan chan bool) {
	if outChan != nil {
		nonBlockingOutput(outChan, true)
	}
	defer func() {
		if outChan != nil {
			nonBlockingOutput(outChan, false)
		}
		innerChan <- true
	}()
	err := f()
	if err != nil {
		nonBlockingOutput(outChan, err)
	}
}

func nonBlockingOutput(ch chan interface{}, data interface{}) {
	select {
	case ch <- data:
	case <-time.After(time.Minute):
		log.Println("Failed to send data to outside channel (timeout: 1m)")
	}
}
