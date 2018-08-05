package mybot

import (
	"testing"

	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/utils"
	"github.com/stretchr/testify/require"
)

func TestLanguageCondition_IsEmpty(t *testing.T) {
	c := &models.LanguageCondition{}
	require.True(t, c.IsEmpty())
}

func TestLanguageCondition_IsEmpty_ReturnsFalse(t *testing.T) {
	c := &models.LanguageCondition{}
	c.MinSentiment = utils.Float64Ptr(0.2)
	c.MaxSentiment = utils.Float64Ptr(0.5)
	require.False(t, c.IsEmpty())
}

func TestLanguageAPI_Enabled(t *testing.T) {
	a := &LanguageAPI{}
	require.False(t, a.Enabled())
}
