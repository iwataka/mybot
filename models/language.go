package models

type LanguageConditionProperties struct {
	MinSentiment *float64 `json:"min_sentiment,omitempty" toml:"min_sentiment,omitempty"`
	MaxSentiment *float64 `json:"max_sentiment,omitempty" toml:"max_sentiment,omitempty"`
}
