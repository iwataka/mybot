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

	MonitorConfigStatus      bool
	MonitorTwitterCred       bool
	MonitorGCloudCred        bool
	monitorConfigStatusChans []chan bool
	monitorTwitterCredChans  []chan bool
	monitorGCloudCredChans   []chan bool
	MonitorConfigStatusMutex *sync.Mutex
	MonitorTwitterCredMutex  *sync.Mutex
	MonitorGCloudCredMutex   *sync.Mutex

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

		false,
		false,
		false,
		[]chan bool{},
		[]chan bool{},
		[]chan bool{},
		new(sync.Mutex),
		new(sync.Mutex),
		new(sync.Mutex),

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

func (s *Status) AddMonitorConfigStatusChan(c chan bool) {
	s.MonitorConfigStatusMutex.Lock()
	s.monitorConfigStatusChans = append(s.monitorConfigStatusChans, c)
	s.MonitorConfigStatusMutex.Unlock()
}

func (s *Status) AddMonitorTwitterCredChan(c chan bool) {
	s.MonitorTwitterCredMutex.Lock()
	s.monitorTwitterCredChans = append(s.monitorTwitterCredChans, c)
	s.MonitorTwitterCredMutex.Unlock()
}

func (s *Status) AddMonitorGCloudCredChan(c chan bool) {
	s.MonitorGCloudCredMutex.Lock()
	s.monitorGCloudCredChans = append(s.monitorGCloudCredChans, c)
	s.MonitorGCloudCredMutex.Unlock()
}

func (s *Status) SendToMonitorConfigStatusChans(flag bool) {
	for _, c := range s.monitorConfigStatusChans {
		c <- flag
	}
	s.monitorConfigStatusChans = []chan bool{}
}

func (s *Status) SendToMonitorTwitterCredChans(flag bool) {
	for _, c := range s.monitorTwitterCredChans {
		c <- flag
	}
	s.monitorTwitterCredChans = []chan bool{}
}

func (s *Status) SendToMonitorGCloudCredChans(flag bool) {
	for _, c := range s.monitorGCloudCredChans {
		c <- flag
	}
	s.monitorGCloudCredChans = []chan bool{}
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
