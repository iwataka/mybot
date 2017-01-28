package mybot

import (
	"testing"
)

func TestLanguageConditionIsEmpty(t *testing.T) {
	c := &LanguageCondition{}
	if !c.isEmpty() {
		t.Fatalf("%v expected empty but not", c)
	}
}

func TestLanguageConditionIsNotEmpty(t *testing.T) {
	min := 0.2
	max := 0.5
	c := &LanguageCondition{&min, &max}
	if c.isEmpty() {
		t.Fatalf("%v expected not empty but empty", c)
	}
}
