package main

import (
	gomock "github.com/golang/mock/gomock"
	mybot "github.com/iwataka/mybot/lib"
	"github.com/iwataka/mybot/mocks"
	worker "github.com/iwataka/mybot/worker"

	"testing"
	"time"
)

func TestTwitterPeriodicWorkerStart(t *testing.T) {
	errMsg := "expected error"
	times := 50
	duration := "0.01s"
	id := "id"
	worker := generatePeriodicWorker(t, times, duration, id, mybot.Errorf(errMsg), nil)
	err := worker.Start()
	if err == nil || err.Error() != errMsg {
		t.Fatal("Error not found or not expected error")
	}
}

func TestManageTwitterPeriodicWorker(t *testing.T) {
	times := -1
	duration := "0.01s"
	id := "id"
	w := generatePeriodicWorker(t, times, duration, id, nil, nil)

	key := 0
	workerChans := make(map[int]chan *worker.WorkerSignal)
	statuses := make(map[int]*bool)
	status := false
	statuses[key] = &status
	manageWorkerWithStart(key, workerChans, statuses, w)
	workerChans[key] <- worker.NewWorkerSignal(worker.RestartSignal)
	workerChans[key] <- worker.NewWorkerSignal(worker.RestartSignal)
	workerChans[key] <- worker.NewWorkerSignal(worker.KillSignal)
}

func TestManageTwitterPeriodicWorkerWithVerificationFailure(t *testing.T) {
	errMsg := "expected error"
	times := -1
	duration := "0.01s"
	id := "id"
	w := generatePeriodicWorker(t, times, duration, id, mybot.Errorf(errMsg), mybot.Errorf(errMsg))

	key := 0
	workerChans := make(map[int]chan *worker.WorkerSignal)
	statuses := make(map[int]*bool)
	status := false
	statuses[key] = &status
	manageWorkerWithStart(key, workerChans, statuses, w)
	workerChans[key] <- worker.NewWorkerSignal(worker.RestartSignal)
	workerChans[key] <- worker.NewWorkerSignal(worker.RestartSignal)
	workerChans[key] <- worker.NewWorkerSignal(worker.KillSignal)
}

func TestTwitterPeriodicWorkerStartWithVerificationFalure(t *testing.T) {
	errMsg := "expected error"
	times := 1
	duration := "0.01s"
	id := "id"
	w := generatePeriodicWorker(t, times, duration, id, mybot.Errorf(errMsg), mybot.Errorf(errMsg))

	err := w.Start()
	if err == nil || err.Error() != errMsg {
		t.Fatal("Error not found or not expected error")
	}
}

func generatePeriodicWorker(t *testing.T, times int, duration string, id string, runErr error, verifyErr error) *twitterPeriodicWorker {
	ctrl := gomock.NewController(t)
	runner := generateRunner(ctrl, times, runErr, verifyErr)
	cache := generateCache(ctrl, times)
	return newTwitterPeriodicWorker(runner, cache, duration, time.Second, id)
}

func generateRunner(ctrl *gomock.Controller, times int, runErr error, verifyErr error) mybot.BatchRunner {
	runner := mocks.NewMockBatchRunner(ctrl)
	var runCall *gomock.Call
	if times < 0 {
		runCall = runner.EXPECT().Run().AnyTimes().Return(nil)
	} else {
		runCall = runner.EXPECT().Run().Times(times).Return(nil)
	}
	runner.EXPECT().Run().After(runCall).Return(runErr)
	if times < 0 {
		runner.EXPECT().Verify().AnyTimes().Return(verifyErr)
	} else {
		runner.EXPECT().Verify().Times(times).Return(verifyErr)
	}
	return runner
}

func generateCache(ctrl *gomock.Controller, times int) mybot.Savable {
	cache := mocks.NewMockSavable(ctrl)
	if times < 0 {
		cache.EXPECT().Save().AnyTimes().Return(nil)
	} else {
		cache.EXPECT().Save().Times(times).Return(nil)
	}
	return cache
}
