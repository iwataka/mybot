// Automatically generated by MockGen. DO NOT EDIT!
// Source: runner/batch.go

package mocks

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of BatchRunner interface
type MockBatchRunner struct {
	ctrl     *gomock.Controller
	recorder *_MockBatchRunnerRecorder
}

// Recorder for MockBatchRunner (not exported)
type _MockBatchRunnerRecorder struct {
	mock *MockBatchRunner
}

func NewMockBatchRunner(ctrl *gomock.Controller) *MockBatchRunner {
	mock := &MockBatchRunner{ctrl: ctrl}
	mock.recorder = &_MockBatchRunnerRecorder{mock}
	return mock
}

func (_m *MockBatchRunner) EXPECT() *_MockBatchRunnerRecorder {
	return _m.recorder
}

func (_m *MockBatchRunner) Run() error {
	ret := _m.ctrl.Call(_m, "Run")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockBatchRunnerRecorder) Run() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Run")
}

func (_m *MockBatchRunner) IsAvailable() error {
	ret := _m.ctrl.Call(_m, "IsAvailable")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockBatchRunnerRecorder) IsAvailable() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "IsAvailable")
}
