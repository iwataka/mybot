package models_test

import (
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/utils"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestLanguageCondition_IsEmpty(t *testing.T) {
	var c models.LanguageCondition

	c = models.LanguageCondition{}
	require.True(t, c.IsEmpty())

	c = models.LanguageCondition{
		MinSentiment: utils.Float64Ptr(0.2),
		MaxSentiment: utils.Float64Ptr(0.8),
	}
	require.False(t, c.IsEmpty())
}

func TestLanguageCondition_LanguageFeatures(t *testing.T) {
	var c models.LanguageCondition

	c = models.LanguageCondition{}
	require.False(t, c.LanguageFeatures().ExtractDocumentSentiment)

	c = models.LanguageCondition{
		MinSentiment: utils.Float64Ptr(0.2),
		MaxSentiment: utils.Float64Ptr(0.8),
	}
	require.True(t, c.LanguageFeatures().ExtractDocumentSentiment)

}
