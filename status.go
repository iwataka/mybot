package main

// MybotStatus represents a status of this app
type MybotStatus struct {
	TwitterListenMyselfStatus bool
	TwitterListenUsersStatus  bool
	GithubStatus              bool
	TwitterStatus             bool
	MonitorConfigStatus       bool
	MonitorTwitterCred        bool
	MonitorGCloudCred         bool
	HttpStatus                bool
}
