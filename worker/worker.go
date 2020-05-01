/*
Package worker provides a way to manipulate concurrent processing.
This guarantees all start/restart/stop operation for worker is always
thread-safe by using Go channel feature.
*/
package worker

import (
	"context"
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
	Apply(ctx context.Context, inChan chan<- WorkerSignal, outChan <-chan interface{}, bufSize int, wg *sync.WaitGroup) (chan<- WorkerSignal, <-chan interface{})
}

// StrategicRestarter is a channel layer applied to restart a worker
// automatically. This restarts a worker when some error happens.
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

func (sr StrategicRestarter) Apply(ctx context.Context, inChan chan<- WorkerSignal, outChan <-chan interface{}, bufSize int, wg *sync.WaitGroup) (chan<- WorkerSignal, <-chan interface{}) {
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
			case <-ctx.Done():
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
	inChan  chan<- WorkerSignal
	outChan <-chan interface{}
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
	status  WorkerStatus
}

func NewWorkerManager(worker models.Worker, bufSize int, layers ...WorkerChannelLayer) *WorkerManager {
	ctx, cancel := context.WithCancel(context.Background())
	inChan, outChan, wg := activateWorkerWithBuffer(ctx, worker, bufSize, layers...)
	return &WorkerManager{inChan, outChan, cancel, wg, StatusStopped}
}

func (wm *WorkerManager) Close() {
	wm.cancel()
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
			h.Handle(out)
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
	Handle(out interface{})
}

// ActivateWorker activates worker, which means worker gets ready to receive
// WorkerSignal to inChan.
// When worker receives WorkerSignal and changes its status, then return
// corresponded WorkerStatus or error via outChan.
// inChan and outChan are created with a given buffer size
func activateWorkerWithBuffer(ctx context.Context, worker models.Worker, bufSize int, layers ...WorkerChannelLayer) (chan<- WorkerSignal, <-chan interface{}, *sync.WaitGroup) {
	inChan, outChan := activateWorker(worker, bufSize)
	wg := &sync.WaitGroup{}
	for _, layer := range layers {
		wg.Add(1)
		inChan, outChan = layer.Apply(ctx, inChan, outChan, bufSize, wg)
	}
	return inChan, outChan, wg
}

func activateWorker(worker models.Worker, bufSize int) (chan<- WorkerSignal, <-chan interface{}) {
	inChan := make(chan WorkerSignal, bufSize)
	outChan := make(chan interface{}, bufSize)

	go func() {
		defer close(outChan)
		wg := &sync.WaitGroup{}
		var cancel context.CancelFunc
		stop := func() {
			if cancel != nil {
				cancel()
				cancel = nil
				wg.Wait()
			}
		}
		defer stop()
		start := func() {
			if cancel == nil {
				var ctx context.Context
				ctx, cancel = context.WithCancel(context.Background())
				wg.Add(1)
				go func(ctx context.Context) {
					startWorkerAndNotify(ctx, worker, outChan)
					wg.Done()
				}(ctx)
			}
		}
		for signal := range inChan {
			// Stop worker
			if signal == RestartSignal || signal == StopSignal {
				stop()
			}
			// Start worker
			if signal == StartSignal || signal == RestartSignal {
				start()
			}
		}
		outChan <- StatusFinished
	}()

	return inChan, outChan
}

func startWorkerAndNotify(ctx context.Context, w models.Worker, outChan chan<- interface{}) {
	outChan <- StatusStarted
	err := w.Start(ctx, outChan)
	if err != nil {
		outChan <- err
	}
	outChan <- StatusStopped
}
