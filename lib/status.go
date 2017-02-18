package mybot

import (
	"github.com/iwataka/anaconda"
	"sync"
)

// MybotStatus represents a status of this app
type Status struct {
	twitterListenDMStatus         bool
	twitterListenDMStatusMutex    *sync.Mutex
	twitterListenUsersStatus      bool
	twitterListenUsersStatusMutex *sync.Mutex
	slackListen                   bool
	slackListenMutex              *sync.Mutex
	TwitterStatus                 bool
	ServerStatus                  bool
	PassTwitterApp                bool
	PassTwitterAuth               bool
}

func NewStatus() *Status {
	return &Status{
		false,
		new(sync.Mutex),
		false,
		new(sync.Mutex),
		false,
		new(sync.Mutex),
		false,
		false,
		false,
		false,
	}
}

func (s *Status) LockListenDMRoutine() {
	s.twitterListenDMStatusMutex.Lock()
	s.twitterListenDMStatus = true
}

func (s *Status) UnlockListenDMRoutine() {
	s.twitterListenDMStatus = false
	s.twitterListenDMStatusMutex.Unlock()
}

func (s *Status) CheckTwitterListenDMStatus() bool {
	return s.twitterListenDMStatus
}

func (s *Status) LockListenUsersRoutine() {
	s.twitterListenUsersStatusMutex.Lock()
	s.twitterListenUsersStatus = true
}

func (s *Status) UnlockListenUsersRoutine() {
	s.twitterListenUsersStatus = false
	s.twitterListenUsersStatusMutex.Unlock()
}

func (s *Status) CheckSlackListen() bool {
	return s.slackListen
}

func (s *Status) LockSlackListenRoutine() {
	s.slackListenMutex.Lock()
	s.slackListen = true
}

func (s *Status) UnlockSlackListenRoutine() {
	s.slackListen = false
	s.slackListenMutex.Unlock()
}

func (s *Status) CheckTwitterListenUsersStatus() bool {
	return s.twitterListenUsersStatus
}

func (s *Status) UpdateTwitterAuth(api *TwitterAPI) {
	s.PassTwitterApp = anaconda.GetConsumerKey() != "" &&
		anaconda.GetConsumerSecret() != ""
	if s.PassTwitterApp {
		ok, err := api.VerifyCredentials()
		s.PassTwitterAuth = ok && err == nil
	} else {
		s.PassTwitterAuth = false
	}
}
