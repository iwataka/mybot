package models

import (
	language "google.golang.org/api/language/v1"
)

func LanguageFeatures(c *LanguageCondition) *language.Features {
	f := &language.Features{}
	if c.MinSentiment != nil || c.MaxSentiment != nil {
		f.ExtractDocumentSentiment = true
	}
	return f
}
