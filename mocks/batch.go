// Code generated by MockGen. DO NOT EDIT.
// Source: runner/batch.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	data "github.com/iwataka/mybot/data"
	reflect "reflect"
)

// MockBatchRunner is a mock of BatchRunner interface
type MockBatchRunner struct {
	ctrl     *gomock.Controller
	recorder *MockBatchRunnerMockRecorder
}

// MockBatchRunnerMockRecorder is the mock recorder for MockBatchRunner
type MockBatchRunnerMockRecorder struct {
	mock *MockBatchRunner
}

// NewMockBatchRunner creates a new mock instance
func NewMockBatchRunner(ctrl *gomock.Controller) *MockBatchRunner {
	mock := &MockBatchRunner{ctrl: ctrl}
	mock.recorder = &MockBatchRunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBatchRunner) EXPECT() *MockBatchRunnerMockRecorder {
	return m.recorder
}

// Run mocks base method
func (m *MockBatchRunner) Run() ([]interface{}, []data.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run")
	ret0, _ := ret[0].([]interface{})
	ret1, _ := ret[1].([]data.Action)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Run indicates an expected call of Run
func (mr *MockBatchRunnerMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockBatchRunner)(nil).Run))
}

// IsAvailable mocks base method
func (m *MockBatchRunner) IsAvailable() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAvailable")
	ret0, _ := ret[0].(error)
	return ret0
}

// IsAvailable indicates an expected call of IsAvailable
func (mr *MockBatchRunnerMockRecorder) IsAvailable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAvailable", reflect.TypeOf((*MockBatchRunner)(nil).IsAvailable))
}
