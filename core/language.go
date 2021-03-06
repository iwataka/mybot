package core

import (
	"os"

	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/utils"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/language/v1"
)

type LanguageAPI struct {
	api *language.Service
}

func NewLanguageMatcher(file string) (LanguageMatcher, error) {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" && len(file) != 0 {
		err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", file)
		if err != nil {
			return nil, utils.WithStack(err)
		}
	}
	c, err := google.DefaultClient(context.Background(), language.CloudPlatformScope)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	a, err := language.New(c)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	return &LanguageAPI{a}, nil
}

type LanguageMatcher interface {
	MatchText(string, models.LanguageCondition) (string, bool, error)
	Enabled() bool
}

func (a *LanguageAPI) MatchText(
	text string,
	cond models.LanguageCondition,
) (string, bool, error) {
	if cond.IsEmpty() {
		return "", true, nil
	}

	f := cond.LanguageFeatures()
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
		return "", false, utils.WithStack(err)
	}

	bytes, err := res.MarshalJSON()
	if err != nil {
		return "", false, utils.WithStack(err)
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
