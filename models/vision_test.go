package models_test

import (
	"testing"

	. "github.com/iwataka/mybot/models"
	"github.com/stretchr/testify/require"
)

func TestVisionCondition_VisionFeatures(t *testing.T) {
	cond := NewVisionCondition()
	cond.Label = []string{"label"}
	cond.Face.BlurredLikelihood = "VERY_LIKELY"
	cond.Text = []string{"text"}
	cond.Landmark = []string{"landmark"}
	cond.Logo = []string{"logo"}

	fs := cond.VisionFeatures()
	require.Len(t, fs, 5)
}

func TestVisionCondition_IsEmpty(t *testing.T) {
	cond := NewVisionCondition()
	cond.Label = []string{}
	cond.Face.BlurredLikelihood = ""
	cond.Text = []string{}
	require.True(t, cond.IsEmpty())
}
