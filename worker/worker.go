/*
Package worker provides a way to manipulate concurrent processing.
This guarantees all start/restart/stop operation for worker is always
thread-safe by using Go channel feature.
*/
package worker

import (
	"sync"
	"time"

	"github.com/iwataka/mybot/models"
)

// WorkerSignal is a signal sent to Worker.
// Worker should behave as per the content of it and respond.
type WorkerSignal int

// These constants indicate signal type sent to worker
const (
	StartSignal WorkerSignal = iota
	StopSignal
	RestartSignal
)

// String returns a text indicating a type of this WorkerSignal.
func (s WorkerSignal) String() string {
	switch s {
	case StartSignal:
		return "Start"
	case StopSignal:
		return "Stop"
	case RestartSignal:
		return "Restart"
	default:
		return ""
	}
}

// WorkerStatus is a type indicating Worker status
type WorkerStatus int

// These constants indicate status type of Worker
const (
	StatusStarted WorkerStatus = iota
	StatusStopped
	// StatusFinished means worker was finished successfully.
	StatusFinished
)

// String returns a text to indicating a type of this WorkerStatus.
func (s WorkerStatus) String() string {
	switch s {
	case StatusFinished:
		return "Finished"
	case StatusStarted:
		return "Started"
	case StatusStopped:
		return "Stopped"
	default:
		return ""
	}
}

// WorkerChannelLayer represents a layer for worker channels to catch
// inChan/outChan outputs, apply some filter or conversion and then rethrow
// them.
type WorkerChannelLayer interface {
	// Apply applies this layer to inChan and outChan asynchronously.
	Apply(inChan chan<- WorkerSignal, outChan <-chan interface{}, bufSize int, stopChan <-chan interface{}, wg *sync.WaitGroup) (chan<- WorkerSignal, <-chan interface{})
}

// StrategicRestarter is a channel layer applied to restart a worker
// automatically. This restarts a worker when some error happenes.
// If error happens more than `count` times in `interval` duration, this stops.
type StrategicRestarter struct {
	interval      time.Duration
	count         int
	suppressError bool
}

// NewStrategicRestarter creates a new StrategicRestarter.
func NewStrategicRestarter(interval time.Duration, count int, suppressError bool) StrategicRestarter {
	return StrategicRestarter{interval, count, suppressError}
}

func (sr StrategicRestarter) Apply(inChan chan<- WorkerSignal, outChan <-chan interface{}, bufSize int, stopChan <-chan interface{}, wg *sync.WaitGroup) (chan<- WorkerSignal, <-chan interface{}) {
	oc := make(chan interface{}, bufSize)
	go func() {
		_wg := &sync.WaitGroup{}
		errTimestamps := []time.Time{}
		for {
			select {
			case msg := <-outChan:
				switch msg.(type) {
				case error:
					ts := time.Now()
					errTimestamps = append(errTimestamps, ts)
					if len(errTimestamps) > sr.count {
						errTimestamps = errTimestamps[len(errTimestamps)-sr.count:]
					}
					if len(errTimestamps) < sr.count || ts.Sub(errTimestamps[0]).Nanoseconds() > sr.interval.Nanoseconds() {
						_wg.Add(1)
						go func() {
							inChan <- RestartSignal
							_wg.Done()
						}()
						if sr.suppressError {
							continue
						}
					}
				}
				oc <- msg
			case <-stopChan:
				go func() {
					defer close(oc)
					for msg := range outChan {
						oc <- msg
					}
				}()
				goto done
			}
		}
	done:
		_wg.Wait()
		wg.Done()
	}()
	return inChan, oc
}

type WorkerManager struct {
	inChan   chan<- WorkerSignal
	outChan  <-chan interface{}
	stopChan chan<- interface{}
	wg       *sync.WaitGroup
	status   WorkerStatus
}

func NewWorkerManager(worker models.Worker, bufSize int, layers ...WorkerChannelLayer) *WorkerManager {
	inChan, outChan, stopChan, wg := activateWorkerWithBuffer(worker, bufSize, layers...)
	return &WorkerManager{inChan, outChan, stopChan, wg, StatusStopped}
}

func (wm *WorkerManager) Close() {
	close(wm.stopChan)
	wm.wg.Wait()
	close(wm.inChan)
}

func (wm *WorkerManager) Send(s WorkerSignal) {
	wm.inChan <- s
}

func (wm *WorkerManager) Receive() interface{} {
	out := <-wm.outChan
	wm.setStatus(out)
	return out
}

func (wm *WorkerManager) HandleOutput(h WorkerManagerOutHandler) {
	for out := range wm.outChan {
		wm.setStatus(out)
		if h != nil {
			switch o := out.(type) {
			case WorkerStatus:
				h.HandleWorkerStatus(o)
			case error:
				h.HandleError(o)
			}
		}
	}
}

func (wm *WorkerManager) setStatus(o interface{}) {
	if s, ok := o.(WorkerStatus); ok {
		wm.status = s
	}
}

func (wm *WorkerManager) Status() WorkerStatus {
	return wm.status
}

type WorkerManagerOutHandler interface {
	HandleWorkerStatus(s WorkerStatus)
	HandleError(err error)
}

// ActivateWorker activates worker, which means worker gets ready to receive
// WorkerSignal to inChan.
// When worker receives WorkerSignal and changes its status, then return
// corresponded WorkerStatus or error via outChan.
// inChan and outChan are created with a given buffer size
func activateWorkerWithBuffer(worker models.Worker, bufSize int, layers ...WorkerChannelLayer) (chan<- WorkerSignal, <-chan interface{}, chan<- interface{}, *sync.WaitGroup) {
	inChan, outChan := activateWorker(worker, bufSize)
	stopChan := make(chan interface{})
	wg := &sync.WaitGroup{}
	for _, layer := range layers {
		wg.Add(1)
		inChan, outChan = layer.Apply(inChan, outChan, bufSize, stopChan, wg)
	}
	return inChan, outChan, stopChan, wg
}

func activateWorker(worker models.Worker, bufSize int) (chan<- WorkerSignal, <-chan interface{}) {
	inChan := make(chan WorkerSignal, bufSize)
	outChan := make(chan interface{}, bufSize)

	go func() {
		defer close(outChan)
		waitStopChan := make(chan bool)
		defer close(waitStopChan)
		var stopChan chan interface{}
		stop := func() {
			if stopChan != nil {
				close(stopChan)
				stopChan = nil
				<-waitStopChan
			}
		}
		defer stop()
		for signal := range inChan {
			// Stop worker
			if signal == RestartSignal || signal == StopSignal {
				stop()
			}

			// Start worker
			if signal == StartSignal || signal == RestartSignal {
				if stopChan == nil {
					stopChan = make(chan interface{})
					go func(stopChan <-chan interface{}) {
						startWorkerAndNotify(worker, stopChan, outChan)
						waitStopChan <- true
					}(stopChan)
				}
			}
		}

		outChan <- StatusFinished
	}()

	return inChan, outChan
}

func startWorkerAndNotify(w models.Worker, stopChan <-chan interface{}, outChan chan<- interface{}) {
	outChan <- StatusStarted
	err := w.Start(stopChan)
	if err != nil {
		outChan <- err
	}
	outChan <- StatusStopped
}
