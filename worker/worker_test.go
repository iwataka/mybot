package worker_test

import (
	"context"
	"errors"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/mybot/mocks"
	"github.com/iwataka/mybot/worker"
	"github.com/stretchr/testify/require"
)

// testWorker is a simple worker which just counts how many times this worker
// has started and records whether this worker is started or stopped.
type testWorker struct {
	// count is &1 if this is started, otherwise &0.
	count *int32
	// totalCount is how many times this has started.
	totalCount *int32
}

func newTestWorker() *testWorker {
	var count int32 = 0
	var totalCount int32 = 0
	return &testWorker{&count, &totalCount}
}

func (w *testWorker) Start(ctx context.Context, outChan chan<- interface{}) error {
	atomic.AddInt32(w.count, 1)
	atomic.AddInt32(w.totalCount, 1)
	defer func() { atomic.AddInt32(w.count, -1) }()
	// To notify Start() processing is finished.
	outChan <- true
	<-ctx.Done()
	return nil
}

func (w *testWorker) Name() string {
	return ""
}

func TestKeepSingleWorkerProcessAsItIsWhenMultipleStartSignalSent(t *testing.T) {
	w := newTestWorker()
	wm := worker.NewWorkerManager(w, 0)
	defer wm.Close()
	for i := 0; i < 5; i++ {
		wm.Send(worker.StartSignal)
		if i == 0 {
			checkStatus(t, worker.StatusStarted, wm)
			require.Equal(t, true, wm.Receive())
		}
	}
	require.EqualValues(t, 1, *w.count)
	require.EqualValues(t, 1, *w.totalCount)
}

func TestStopAndStartSignalSentAlternately(t *testing.T) {
	w := newTestWorker()
	wm := worker.NewWorkerManager(w, 0)
	defer wm.Close()
	var totalCount int32 = 5
	var i int32 = 0
	for ; i < totalCount*2; i++ {
		if i%2 == 0 {
			wm.Send(worker.StartSignal)
			checkStatus(t, worker.StatusStarted, wm)
			require.Equal(t, true, wm.Receive())
		} else {
			wm.Send(worker.StopSignal)
			checkStatus(t, worker.StatusStopped, wm)
		}
	}
	require.EqualValues(t, 0, *w.count)
	require.EqualValues(t, totalCount, *w.totalCount)
}

func TestStopSignal(t *testing.T) {
	w := newTestWorker()
	wm := worker.NewWorkerManager(w, 0)
	defer wm.Close()
	wm.Send(worker.StartSignal)
	checkStatus(t, worker.StatusStarted, wm)
	require.Equal(t, true, wm.Receive())
	wm.Send(worker.StopSignal)
	checkStatus(t, worker.StatusStopped, wm)
	require.EqualValues(t, 0, *w.count)
	require.EqualValues(t, 1, *w.totalCount)
}

func TestRestartSignalForWorker(t *testing.T) {
	w := newTestWorker()
	wm := worker.NewWorkerManager(w, 0)
	defer wm.Close()
	var totalCount int32 = 5
	var i int32 = 0
	wm.Send(worker.StartSignal)
	checkStatus(t, worker.StatusStarted, wm)
	require.Equal(t, true, wm.Receive())
	for ; i < totalCount; i++ {
		wm.Send(worker.RestartSignal)
		checkStatus(t, worker.StatusStopped, wm)
		checkStatus(t, worker.StatusStarted, wm)
		require.Equal(t, true, wm.Receive())
	}
	require.EqualValues(t, 1, *w.count)
	require.EqualValues(t, totalCount+1, *w.totalCount)
}

func TestMultipleRandomWorkerSignals(t *testing.T) {
	w := newTestWorker()
	wm := worker.NewWorkerManager(w, 0)
	defer wm.Close()
	isActive := false
	for i := 0; i < 100; i++ {
		signalSign := rand.Intn(3)
		signal := worker.WorkerSignal(signalSign)
		wm.Send(signal)
		started := signal == worker.RestartSignal || (!isActive && signal == worker.StartSignal)
		stopped := isActive && (signal == worker.StopSignal || signal == worker.RestartSignal)
		if stopped {
			checkStatus(t, worker.StatusStopped, wm)
		}
		if started {
			checkStatus(t, worker.StatusStarted, wm)
			require.Equal(t, true, wm.Receive())
		}
		isActive = signal != worker.StopSignal
	}
}

func TestWorkerFinished(t *testing.T) {
	w := newTestWorker()
	wm := worker.NewWorkerManager(w, 0)
	wm.Close()
	require.EqualValues(t, worker.StatusFinished, wm.Receive())
}

func TestStartSignalWhenWorkerStartFuncThrowAnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	w := mocks.NewMockWorker(ctrl)
	err := errors.New("foo")
	w.EXPECT().Start(gomock.Any(), gomock.Any()).Return(err)
	wm := worker.NewWorkerManager(w, 0)
	defer wm.Close()
	wm.Send(worker.StartSignal)
	checkStatus(t, worker.StatusStarted, wm)
	require.Equal(t, err, wm.Receive())
	checkStatus(t, worker.StatusStopped, wm)
}

func TestStrategicRestarter(t *testing.T) {
	testStrategicRestarter(t, true)
	testStrategicRestarter(t, false)
}

func testStrategicRestarter(t *testing.T, suppressError bool) {
	ctrl := gomock.NewController(t)
	w := mocks.NewMockWorker(ctrl)
	err := errors.New("error")
	w.EXPECT().Start(gomock.Any(), gomock.Any()).Times(5).Return(err)
	interval, _ := time.ParseDuration("60m")
	l := worker.NewStrategicRestarter(interval, 5, suppressError)
	wm := worker.NewWorkerManager(w, 0, l)
	defer wm.Close()
	wm.Send(worker.StartSignal)
	for i := 0; i < 5; i++ {
		checkStatus(t, worker.StatusStarted, wm)
		if i < 4 {
			if !suppressError {
				require.Equal(t, err, wm.Receive())
			}
		} else {
			require.Equal(t, err, wm.Receive())
		}
		checkStatus(t, worker.StatusStopped, wm)
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
	w.EXPECT().Start(gomock.Any(), gomock.Any()).Times(7).Return(err)
	w.EXPECT().Start(gomock.Any(), gomock.Any()).Times(1).Return(nil)
	interval, _ := time.ParseDuration("0ns")
	l := worker.NewStrategicRestarter(interval, 5, suppressError)
	wm := worker.NewWorkerManager(w, 0, l)
	defer wm.Close()
	wm.Send(worker.StartSignal)
	for i := 0; i < 7; i++ {
		checkStatus(t, worker.StatusStarted, wm)
		if !suppressError {
			require.Equal(t, err, wm.Receive())
		}
		checkStatus(t, worker.StatusStopped, wm)
	}
	checkStatus(t, worker.StatusStarted, wm)
	checkStatus(t, worker.StatusStopped, wm)
}

func checkStatus(t *testing.T, s worker.WorkerStatus, wm *worker.WorkerManager) {
	require.Equal(t, s, wm.Receive())
	require.Equal(t, s, wm.Status())
}
