package mybot

import (
	"context"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/language/v1"
)

type LanguageAPI struct {
	api    *language.Service
	cache  *MybotCache
	config *MybotConfig
	File   string
}

func NewLanguageAPI(cache *MybotCache, config *MybotConfig, file string) (*LanguageAPI, error) {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" && len(file) != 0 {
		err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", file)
		if err != nil {
			return nil, err
		}
	}
	c, err := google.DefaultClient(context.Background(), language.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	a, err := language.New(c)
	if err != nil {
		return nil, err
	}
	return &LanguageAPI{a, cache, config, file}, nil
}

type LanguageCondition struct {
	MinSentiment *float64 `toml:"min_sentiment,omitempty"`
	MaxSentiment *float64 `toml:"max_sentiment,omitempty"`
}

func (a *LanguageAPI) MatchText(
	text string,
	cond *LanguageCondition,
) (string, bool, error) {
	f := LanguageFeatures(cond)
	if !f.ExtractDocumentSentiment && !f.ExtractEntities && !f.ExtractSyntax {
		return "", false, nil
	}

	doc := &language.Document{
		Content: text,
		Type:    "PLAIN_TEXT",
	}
	req := &language.AnnotateTextRequest{
		Document: doc,
		Features: f,
	}

	res, err := a.api.Documents.AnnotateText(req).Do(nil)
	if err != nil {
		return "", false, err
	}

	bytes, err := res.MarshalJSON()
	if err != nil {
		return "", false, err
	}

	if f.ExtractDocumentSentiment {
		score := res.DocumentSentiment.Score
		if cond.MinSentiment != nil && score < *cond.MinSentiment {
			return string(bytes), false, nil
		}
		if cond.MaxSentiment != nil && *cond.MaxSentiment < score {
			return string(bytes), false, nil
		}
	}

	return string(bytes), true, nil
}

func LanguageFeatures(c *LanguageCondition) *language.Features {
	f := &language.Features{}
	if c.MinSentiment != nil || c.MaxSentiment != nil {
		f.ExtractDocumentSentiment = true
	}
	return f
}