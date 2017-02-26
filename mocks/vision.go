// Automatically generated by MockGen. DO NOT EDIT!
// Source: lib/vision.go

package mocks

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/iwataka/mybot/models"
)

// Mock of VisionMatcher interface
type MockVisionMatcher struct {
	ctrl     *gomock.Controller
	recorder *_MockVisionMatcherRecorder
}

// Recorder for MockVisionMatcher (not exported)
type _MockVisionMatcherRecorder struct {
	mock *MockVisionMatcher
}

func NewMockVisionMatcher(ctrl *gomock.Controller) *MockVisionMatcher {
	mock := &MockVisionMatcher{ctrl: ctrl}
	mock.recorder = &_MockVisionMatcherRecorder{mock}
	return mock
}

func (_m *MockVisionMatcher) EXPECT() *_MockVisionMatcherRecorder {
	return _m.recorder
}

func (_m *MockVisionMatcher) MatchImages(_param0 []string, _param1 *models.VisionCondition) ([]string, []bool, error) {
	ret := _m.ctrl.Call(_m, "MatchImages", _param0, _param1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].([]bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockVisionMatcherRecorder) MatchImages(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MatchImages", arg0, arg1)
}

func (_m *MockVisionMatcher) Enabled() bool {
	ret := _m.ctrl.Call(_m, "Enabled")
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockVisionMatcherRecorder) Enabled() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Enabled")
}
