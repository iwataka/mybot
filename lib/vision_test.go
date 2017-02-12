package mybot

import (
	"testing"

	"google.golang.org/api/vision/v1"
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
	if !cond.isEmpty() {
		t.Fatalf("%v is expected to be empty but not", cond)
	}
}

func TestMatchEntity(t *testing.T) {
	as := []*vision.EntityAnnotation{
		&vision.EntityAnnotation{
			Description: "foo",
		},
		&vision.EntityAnnotation{
			Description: "bar",
		},
	}
	ds := []string{
		"foo",
		"bar|any",
	}
	flag, err := matchEntity(as, ds)
	if err != nil {
		t.Fatal(err)
	}
	if !flag {
		t.Fatalf("%v expected but %v found", true, flag)
	}
}

func TestMatchFace(t *testing.T) {
	as := []*vision.FaceAnnotation{
		&vision.FaceAnnotation{
			AngerLikelihood:   "LIKELY",
			BlurredLikelihood: "VERY_UNLIKELY",
		},
	}
	face := &VisionFaceCondition{}
	face.AngerLikelihood = "LIKELY|VERY_LIKELY"
	flag, err := matchFace(as, face)
	if err != nil {
		t.Fatal(err)
	}
	if !flag {
		t.Fatalf("%v expeted but %v found", true, flag)
	}
}

func TestVisionAPIEnabled(t *testing.T) {
	a := &VisionAPI{}
	if a.Enabled() {
		t.Fatalf("%v expected to be enabled, but not", a)
	}
}
