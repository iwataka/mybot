package main

type Notification struct {
	Place *PlaceNotification
}

type PlaceNotification struct {
	AllowSelf bool `toml:"allow_self"`
	Users     []string
}
