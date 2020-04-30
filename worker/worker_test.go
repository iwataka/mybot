package worker

import (
	"context"
	"errors"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/mybot/mocks"
	"github.com/stretchr/testify/require"
)

// testWorker is a simple worker which just counts how many times this worker
// has started and records whether this worker is started or stopped.
type testWorker struct {
	// count is &1 if this is started, otherwise &0.
	count *int32
	// totalCount is how many times this has started.
	totalCount *int32
	outChan    chan bool
}

func newTestWorker() *testWorker {
	var count int32 = 0
	var totalCount int32 = 0
	return &testWorker{&count, &totalCount, make(chan bool)}
}

func (w *testWorker) Start(ctx context.Context) error {
	atomic.AddInt32(w.count, 1)
	atomic.AddInt32(w.totalCount, 1)
	defer func() { atomic.AddInt32(w.count, -1) }()
	// To notify Start() processing is finished.
	w.outChan <- true
	<-ctx.Done()
	return nil
}

func (w *testWorker) Name() string {
	return ""
}

func TestKeepSingleWorkerProcessAsItIsWhenMultipleStartSignalSent(t *testing.T) {
	w := newTestWorker()
	wm := NewWorkerManager(w, 0)
	defer wm.Close()
	for i := 0; i < 5; i++ {
		wm.Send(StartSignal)
		if i == 0 {
			checkStatus(t, StatusStarted, wm)
			<-w.outChan
		}
	}
	require.EqualValues(t, 1, *w.count)
	require.EqualValues(t, 1, *w.totalCount)
}

func TestStopAndStartSignalSentAlternately(t *testing.T) {
	w := newTestWorker()
	wm := NewWorkerManager(w, 0)
	defer wm.Close()
	var totalCount int32 = 5
	var i int32 = 0
	for ; i < totalCount*2; i++ {
		if i%2 == 0 {
			wm.Send(StartSignal)
			checkStatus(t, StatusStarted, wm)
			<-w.outChan
		} else {
			wm.Send(StopSignal)
			checkStatus(t, StatusStopped, wm)
		}
	}
	require.EqualValues(t, 0, *w.count)
	require.EqualValues(t, totalCount, *w.totalCount)
}

func TestStopSignal(t *testing.T) {
	w := newTestWorker()
	wm := NewWorkerManager(w, 0)
	defer wm.Close()
	wm.Send(StartSignal)
	checkStatus(t, StatusStarted, wm)
	<-w.outChan
	wm.Send(StopSignal)
	checkStatus(t, StatusStopped, wm)
	require.EqualValues(t, 0, *w.count)
	require.EqualValues(t, 1, *w.totalCount)
}

func TestRestartSignalForWorker(t *testing.T) {
	w := newTestWorker()
	wm := NewWorkerManager(w, 0)
	defer wm.Close()
	var totalCount int32 = 5
	var i int32 = 0
	wm.Send(StartSignal)
	checkStatus(t, StatusStarted, wm)
	<-w.outChan
	for ; i < totalCount; i++ {
		wm.Send(RestartSignal)
		checkStatus(t, StatusStopped, wm)
		checkStatus(t, StatusStarted, wm)
		<-w.outChan
	}
	require.EqualValues(t, 1, *w.count)
	require.EqualValues(t, totalCount+1, *w.totalCount)
}

func TestMultipleRandomWorkerSignals(t *testing.T) {
	w := newTestWorker()
	wm := NewWorkerManager(w, 0)
	defer wm.Close()
	isActive := false
	for i := 0; i < 100; i++ {
		signalSign := rand.Intn(3)
		signal := WorkerSignal(signalSign)
		wm.Send(signal)
		started := signal == RestartSignal || (!isActive && signal == StartSignal)
		stopped := isActive && (signal == StopSignal || signal == RestartSignal)
		if stopped {
			checkStatus(t, StatusStopped, wm)
		}
		if started {
			checkStatus(t, StatusStarted, wm)
			<-w.outChan
		}
		isActive = signal != StopSignal
	}
}

func TestWorkerFinished(t *testing.T) {
	w := newTestWorker()
	wm := NewWorkerManager(w, 0)
	wm.Close()
	require.EqualValues(t, StatusFinished, wm.Receive())
}

func TestStartSignalWhenWorkerStartFuncThrowAnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	w := mocks.NewMockWorker(ctrl)
	err := errors.New("foo")
	w.EXPECT().Start(gomock.Any()).Return(err)
	wm := NewWorkerManager(w, 0)
	defer wm.Close()
	wm.Send(StartSignal)
	checkStatus(t, StatusStarted, wm)
	require.Equal(t, err, wm.Receive())
	checkStatus(t, StatusStopped, wm)
}

func TestStrategicRestarter(t *testing.T) {
	testStrategicRestarter(t, true)
	testStrategicRestarter(t, false)
}

func testStrategicRestarter(t *testing.T, suppressError bool) {
	ctrl := gomock.NewController(t)
	w := mocks.NewMockWorker(ctrl)
	err := errors.New("error")
	w.EXPECT().Start(gomock.Any()).Times(5).Return(err)
	interval, _ := time.ParseDuration("60m")
	l := NewStrategicRestarter(interval, 5, suppressError)
	wm := NewWorkerManager(w, 0, l)
	defer wm.Close()
	wm.Send(StartSignal)
	for i := 0; i < 5; i++ {
		checkStatus(t, StatusStarted, wm)
		if i < 4 {
			if !suppressError {
				require.Equal(t, err, wm.Receive())
			}
		} else {
			require.Equal(t, err, wm.Receive())
		}
		checkStatus(t, StatusStopped, wm)
	}
}

func TestStrategicRestarterWithSmallInterval(t *testing.T) {
	testStrategicRestarterWithSmallInterval(t, true)
	testStrategicRestarterWithSmallInterval(t, false)
}

func testStrategicRestarterWithSmallInterval(t *testing.T, suppressError bool) {
	ctrl := gomock.NewController(t)
	w := mocks.NewMockWorker(ctrl)
	err := errors.New("error")
	w.EXPECT().Start(gomock.Any()).Times(7).Return(err)
	w.EXPECT().Start(gomock.Any()).Times(1).Return(nil)
	interval, _ := time.ParseDuration("0ns")
	l := NewStrategicRestarter(interval, 5, suppressError)
	wm := NewWorkerManager(w, 0, l)
	defer wm.Close()
	wm.Send(StartSignal)
	for i := 0; i < 7; i++ {
		checkStatus(t, StatusStarted, wm)
		if !suppressError {
			require.Equal(t, err, wm.Receive())
		}
		checkStatus(t, StatusStopped, wm)
	}
	checkStatus(t, StatusStarted, wm)
	checkStatus(t, StatusStopped, wm)
}

func checkStatus(t *testing.T, s WorkerStatus, wm *WorkerManager) {
	require.Equal(t, s, wm.Receive())
	require.Equal(t, s, wm.Status())
}
