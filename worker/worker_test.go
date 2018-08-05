package worker_test

import (
	"errors"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/mybot/mocks"
	. "github.com/iwataka/mybot/worker"
	"github.com/stretchr/testify/require"
)

const (
	timeout = time.Minute
)

// testWorker is a simple worker which just counts how many times this worker
// has started and records whether this worker is started or stopped.
type testWorker struct {
	// count is &1 if this is started, otherwise &0.
	count *int32
	// totalCount is how many times this has started.
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
	// To notify Start() processing is finished.
	w.outChan <- true
	<-w.innerChan
	return nil
}

func (w *testWorker) Stop() error {
	if *w.count == 1 {
		w.innerChan <- true
	}
	return nil
}

func (w *testWorker) Name() string {
	return ""
}

func TestKeepSingleWorkerProcessAsItIsWhenMultipleStartSignalSent(t *testing.T) {
	w := newTestWorker()
	inChan, outChan := ActivateWorker(w, timeout)
	for i := 0; i < 5; i++ {
		inChan <- NewWorkerSignal(StartSignal)
		if i == 0 {
			require.Equal(t, StatusStarted, <-outChan)
			<-w.outChan
		}
	}
	require.EqualValues(t, 1, *w.count)
	require.EqualValues(t, 1, *w.totalCount)
}

func TestStopAndStartSignalSentAlternately(t *testing.T) {
	w := newTestWorker()
	inChan, outChan := ActivateWorker(w, timeout)
	var totalCount int32 = 5
	var i int32 = 0
	for ; i < totalCount*2; i++ {
		if i%2 == 0 {
			inChan <- NewWorkerSignal(StartSignal)
			require.Equal(t, StatusStarted, <-outChan)
			<-w.outChan
		} else {
			inChan <- NewWorkerSignal(StopSignal)
			require.Equal(t, StatusStopped, <-outChan)
		}
	}
	inChan <- NewWorkerSignal(KillSignal)
	require.EqualValues(t, 0, *w.count)
	require.EqualValues(t, totalCount, *w.totalCount)
}

func TestStopSignal(t *testing.T) {
	w := newTestWorker()
	inChan, outChan := ActivateWorker(w, timeout)
	inChan <- NewWorkerSignal(StartSignal)
	require.Equal(t, StatusStarted, <-outChan)
	<-w.outChan
	inChan <- NewWorkerSignal(StopSignal)
	require.Equal(t, StatusStopped, <-outChan)
	inChan <- NewWorkerSignal(KillSignal)
	require.EqualValues(t, 0, *w.count)
	require.EqualValues(t, 1, *w.totalCount)
}

func TestRestartSignalForWorker(t *testing.T) {
	w := newTestWorker()
	inChan, outChan := ActivateWorker(w, timeout)
	var totalCount int32 = 5
	var i int32 = 0
	inChan <- NewWorkerSignal(StartSignal)
	require.Equal(t, StatusStarted, <-outChan)
	<-w.outChan
	for ; i < totalCount; i++ {
		inChan <- NewWorkerSignal(RestartSignal)
		require.Equal(t, StatusStopped, <-outChan)
		require.Equal(t, StatusStarted, <-outChan)
		<-w.outChan
	}
	require.EqualValues(t, 1, *w.count)
	require.EqualValues(t, totalCount+1, *w.totalCount)
}

func TestKillSignalForWorker(t *testing.T) {
	w := newTestWorker()
	inChan, outChan := ActivateWorker(w, timeout)
	inChan <- NewWorkerSignal(StartSignal)
	require.Equal(t, StatusStarted, <-outChan)
	<-w.outChan
	inChan <- NewWorkerSignal(KillSignal)
	require.Equal(t, StatusStopped, <-outChan)
	select {
	case inChan <- NewWorkerSignal(KillSignal):
		t.Fatal("Sent kill signal but worker manager process still wait for signals")
	case <-time.After(time.Second):
	}
	require.EqualValues(t, 0, *w.count)
	require.EqualValues(t, 1, *w.totalCount)
}

func TestWorkerWithoutOutChannel(t *testing.T) {
	w := newTestWorker()
	inChan := ActivateWorkerWithoutOutChan(w, timeout)
	inChan <- NewWorkerSignal(StartSignal)
	<-w.outChan
	inChan <- NewWorkerSignal(RestartSignal)
	<-w.outChan
	inChan <- NewWorkerSignal(StopSignal)
	inChan <- NewWorkerSignal(KillSignal)
	require.EqualValues(t, 0, *w.count)
	require.EqualValues(t, 2, *w.totalCount)
}

func TestWorkerSignalWithOldTimestamp(t *testing.T) {
	w := newTestWorker()
	inChan := ActivateWorkerWithoutOutChan(w, timeout)
	oldRestartSignal := NewWorkerSignal(RestartSignal)
	inChan <- NewWorkerSignal(StartSignal)
	<-w.outChan
	for i := 0; i < 10; i++ {
		inChan <- oldRestartSignal
	}
	require.EqualValues(t, 1, *w.count)
	require.EqualValues(t, 1, *w.totalCount)
}

func TestMultipleRandomWorkerSignals(t *testing.T) {
	w := newTestWorker()
	inChan := ActivateWorkerWithoutOutChan(w, timeout)
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

func TestPingSignal(t *testing.T) {
	w := newTestWorker()
	inChan, outChan := ActivateWorker(w, timeout)
	inChan <- NewWorkerSignal(PingSignal)
	require.Equal(t, StatusActive, <-outChan)
	require.EqualValues(t, 0, *w.count)
	require.EqualValues(t, 0, *w.totalCount)
}

func TestWorkerFinished(t *testing.T) {
	w := newTestWorker()
	inChan, outChan := ActivateWorker(w, timeout)
	close(inChan)
	require.EqualValues(t, StatusFinished, <-outChan)
}

func TestStartSignalWhenWorkerStartFuncThrowAnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	w := mocks.NewMockWorker(ctrl)
	err := errors.New("foo")
	w.EXPECT().Start().Return(err)
	inChan, outChan := ActivateWorker(w, timeout)
	inChan <- NewWorkerSignal(StartSignal)
	require.Equal(t, StatusStarted, <-outChan)
	require.Equal(t, err, <-outChan)
	require.Equal(t, StatusStopped, <-outChan)
}
