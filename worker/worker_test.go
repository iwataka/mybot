package worker

import (
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
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

	for {
		select {
		case <-w.innerChan:
			return nil
		}
	}
	return nil
}

func (w *testWorker) Stop() {
	w.innerChan <- true
}

func (w *testWorker) Name() string {
	return ""
}

func TestKeepSingleWorkerProcessIfMultipleStartSignal(t *testing.T) {
	w := newTestWorker()
	inChan := make(chan *WorkerSignal)
	outChan := make(chan interface{})
	go ManageWorker(inChan, outChan, w)
	defer func() { inChan <- NewWorkerSignal(KillSignal); close(w.outChan) }()
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
	inChan := make(chan *WorkerSignal)
	outChan := make(chan interface{})
	go ManageWorker(inChan, outChan, w)
	defer func() { close(w.outChan) }()
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
	inChan := make(chan *WorkerSignal)
	outChan := make(chan interface{})
	go ManageWorker(inChan, outChan, w)
	defer func() { close(w.outChan) }()
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
	inChan := make(chan *WorkerSignal)
	outChan := make(chan interface{})
	go ManageWorker(inChan, outChan, w)
	defer func() { inChan <- NewWorkerSignal(KillSignal); close(w.outChan) }()
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
	inChan := make(chan *WorkerSignal)
	outChan := make(chan interface{})
	go ManageWorker(inChan, outChan, w)
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Channel is still open because it receives signals successfully")
		}
		assertCount(t, *w.count, 0)
		assertTotalCount(t, *w.totalCount, 1)
		close(w.outChan)
	}()
	inChan <- NewWorkerSignal(StartSignal)
	assertMessage(t, outChan, true)
	<-w.outChan
	inChan <- NewWorkerSignal(KillSignal)
	assertMessage(t, outChan, false)
	inChan <- NewWorkerSignal(KillSignal)
}

func TestWorkerWithoutOutChannel(t *testing.T) {
	w := newTestWorker()
	inChan := make(chan *WorkerSignal)
	go ManageWorker(inChan, nil, w)
	defer func() { close(w.outChan) }()
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
	inChan := make(chan *WorkerSignal)
	go ManageWorker(inChan, nil, w)
	defer func() { inChan <- NewWorkerSignal(KillSignal); close(w.outChan) }()
	inChan <- NewWorkerSignal(StartSignal)
	<-w.outChan
	oldRestartSignal := &WorkerSignal{RestartSignal, time.Now().Add(-1 * time.Hour)}
	for i := 0; i < 10; i++ {
		inChan <- oldRestartSignal
	}
	assertCount(t, *w.count, 1)
	assertTotalCount(t, *w.totalCount, 1)
}

func TestMultipleRandomWorkerSignals(t *testing.T) {
	w := newTestWorker()
	inChan := make(chan *WorkerSignal)
	go ManageWorker(inChan, nil, w)
	defer func() { inChan <- NewWorkerSignal(KillSignal); close(w.outChan) }()
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
