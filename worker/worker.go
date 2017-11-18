package worker

import (
	"log"
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

	start := func(t time.Time) {
		if !innerStatus {
			go wrapWithStatusManagement(worker.Start, outChan, innerChan)
			innerStatus = true
			timestamp = t
		}
	}

	stop := func(t time.Time, force bool) {
		if innerStatus && (force || timestamp.Before(t)) {
			worker.Stop()
			<-innerChan
			innerStatus = false
		}
	}

	defer func() {
		close(innerChan)
	}()

	for workerSignal := range inChan {
		signal := workerSignal.signal
		t := workerSignal.timestamp
		switch signal {
		case StartSignal:
			start(t)
		case RestartSignal:
			stop(t, false)
			start(t)
		case StopSignal:
			stop(t, true)
		case KillSignal:
			stop(t, true)
			return
		case PingSignal:
			if outChan != nil {
				outChan <- StatusAlive
			}
		}
	}
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
