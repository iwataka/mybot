package mybot

import (
	"testing"

	"github.com/iwataka/mybot/models"
)

func TestLanguageConditionIsEmpty(t *testing.T) {
	c := &models.LanguageCondition{}
	if !c.IsEmpty() {
		t.Fatalf("%v expected empty but not", c)
	}
}

func TestLanguageConditionIsNotEmpty(t *testing.T) {
	min := 0.2
	max := 0.5
	c := &models.LanguageCondition{}
	c.MinSentiment = &min
	c.MaxSentiment = &max
	if c.IsEmpty() {
		t.Fatalf("%v expected not empty but empty", c)
	}
}

func TestLanguageAPIEnabled(t *testing.T) {
	a := &LanguageAPI{}
	if a.Enabled() {
		t.Fatalf("%v expected to be enabled, but not", a)
	}
}
