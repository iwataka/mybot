package models

type TweetFilterProperties struct {
	HasMedia           *bool  `json:"has_media,omitempty" toml:"has_media,omitempty"`
	HasURL             *bool  `json:"has_url,omitempty" toml:"has_url,omitempty"`
	Retweeted          *bool  `json:"retweeted,omitempty" toml:"retweeted,omitempty"`
	FavoriteThreshold  *int   `json:"favorite_threshold" toml:"favorite_threshold"`
	RetweetedThreshold *int   `json:"retweeted_threshold" toml:"retweeted_threshold"`
	Lang               string `json:"lang,omitempty" toml:"lang,omitempty"`
}
