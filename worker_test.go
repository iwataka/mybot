package main

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/mybot/mocks"

	"fmt"
	"testing"
)

func TestTwitterPeriodicWorkerStart(t *testing.T) {
	errMsg := "expected error"
	times := 50
	duration := "0.01s"
	id := "id"

	ctrl := gomock.NewController(t)
	runner := mocks.NewMockBatchRunner(ctrl)
	runnerCall := runner.EXPECT().Run().Times(times).Return(nil)
	runner.EXPECT().Run().After(runnerCall).Return(fmt.Errorf(errMsg))
	runner.EXPECT().Verify().Return(nil)
	cache := mocks.NewMockSavable(ctrl)
	cache.EXPECT().Save().Times(times).Return(nil)

	worker := newTwitterPeriodicWorker(runner, cache, duration, id)
	err := worker.Start()
	if err == nil || err.Error() != errMsg {
		t.Fatal("Error not found or not expected error")
	}
}

func TestTwitterPeriodicWorkerStartWithVerificationFalure(t *testing.T) {
	errMsg := "expected error"
	duration := "0.01s"
	id := "id"

	ctrl := gomock.NewController(t)
	runner := mocks.NewMockBatchRunner(ctrl)
	runner.EXPECT().Verify().Return(fmt.Errorf(errMsg))
	cache := mocks.NewMockSavable(ctrl)

	worker := newTwitterPeriodicWorker(runner, cache, duration, id)
	err := worker.Start()
	if err == nil || err.Error() != errMsg {
		t.Fatal("Error not found or not expected error")
	}
}
