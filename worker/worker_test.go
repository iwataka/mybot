package worker_test

import (
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/iwataka/mybot/worker"
)

const (
	timeout         = time.Minute
	timeoutTooSmall = time.Second / 100
)

type testWorker struct {
	count      *int32
	totalCount *int32
	outChan    chan bool
	innerChan  chan bool
}

func newTestWorker() *testWorker {
	var count int32 = 0
	var totalCount int32 = 0
	return &testWorker{&count, &totalCount, make(chan bool), make(chan bool)}
}

func (w *testWorker) Start() error {
	atomic.AddInt32(w.count, 1)
	atomic.AddInt32(w.totalCount, 1)
	defer func() { atomic.AddInt32(w.count, -1) }()
	w.outChan <- true
	<-w.innerChan
	return nil
}

func (w *testWorker) Stop() {
	if *w.count == 1 {
		w.innerChan <- true
	}
}

func (w *testWorker) Name() string {
	return ""
}

func TestKeepSingleWorkerProcessIfMultipleStartSignal(t *testing.T) {
	w := newTestWorker()
	inChan, outChan := ActivateWorker(w, timeout)
	for i := 0; i < 5; i++ {
		inChan <- NewWorkerSignal(StartSignal)
		if i == 0 {
			assertMessage(t, outChan, true)
			<-w.outChan
		}
	}
	assertCount(t, *w.count, 1)
	assertTotalCount(t, *w.totalCount, 1)
}

func TestStopAndStartSignal(t *testing.T) {
	w := newTestWorker()
	inChan, outChan := ActivateWorker(w, timeout)
	var totalCount int32 = 5
	var i int32 = 0
	for ; i < totalCount*2; i++ {
		if i%2 == 0 {
			inChan <- NewWorkerSignal(StartSignal)
			assertMessage(t, outChan, true)
			<-w.outChan
		} else {
			inChan <- NewWorkerSignal(StopSignal)
			assertMessage(t, outChan, false)
		}
	}
	inChan <- NewWorkerSignal(KillSignal)
	assertCount(t, *w.count, 0)
	assertTotalCount(t, *w.totalCount, totalCount)
}

func TestStopSignalForWorker(t *testing.T) {
	w := newTestWorker()
	inChan, outChan := ActivateWorker(w, timeout)
	inChan <- NewWorkerSignal(StartSignal)
	assertMessage(t, outChan, true)
	<-w.outChan
	inChan <- NewWorkerSignal(StopSignal)
	assertMessage(t, outChan, false)
	inChan <- NewWorkerSignal(KillSignal)
	assertCount(t, *w.count, 0)
	assertTotalCount(t, *w.totalCount, 1)
}

func TestRestartSignalForWorker(t *testing.T) {
	w := newTestWorker()
	inChan, outChan := ActivateWorker(w, timeout)
	var totalCount int32 = 5
	var i int32 = 0
	inChan <- NewWorkerSignal(StartSignal)
	assertMessage(t, outChan, true)
	<-w.outChan
	for ; i < totalCount; i++ {
		inChan <- NewWorkerSignal(RestartSignal)
		assertMessage(t, outChan, false)
		assertMessage(t, outChan, true)
		<-w.outChan
	}
	assertCount(t, *w.count, 1)
	assertTotalCount(t, *w.totalCount, totalCount+1)
}

func TestKillSignalForWorker(t *testing.T) {
	w := newTestWorker()
	inChan, outChan := ActivateWorker(w, timeout)
	inChan <- NewWorkerSignal(StartSignal)
	assertMessage(t, outChan, true)
	<-w.outChan
	inChan <- NewWorkerSignal(KillSignal)
	assertMessage(t, outChan, false)
	select {
	case inChan <- NewWorkerSignal(KillSignal):
		t.Fatal("Sent kill signal but worker manager process still wait for signals")
	case <-time.After(time.Second):
	}
	assertCount(t, *w.count, 0)
	assertTotalCount(t, *w.totalCount, 1)
}

func TestWorkerWithoutOutChannel(t *testing.T) {
	w := newTestWorker()
	inChan := ActivateWorkerWithoutOutChan(w, timeoutTooSmall)
	inChan <- NewWorkerSignal(StartSignal)
	<-w.outChan
	inChan <- NewWorkerSignal(RestartSignal)
	<-w.outChan
	inChan <- NewWorkerSignal(StopSignal)
	inChan <- NewWorkerSignal(KillSignal)
	assertCount(t, *w.count, 0)
	assertTotalCount(t, *w.totalCount, 2)
}

func TestWorkerSignalWithOldTimestamp(t *testing.T) {
	w := newTestWorker()
	inChan := ActivateWorkerWithoutOutChan(w, timeoutTooSmall)
	oldRestartSignal := NewWorkerSignal(RestartSignal)
	inChan <- NewWorkerSignal(StartSignal)
	<-w.outChan
	for i := 0; i < 10; i++ {
		inChan <- oldRestartSignal
	}
	assertCount(t, *w.count, 1)
	assertTotalCount(t, *w.totalCount, 1)
}

func TestMultipleRandomWorkerSignals(t *testing.T) {
	w := newTestWorker()
	inChan := ActivateWorkerWithoutOutChan(w, timeoutTooSmall)
	prevSignal := StopSignal
	for i := 0; i < 100; i++ {
		signalSign := rand.Intn(3)
		signal := NewWorkerSignal(signalSign)
		inChan <- signal
		if (signalSign == StartSignal && prevSignal == StopSignal) || signalSign == RestartSignal {
			<-w.outChan
		}
		prevSignal = signalSign
	}
}

func assertMessage(t *testing.T, outChan chan interface{}, expected bool) {
	if msg := <-outChan; msg != expected {
		t.Fatal("Invalid message: ", msg)
	}
}

func assertCount(t *testing.T, count int32, expected int32) {
	if count != expected {
		t.Fatalf("Invalid worker process count: %d (%d expected)", count, expected)
	}
}

func assertTotalCount(t *testing.T, count int32, expected int32) {
	if count != expected {
		t.Fatalf("Invalid worker process total count: %d (%d expected)", count, expected)
	}
}
