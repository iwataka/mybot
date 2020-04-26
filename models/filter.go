package models

type FilterProperties struct {
	HasMedia           *bool  `json:"has_media,omitempty" toml:"has_media,omitempty" bson:"has_media,omitempty" yaml:"has_media,omitempty"`
	FavoriteThreshold  *int   `json:"favorite_threshold" toml:"favorite_threshold" bson:"favorite_threshold" yaml:"favorite_threshold"`
	RetweetedThreshold *int   `json:"retweeted_threshold" toml:"retweeted_threshold" bson:"retweeted_threshold" yaml:"retweeted_threshold"`
	Lang               string `json:"lang,omitempty" toml:"lang,omitempty" bson:"lang,omitempty" yaml:"lang,omitempty"`
}
