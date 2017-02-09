package mybot

import (
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/language/v1"
)

type LanguageAPI struct {
	api  *language.Service
	File string
}

func NewLanguageAPI(file string) (*LanguageAPI, error) {
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
	return &LanguageAPI{a, file}, nil
}

type LanguageCondition struct {
	MinSentiment *float64 `json:"min_sentiment,omitempty" toml:"min_sentiment,omitempty"`
	MaxSentiment *float64 `json:"max_sentiment,omitempty" toml:"max_sentiment,omitempty"`
}

func (c *LanguageCondition) isEmpty() bool {
	return c.MinSentiment == nil && c.MaxSentiment == nil
}

type LanguageMatcher interface {
	MatchText(string, *LanguageCondition) (string, bool, error)
	Enabled() bool
}

func (a *LanguageAPI) MatchText(
	text string,
	cond *LanguageCondition,
) (string, bool, error) {
	f := LanguageFeatures(cond)
	// This means that nothing to do with language API.
	if !f.ExtractDocumentSentiment && !f.ExtractEntities && !f.ExtractSyntax {
		return "", true, nil
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

func (a *LanguageAPI) Enabled() bool {
	return a.api != nil
}

func LanguageFeatures(c *LanguageCondition) *language.Features {
	f := &language.Features{}
	if c.MinSentiment != nil || c.MaxSentiment != nil {
		f.ExtractDocumentSentiment = true
	}
	return f
}
