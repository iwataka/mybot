// Automatically generated by MockGen. DO NOT EDIT!
// Source: lib/language.go

package mocks

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/iwataka/mybot/models"
)

// Mock of LanguageMatcher interface
type MockLanguageMatcher struct {
	ctrl     *gomock.Controller
	recorder *_MockLanguageMatcherRecorder
}

// Recorder for MockLanguageMatcher (not exported)
type _MockLanguageMatcherRecorder struct {
	mock *MockLanguageMatcher
}

func NewMockLanguageMatcher(ctrl *gomock.Controller) *MockLanguageMatcher {
	mock := &MockLanguageMatcher{ctrl: ctrl}
	mock.recorder = &_MockLanguageMatcherRecorder{mock}
	return mock
}

func (_m *MockLanguageMatcher) EXPECT() *_MockLanguageMatcherRecorder {
	return _m.recorder
}

func (_m *MockLanguageMatcher) MatchText(_param0 string, _param1 *models.LanguageCondition) (string, bool, error) {
	ret := _m.ctrl.Call(_m, "MatchText", _param0, _param1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockLanguageMatcherRecorder) MatchText(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MatchText", arg0, arg1)
}

func (_m *MockLanguageMatcher) Enabled() bool {
	ret := _m.ctrl.Call(_m, "Enabled")
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockLanguageMatcherRecorder) Enabled() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Enabled")
}