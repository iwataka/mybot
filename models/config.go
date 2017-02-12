package models

type TimelineProperties struct {
	ExcludeReplies *bool `json:"exclude_replies" toml:"exclude_replies"`
	IncludeRts     *bool `json:"include_rts" toml:"include_rts"`
}

type FavoriteProperties struct {
}

type SearchProperties struct {
	Queries    []string `json:"queries" toml:"queries"`
	ResultType string   `json:"result_type,omitempty" toml:"result_type,omitempty"`
}

type AccountProperties struct {
	ScreenNames []string `json:"screen_names" toml:"screen_names"`
}

type SourceProperties struct {
	Count *int `json:"count,omitempty" toml:"count,omitempty"`
}
