package worker

import (
	"time"
)

const (
	StartSignal = iota
	RestartSignal
	StopSignal
	KillSignal
)

type WorkerSignal struct {
	signal    int
	timestamp time.Time
}

func NewWorkerSignal(signal int) *WorkerSignal {
	return &WorkerSignal{signal, time.Now()}
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

	clean := func() {
		close(innerChan)
		close(inChan)
		if outChan != nil {
			close(outChan)
		}
	}

	for {
		select {
		case workerSignal := <-inChan:
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
				clean()
				return
			}
		}
	}
}

func wrapWithStatusManagement(f func() error, outChan chan interface{}, ch chan bool) {
	if outChan != nil {
		outChan <- true
	}
	defer func() {
		if outChan != nil {
			outChan <- false
		}
		ch <- true
	}()
	err := f()
	if err != nil {
		outChan <- err
	}
}
