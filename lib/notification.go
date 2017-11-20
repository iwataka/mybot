package mybot

// Notification represents a notification about tweets
type Notification struct {
	Place PlaceNotification
}

func NewNotification() Notification {
	return Notification{
		Place: PlaceNotification{},
	}
}

// PlaceNotification represents a place notification about tweets
type PlaceNotification struct {
	AllowSelf bool     `json:"allow_self" toml:"allow_self" bson:"allow_self"`
	Users     []string `json:"users,omitempty" toml:"users,omitempty" bson:"users,omitempty"`
}
