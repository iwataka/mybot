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
	TwitterStatus                 bool

	MonitorStatus map[string]bool
	MonitorChans  map[string][]chan bool
	MonitorMutex  map[string]*sync.Mutex

	ServerStatus bool

	PassTwitterApp  bool
	PassTwitterAuth bool
}

func NewStatus() *Status {
	return &Status{
		false,
		new(sync.Mutex),
		false,
		new(sync.Mutex),
		false,

		make(map[string]bool),
		make(map[string][]chan bool),
		make(map[string]*sync.Mutex),

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

func (s *Status) initStatus(file string) {
	_, exists := s.MonitorStatus[file]
	if !exists {
		s.MonitorStatus[file] = false
	}
}

func (s *Status) GetMonitorStatus(file string) bool {
	s.initStatus(file)
	return s.MonitorStatus[file]
}

func (s *Status) SetMonitorStatus(file string, flag bool) {
	s.initStatus(file)
	s.MonitorStatus[file] = flag
}

func (s *Status) initMutex(file string) {
	_, exists := s.MonitorMutex[file]
	if !exists {
		s.MonitorMutex[file] = new(sync.Mutex)
	}
}

func (s *Status) LockMonitor(file string) {
	s.initMutex(file)
	s.MonitorMutex[file].Lock()
}

func (s *Status) UnlockMonitor(file string) {
	s.initMutex(file)
	s.MonitorMutex[file].Unlock()
}

func (s *Status) initChans(file string) {
	_, exists := s.MonitorChans[file]
	if !exists {
		s.MonitorChans[file] = []chan bool{}
	}
}

func (s *Status) AddMonitorChan(file string, c chan bool) {
	s.initChans(file)
	s.LockMonitor(file)
	s.MonitorChans[file] = append(s.MonitorChans[file], c)
	s.UnlockMonitor(file)
}

func (s *Status) SendToMonitor(file string, flag bool) {
	s.initChans(file)
	for _, c := range s.MonitorChans[file] {
		c <- flag
	}
	s.initChans(file)
}

func (s *Status) UpdateTwitterAuth(api *TwitterAPI) {
	if anaconda.GetConsumerKey() != "" && anaconda.GetConsumerSecret() != "" {
		s.PassTwitterApp = true
	}
	ok, err := api.VerifyCredentials()
	if ok && err == nil {
		s.PassTwitterAuth = true
	}
}
