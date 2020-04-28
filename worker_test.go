package main

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/mybot/core"
	"github.com/iwataka/mybot/mocks"
	"github.com/iwataka/mybot/runner"
	"github.com/iwataka/mybot/utils"
	"github.com/iwataka/mybot/worker"
	"github.com/stretchr/testify/require"

	"fmt"
	"testing"
)

func TestTwitterPeriodicWorker_Start(t *testing.T) {
	errMsg := "expected error"
	times := 50
	duration := "0.01s"
	id := "id"
	worker := generatePeriodicWorker(t, times, duration, id, fmt.Errorf(errMsg), nil)
	ch := make(chan interface{})
	defer close(ch)
	err := worker.Start(ch)
	require.Error(t, err)
	require.Equal(t, errMsg, err.Error())
}

func TestManageTwitterPeriodicWorker(t *testing.T) {
	times := -1
	duration := "0.01s"
	id := "id"
	w := generatePeriodicWorker(t, times, duration, id, nil, nil)

	key := 0
	workerMgrs := make(map[int]*worker.WorkerManager)
	activateWorkerAndStart(key, workerMgrs, w, nil, defaultWorkerBufSize)
	defer workerMgrs[key].Close()
	workerMgrs[key].Send(worker.RestartSignal)
	workerMgrs[key].Send(worker.RestartSignal)
}

func TestManageTwitterPeriodicWorkerWithVerificationFailure(t *testing.T) {
	errMsg := "expected error"
	times := -1
	duration := "0.01s"
	id := "id"
	w := generatePeriodicWorker(t, times, duration, id, fmt.Errorf(errMsg), fmt.Errorf(errMsg))

	key := 0
	workerMgrs := make(map[int]*worker.WorkerManager)
	activateWorkerAndStart(key, workerMgrs, w, nil, defaultWorkerBufSize)
	defer workerMgrs[key].Close()
	workerMgrs[key].Send(worker.RestartSignal)
	workerMgrs[key].Send(worker.RestartSignal)
}

func TestTwitterPeriodicWorkerStartWithVerificationFalure(t *testing.T) {
	errMsg := "expected error"
	times := 1
	duration := "0.01s"
	id := "id"
	w := generatePeriodicWorker(t, times, duration, id, fmt.Errorf(errMsg), fmt.Errorf(errMsg))

	ch := make(chan interface{})
	defer close(ch)
	err := w.Start(ch)
	require.Error(t, err)
	require.Equal(t, errMsg, err.Error())
}

// TODO: Call gomocl.Controlller#Finish to ensure all stub methods are called.
func generatePeriodicWorker(t *testing.T, times int, duration string, id string, runErr error, verifyErr error) *twitterPeriodicWorker {
	ctrl := gomock.NewController(t)
	runner := generateRunner(ctrl, times, runErr, verifyErr)
	cache := generateCache(ctrl, times)
	config := generateConfig(t, duration)
	return newTwitterPeriodicWorker(runner, cache, config, id)
}

func generateRunner(ctrl *gomock.Controller, times int, runErr error, verifyErr error) runner.BatchRunner {
	runner := mocks.NewMockBatchRunner(ctrl)
	var runCall *gomock.Call
	if times < 0 {
		runCall = runner.EXPECT().Run().AnyTimes().Return(nil)
	} else {
		runCall = runner.EXPECT().Run().Times(times).Return(nil)
	}
	runner.EXPECT().Run().After(runCall).Return(runErr)
	if times < 0 {
		runner.EXPECT().IsAvailable().AnyTimes().Return(verifyErr)
	} else {
		runner.EXPECT().IsAvailable().Times(times).Return(verifyErr)
	}
	return runner
}

func generateCache(ctrl *gomock.Controller, times int) utils.Savable {
	cache := mocks.NewMockSavable(ctrl)
	if times < 0 {
		cache.EXPECT().Save().AnyTimes().Return(nil)
	} else {
		cache.EXPECT().Save().Times(times).Return(nil)
	}
	return cache
}

func generateConfig(t *testing.T, duration string) core.Config {
	config := core.NewTestFileConfig("", t)
	config.SetPollingDuration(duration)
	return config
}
