package mybot

// Notification represents a notification about tweets
type Notification struct {
	Place *PlaceNotification
}

// PlaceNotification represents a place notification about tweets
type PlaceNotification struct {
	AllowSelf bool     `toml:"allow_self"`
	Users     []string `toml:"users,omitempty"`
}
