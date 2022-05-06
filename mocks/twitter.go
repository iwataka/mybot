// Code generated by MockGen. DO NOT EDIT.
// Source: models/twitter.go

// Package mocks is a generated GoMock package.
package mocks

import (
	url "net/url"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	anaconda "github.com/iwataka/anaconda"
)

// MockTwitterAPI is a mock of TwitterAPI interface.
type MockTwitterAPI struct {
	ctrl     *gomock.Controller
	recorder *MockTwitterAPIMockRecorder
}

// MockTwitterAPIMockRecorder is the mock recorder for MockTwitterAPI.
type MockTwitterAPIMockRecorder struct {
	mock *MockTwitterAPI
}

// NewMockTwitterAPI creates a new mock instance.
func NewMockTwitterAPI(ctrl *gomock.Controller) *MockTwitterAPI {
	mock := &MockTwitterAPI{ctrl: ctrl}
	mock.recorder = &MockTwitterAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTwitterAPI) EXPECT() *MockTwitterAPIMockRecorder {
	return m.recorder
}

// AddEntryToCollection mocks base method.
func (m *MockTwitterAPI) AddEntryToCollection(arg0 string, arg1 int64, arg2 url.Values) (anaconda.CollectionEntryAddResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddEntryToCollection", arg0, arg1, arg2)
	ret0, _ := ret[0].(anaconda.CollectionEntryAddResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddEntryToCollection indicates an expected call of AddEntryToCollection.
func (mr *MockTwitterAPIMockRecorder) AddEntryToCollection(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEntryToCollection", reflect.TypeOf((*MockTwitterAPI)(nil).AddEntryToCollection), arg0, arg1, arg2)
}

// CreateCollection mocks base method.
func (m *MockTwitterAPI) CreateCollection(arg0 string, arg1 url.Values) (anaconda.CollectionShowResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCollection", arg0, arg1)
	ret0, _ := ret[0].(anaconda.CollectionShowResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCollection indicates an expected call of CreateCollection.
func (mr *MockTwitterAPIMockRecorder) CreateCollection(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCollection", reflect.TypeOf((*MockTwitterAPI)(nil).CreateCollection), arg0, arg1)
}

// Favorite mocks base method.
func (m *MockTwitterAPI) Favorite(arg0 int64) (anaconda.Tweet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Favorite", arg0)
	ret0, _ := ret[0].(anaconda.Tweet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Favorite indicates an expected call of Favorite.
func (mr *MockTwitterAPIMockRecorder) Favorite(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Favorite", reflect.TypeOf((*MockTwitterAPI)(nil).Favorite), arg0)
}

// GetCollectionListByUserId mocks base method.
func (m *MockTwitterAPI) GetCollectionListByUserId(arg0 int64, arg1 url.Values) (anaconda.CollectionListResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCollectionListByUserId", arg0, arg1)
	ret0, _ := ret[0].(anaconda.CollectionListResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCollectionListByUserId indicates an expected call of GetCollectionListByUserId.
func (mr *MockTwitterAPIMockRecorder) GetCollectionListByUserId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCollectionListByUserId", reflect.TypeOf((*MockTwitterAPI)(nil).GetCollectionListByUserId), arg0, arg1)
}

// GetFavorites mocks base method.
func (m *MockTwitterAPI) GetFavorites(arg0 url.Values) ([]anaconda.Tweet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavorites", arg0)
	ret0, _ := ret[0].([]anaconda.Tweet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavorites indicates an expected call of GetFavorites.
func (mr *MockTwitterAPIMockRecorder) GetFavorites(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavorites", reflect.TypeOf((*MockTwitterAPI)(nil).GetFavorites), arg0)
}

// GetFriendsList mocks base method.
func (m *MockTwitterAPI) GetFriendsList(arg0 url.Values) (anaconda.UserCursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFriendsList", arg0)
	ret0, _ := ret[0].(anaconda.UserCursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFriendsList indicates an expected call of GetFriendsList.
func (mr *MockTwitterAPIMockRecorder) GetFriendsList(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFriendsList", reflect.TypeOf((*MockTwitterAPI)(nil).GetFriendsList), arg0)
}

// GetSearch mocks base method.
func (m *MockTwitterAPI) GetSearch(arg0 string, arg1 url.Values) (anaconda.SearchResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSearch", arg0, arg1)
	ret0, _ := ret[0].(anaconda.SearchResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSearch indicates an expected call of GetSearch.
func (mr *MockTwitterAPIMockRecorder) GetSearch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSearch", reflect.TypeOf((*MockTwitterAPI)(nil).GetSearch), arg0, arg1)
}

// GetSelf mocks base method.
func (m *MockTwitterAPI) GetSelf(arg0 url.Values) (anaconda.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSelf", arg0)
	ret0, _ := ret[0].(anaconda.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSelf indicates an expected call of GetSelf.
func (mr *MockTwitterAPIMockRecorder) GetSelf(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSelf", reflect.TypeOf((*MockTwitterAPI)(nil).GetSelf), arg0)
}

// GetUserSearch mocks base method.
func (m *MockTwitterAPI) GetUserSearch(arg0 string, arg1 url.Values) ([]anaconda.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserSearch", arg0, arg1)
	ret0, _ := ret[0].([]anaconda.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserSearch indicates an expected call of GetUserSearch.
func (mr *MockTwitterAPIMockRecorder) GetUserSearch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserSearch", reflect.TypeOf((*MockTwitterAPI)(nil).GetUserSearch), arg0, arg1)
}

// GetUserTimeline mocks base method.
func (m *MockTwitterAPI) GetUserTimeline(arg0 url.Values) ([]anaconda.Tweet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserTimeline", arg0)
	ret0, _ := ret[0].([]anaconda.Tweet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserTimeline indicates an expected call of GetUserTimeline.
func (mr *MockTwitterAPIMockRecorder) GetUserTimeline(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserTimeline", reflect.TypeOf((*MockTwitterAPI)(nil).GetUserTimeline), arg0)
}

// GetUsersLookup mocks base method.
func (m *MockTwitterAPI) GetUsersLookup(arg0 string, arg1 url.Values) ([]anaconda.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersLookup", arg0, arg1)
	ret0, _ := ret[0].([]anaconda.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersLookup indicates an expected call of GetUsersLookup.
func (mr *MockTwitterAPIMockRecorder) GetUsersLookup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersLookup", reflect.TypeOf((*MockTwitterAPI)(nil).GetUsersLookup), arg0, arg1)
}

// GetUsersShow mocks base method.
func (m *MockTwitterAPI) GetUsersShow(arg0 string, arg1 url.Values) (anaconda.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersShow", arg0, arg1)
	ret0, _ := ret[0].(anaconda.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersShow indicates an expected call of GetUsersShow.
func (mr *MockTwitterAPIMockRecorder) GetUsersShow(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersShow", reflect.TypeOf((*MockTwitterAPI)(nil).GetUsersShow), arg0, arg1)
}

// PostDMToScreenName mocks base method.
func (m *MockTwitterAPI) PostDMToScreenName(arg0, arg1 string) (anaconda.DirectMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostDMToScreenName", arg0, arg1)
	ret0, _ := ret[0].(anaconda.DirectMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostDMToScreenName indicates an expected call of PostDMToScreenName.
func (mr *MockTwitterAPIMockRecorder) PostDMToScreenName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostDMToScreenName", reflect.TypeOf((*MockTwitterAPI)(nil).PostDMToScreenName), arg0, arg1)
}

// PostTweet mocks base method.
func (m *MockTwitterAPI) PostTweet(arg0 string, arg1 url.Values) (anaconda.Tweet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostTweet", arg0, arg1)
	ret0, _ := ret[0].(anaconda.Tweet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostTweet indicates an expected call of PostTweet.
func (mr *MockTwitterAPIMockRecorder) PostTweet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostTweet", reflect.TypeOf((*MockTwitterAPI)(nil).PostTweet), arg0, arg1)
}

// PublicStreamFilter mocks base method.
func (m *MockTwitterAPI) PublicStreamFilter(arg0 url.Values) *anaconda.Stream {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublicStreamFilter", arg0)
	ret0, _ := ret[0].(*anaconda.Stream)
	return ret0
}

// PublicStreamFilter indicates an expected call of PublicStreamFilter.
func (mr *MockTwitterAPIMockRecorder) PublicStreamFilter(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublicStreamFilter", reflect.TypeOf((*MockTwitterAPI)(nil).PublicStreamFilter), arg0)
}

// Retweet mocks base method.
func (m *MockTwitterAPI) Retweet(arg0 int64, arg1 bool) (anaconda.Tweet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Retweet", arg0, arg1)
	ret0, _ := ret[0].(anaconda.Tweet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Retweet indicates an expected call of Retweet.
func (mr *MockTwitterAPIMockRecorder) Retweet(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Retweet", reflect.TypeOf((*MockTwitterAPI)(nil).Retweet), arg0, arg1)
}

// SetLogger mocks base method.
func (m *MockTwitterAPI) SetLogger(arg0 anaconda.Logger) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLogger", arg0)
}

// SetLogger indicates an expected call of SetLogger.
func (mr *MockTwitterAPIMockRecorder) SetLogger(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLogger", reflect.TypeOf((*MockTwitterAPI)(nil).SetLogger), arg0)
}

// UserStream mocks base method.
func (m *MockTwitterAPI) UserStream(arg0 url.Values) *anaconda.Stream {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserStream", arg0)
	ret0, _ := ret[0].(*anaconda.Stream)
	return ret0
}

// UserStream indicates an expected call of UserStream.
func (mr *MockTwitterAPIMockRecorder) UserStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserStream", reflect.TypeOf((*MockTwitterAPI)(nil).UserStream), arg0)
}

// VerifyCredentials mocks base method.
func (m *MockTwitterAPI) VerifyCredentials() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyCredentials")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyCredentials indicates an expected call of VerifyCredentials.
func (mr *MockTwitterAPIMockRecorder) VerifyCredentials() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyCredentials", reflect.TypeOf((*MockTwitterAPI)(nil).VerifyCredentials))
}
