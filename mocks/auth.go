// Automatically generated by MockGen. DO NOT EDIT!
// Source: models/auth.go

package mocks

import (
	gomock "github.com/golang/mock/gomock"
	goth "github.com/markbates/goth"
	http "net/http"
)

// Mock of Authenticator interface
type MockAuthenticator struct {
	ctrl     *gomock.Controller
	recorder *_MockAuthenticatorRecorder
}

// Recorder for MockAuthenticator (not exported)
type _MockAuthenticatorRecorder struct {
	mock *MockAuthenticator
}

func NewMockAuthenticator(ctrl *gomock.Controller) *MockAuthenticator {
	mock := &MockAuthenticator{ctrl: ctrl}
	mock.recorder = &_MockAuthenticatorRecorder{mock}
	return mock
}

func (_m *MockAuthenticator) EXPECT() *_MockAuthenticatorRecorder {
	return _m.recorder
}

func (_m *MockAuthenticator) SetProvider(req *http.Request, name string) {
	_m.ctrl.Call(_m, "SetProvider", req, name)
}

func (_mr *_MockAuthenticatorRecorder) SetProvider(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetProvider", arg0, arg1)
}

func (_m *MockAuthenticator) InitProvider(host string, name string, callback string) {
	_m.ctrl.Call(_m, "InitProvider", host, name, callback)
}

func (_mr *_MockAuthenticatorRecorder) InitProvider(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "InitProvider", arg0, arg1, arg2)
}

func (_m *MockAuthenticator) CompleteUserAuth(provider string, w http.ResponseWriter, r *http.Request) (goth.User, error) {
	ret := _m.ctrl.Call(_m, "CompleteUserAuth", provider, w, r)
	ret0, _ := ret[0].(goth.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAuthenticatorRecorder) CompleteUserAuth(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CompleteUserAuth", arg0, arg1, arg2)
}

func (_m *MockAuthenticator) Logout(provider string, w http.ResponseWriter, r *http.Request) error {
	ret := _m.ctrl.Call(_m, "Logout", provider, w, r)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockAuthenticatorRecorder) Logout(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Logout", arg0, arg1, arg2)
}