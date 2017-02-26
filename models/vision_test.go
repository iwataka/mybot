package models

import (
	"testing"
)

func TestVisionFeatures(t *testing.T) {
	cond := NewVisionCondition()
	cond.Label = []string{"label"}
	cond.Face.BlurredLikelihood = "VERY_LIKELY"
	cond.Text = []string{"text"}
	fs := cond.VisionFeatures()
	if len(fs) != 3 {
		t.Fatalf("%v has %d elements but should have %d", fs, len(fs), 3)
	}
}

func TestVisionConditionIsEmpty(t *testing.T) {
	cond := NewVisionCondition()
	cond.Label = []string{}
	cond.Face.BlurredLikelihood = ""
	cond.Text = []string{}
	if !cond.IsEmpty() {
		t.Fatalf("%v is expected to be empty but not", cond)
	}
}
