package main

type Notification struct {
	Place *PlaceNotification
}

type PlaceNotification struct {
	AllowSelf bool `yaml:"allowSelf"`
	Users     []string
}
