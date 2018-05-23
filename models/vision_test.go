package models_test

import (
	"testing"

	. "github.com/iwataka/mybot/models"
	"github.com/stretchr/testify/assert"
)

func TestVisionFeatures(t *testing.T) {
	cond := NewVisionCondition()
	cond.Label = []string{"label"}
	cond.Face.BlurredLikelihood = "VERY_LIKELY"
	cond.Text = []string{"text"}
	cond.Landmark = []string{"landmark"}
	cond.Logo = []string{"logo"}

	fs := cond.VisionFeatures()
	assert.Len(t, fs, 5)
}

func TestVisionConditionIsEmpty(t *testing.T) {
	cond := NewVisionCondition()
	cond.Label = []string{}
	cond.Face.BlurredLikelihood = ""
	cond.Text = []string{}
	assert.True(t, cond.IsEmpty())
}
