package models_test

import (
	. "github.com/iwataka/mybot/models"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestLanguageCondition_IsEmpty(t *testing.T) {
	var c LanguageCondition

	c = LanguageCondition{}
	assert.True(t, c.IsEmpty())

	min, max := 0.2, 0.8
	c = LanguageCondition{&min, &max}
	assert.False(t, c.IsEmpty())
}

func TestLanguageCondition_LanguageFeatures(t *testing.T) {
	var c LanguageCondition

	c = LanguageCondition{}
	assert.False(t, c.LanguageFeatures().ExtractDocumentSentiment)

	min, max := 0.2, 0.8
	c = LanguageCondition{&min, &max}
	assert.True(t, c.LanguageFeatures().ExtractDocumentSentiment)

}
