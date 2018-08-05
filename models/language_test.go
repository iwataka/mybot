package models_test

import (
	. "github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/utils"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestLanguageCondition_IsEmpty(t *testing.T) {
	var c LanguageCondition

	c = LanguageCondition{}
	require.True(t, c.IsEmpty())

	c = LanguageCondition{
		MinSentiment: utils.Float64Ptr(0.2),
		MaxSentiment: utils.Float64Ptr(0.8),
	}
	require.False(t, c.IsEmpty())
}

func TestLanguageCondition_LanguageFeatures(t *testing.T) {
	var c LanguageCondition

	c = LanguageCondition{}
	require.False(t, c.LanguageFeatures().ExtractDocumentSentiment)

	c = LanguageCondition{
		MinSentiment: utils.Float64Ptr(0.2),
		MaxSentiment: utils.Float64Ptr(0.8),
	}
	require.True(t, c.LanguageFeatures().ExtractDocumentSentiment)

}
