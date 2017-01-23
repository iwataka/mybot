package mybot

import (
	"sync"
)

// MybotStatus represents a status of this app
type Status struct {
	twitterListenDMStatus         bool
	twitterListenDMStatusMutex    *sync.Mutex
	twitterListenUsersStatus      bool
	twitterListenUsersStatusMutex *sync.Mutex
	TwitterStatus                 bool
	MonitorConfigStatus           bool
	MonitorTwitterCred            bool
	MonitorGCloudCred             bool
	ServerStatus                  bool
}

func NewStatus() *Status {
	return &Status{
		false,
		new(sync.Mutex),
		false,
		new(sync.Mutex),
		false,
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

func (s *Status) CheckTwitterListenUsersStatus() bool {
	return s.twitterListenUsersStatus
}
