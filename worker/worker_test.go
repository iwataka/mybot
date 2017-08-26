package worker

import (
	"sync/atomic"
	"testing"
	"time"
)

type testWorker struct {
	status     bool
	count      *int32
	totalCount *int32
	outChan    chan bool
}

func newTestWorker() *testWorker {
	var count int32 = 0
	var totalCount int32 = 0
	outChan := make(chan bool)
	return &testWorker{false, &count, &totalCount, outChan}
}

func (w *testWorker) Start() {
	atomic.AddInt32(w.count, 1)
	atomic.AddInt32(w.totalCount, 1)
	defer func() { atomic.AddInt32(w.count, -1) }()
	w.status = true
	w.outChan <- true

	for w.status {
		time.Sleep(1 * time.Second)
	}
}

func (w *testWorker) Stop() {
	w.status = false
}

func TestKeepSingleWorkerProcessIfMultipleStartSignal(t *testing.T) {
	w := newTestWorker()
	inChan := make(chan int)
	status := false
	go ManageWorker(inChan, &status, w)
	defer func() { inChan <- KillSignal; close(w.outChan) }()
	for i := 0; i < 5; i++ {
		inChan <- StartSignal
		if i == 0 {
			<-w.outChan
		}
	}
	inChan <- FlushSignal
	if *w.count != 1 {
		t.Fatal("Invalid worker process count: ", *w.count)
	}
	if *w.totalCount != 1 {
		t.Fatal("Invalid worker process total count: ", *w.totalCount)
	}
	if !status {
		t.Fatal("Invalid worker process status: ", status)
	}
}

func TestStopAndStartSignal(t *testing.T) {
	w := newTestWorker()
	inChan := make(chan int)
	status := false
	go ManageWorker(inChan, &status, w)
	defer func() { inChan <- KillSignal; close(w.outChan) }()
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			inChan <- StartSignal
			<-w.outChan
		} else {
			inChan <- StopSignal
		}
	}
	inChan <- FlushSignal
	if *w.count != 0 {
		t.Fatal("Invalid worker process count: ", *w.count)
	}
	if *w.totalCount != 5 {
		t.Fatal("Invalid worker process total count: ", *w.totalCount)
	}
	if status {
		t.Fatal("Invalid worker process status: ", status)
	}
}

func TestStopSignalForWorker(t *testing.T) {
	w := newTestWorker()
	inChan := make(chan int)
	status := false
	go ManageWorker(inChan, &status, w)
	defer func() { inChan <- KillSignal; close(w.outChan) }()
	inChan <- StartSignal
	<-w.outChan
	inChan <- StopSignal
	inChan <- FlushSignal
	if *w.count != 0 {
		t.Fatal("Invalid worker process count: ", *w.count)
	}
	if *w.totalCount != 1 {
		t.Fatal("Invalid worker process total count: ", *w.totalCount)
	}
	if status {
		t.Fatal("Invalid worker process status: ", status)
	}
}

func TestRestartSignalForWorker(t *testing.T) {
	w := newTestWorker()
	inChan := make(chan int)
	status := false
	go ManageWorker(inChan, &status, w)
	defer func() { inChan <- KillSignal; close(w.outChan) }()
	var totalCount int32 = 2
	var i int32 = 0
	for ; i < totalCount; i++ {
		inChan <- RestartSignal
		<-w.outChan
	}
	inChan <- FlushSignal
	if *w.count != 1 {
		t.Fatal("Invalid worker process count: ", *w.count)
	}
	if *w.totalCount != totalCount {
		t.Fatal("Invalid worker process total count: ", *w.totalCount)
	}
	if !status {
		t.Fatal("Invalid worker process status: ", status)
	}
}

func TestKillSignalForWorker(t *testing.T) {
	w := newTestWorker()
	inChan := make(chan int)
	status := false
	go ManageWorker(inChan, &status, w)
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Channel is still open because it receives signals successfully")
		}
		if *w.count != 0 {
			t.Fatal("Invalid worker process count: ", *w.count)
		}
		if *w.totalCount != 1 {
			t.Fatal("Invalid worker process total count: ", *w.totalCount)
		}
		if status {
			t.Fatal("Invalid worker process status: ", status)
		}
		close(w.outChan)
	}()
	inChan <- StartSignal
	<-w.outChan
	inChan <- KillSignal
	inChan <- KillSignal
	inChan <- FlushSignal
}
