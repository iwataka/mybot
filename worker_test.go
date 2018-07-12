package main

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/mybot/lib"
	"github.com/iwataka/mybot/mocks"
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/runner"
	"github.com/iwataka/mybot/utils"
	"github.com/iwataka/mybot/worker"
	"github.com/stretchr/testify/assert"

	"fmt"
	"testing"
	"time"
)

func TestTwitterPeriodicWorkerStart(t *testing.T) {
	errMsg := "expected error"
	times := 50
	duration := "0.01s"
	id := "id"
	worker := generatePeriodicWorker(t, times, duration, id, fmt.Errorf(errMsg), nil)
	err := worker.Start()
	assert.Error(t, err)
	assert.Equal(t, errMsg, err.Error())
}

func TestManageTwitterPeriodicWorker(t *testing.T) {
	times := -1
	duration := "0.01s"
	id := "id"
	w := generatePeriodicWorker(t, times, duration, id, nil, nil)
	h := generateWorkerMessageHandler(t, -1)

	key := 0
	workerChans := make(map[int]chan *worker.WorkerSignal)
	statuses := make(map[int]bool)
	statuses[key] = false
	activateWorkerAndStart(key, workerChans, statuses, w, h)
	workerChans[key] <- worker.NewWorkerSignal(worker.RestartSignal)
	workerChans[key] <- worker.NewWorkerSignal(worker.RestartSignal)
	workerChans[key] <- worker.NewWorkerSignal(worker.KillSignal)
}

func TestManageTwitterPeriodicWorkerWithVerificationFailure(t *testing.T) {
	errMsg := "expected error"
	times := -1
	duration := "0.01s"
	id := "id"
	w := generatePeriodicWorker(t, times, duration, id, fmt.Errorf(errMsg), fmt.Errorf(errMsg))
	h := generateWorkerMessageHandler(t, -1)

	key := 0
	workerChans := make(map[int]chan *worker.WorkerSignal)
	statuses := make(map[int]bool)
	statuses[key] = false
	activateWorkerAndStart(key, workerChans, statuses, w, h)
	workerChans[key] <- worker.NewWorkerSignal(worker.RestartSignal)
	workerChans[key] <- worker.NewWorkerSignal(worker.RestartSignal)
	workerChans[key] <- worker.NewWorkerSignal(worker.KillSignal)
}

func TestTwitterPeriodicWorkerStartWithVerificationFalure(t *testing.T) {
	errMsg := "expected error"
	times := 1
	duration := "0.01s"
	id := "id"
	w := generatePeriodicWorker(t, times, duration, id, fmt.Errorf(errMsg), fmt.Errorf(errMsg))

	err := w.Start()
	assert.Error(t, err)
	assert.Equal(t, errMsg, err.Error())
}

// TODO: Call gomocl.Controlller#Finish to ensure all stub methods are called.
func generatePeriodicWorker(t *testing.T, times int, duration string, id string, runErr error, verifyErr error) *twitterPeriodicWorker {
	ctrl := gomock.NewController(t)
	runner := generateRunner(ctrl, times, runErr, verifyErr)
	cache := generateCache(ctrl, times)
	config := generateConfig(t, duration)
	return newTwitterPeriodicWorker(runner, cache, config, time.Second, id)
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

func generateConfig(t *testing.T, duration string) mybot.Config {
	config := mybot.NewTestFileConfig("", t)
	config.SetPollingDuration(duration)
	return config
}

func generateWorkerMessageHandler(t *testing.T, times int) models.WorkerMessageHandler {
	ctrl := gomock.NewController(t)
	h := mocks.NewMockWorkerMessageHandler(ctrl)
	if times < 0 {
		h.EXPECT().Handle(gomock.Any()).AnyTimes().Return(nil)
	} else {
		h.EXPECT().Handle(gomock.Any()).Times(times).Return(nil)
	}
	return h
}
