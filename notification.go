package main

type Notification struct {
	Place *PlaceNotification
}

type PlaceNotification struct {
	AllowSelf bool `toml:"allowSelf"`
	Users     []string
}
