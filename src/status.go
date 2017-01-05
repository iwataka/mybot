package mybot

// MybotStatus represents a status of this app
type MybotStatus struct {
	TwitterListenMyselfStatus bool
	TwitterListenUsersStatus  bool
	TwitterStatus             bool
	MonitorConfigStatus       bool
	MonitorTwitterCred        bool
	MonitorGCloudCred         bool
	ServerStatus              bool
}
