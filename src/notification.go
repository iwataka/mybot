package mybot

// Notification represents a notification about tweets
type Notification struct {
	Place *PlaceNotification
}

// PlaceNotification represents a place notification about tweets
type PlaceNotification struct {
	AllowSelf bool     `json:"allow_self" toml:"allow_self"`
	Users     []string `json:"users,omitempty" toml:"users,omitempty"`
}
