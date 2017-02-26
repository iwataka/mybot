package models

import (
	language "google.golang.org/api/language/v1"
)

type LanguageCondition struct {
	MinSentiment *float64 `json:"min_sentiment,omitempty" toml:"min_sentiment,omitempty"`
	MaxSentiment *float64 `json:"max_sentiment,omitempty" toml:"max_sentiment,omitempty"`
}

func (c *LanguageCondition) IsEmpty() bool {
	return c.MinSentiment == nil && c.MaxSentiment == nil
}

func (c *LanguageCondition) LanguageFeatures() *language.Features {
	f := &language.Features{}
	if c.MinSentiment != nil || c.MaxSentiment != nil {
		f.ExtractDocumentSentiment = true
	}
	return f
}
