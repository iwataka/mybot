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

type RoutineWorker interface {
	Start() error
	Stop()
	Name() string
}

func ManageWorker(inChan chan *WorkerSignal, outChan chan interface{}, worker RoutineWorker) {
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
					msg := fmt.Sprintf("Faield to wait stopping worker %s (timeout: 1m)", worker.Name())
					nonBlockingOutput(outChan, msg)
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
		return
	}
}
