package mybot

import (
	"testing"

	"google.golang.org/api/vision/v1"
)

func TestVisionFeatures(t *testing.T) {
	cond := &VisionCondition{
		Label: []string{"label"},
		Face: &VisionFaceCondition{
			BlurredLikelihood: "VERY_LIKELY",
		},
		Text: []string{"text"},
	}
	fs := cond.VisionFeatures()
	if len(fs) != 3 {
		t.Fatalf("%v has %d elements but should have %d", fs, len(fs), 3)
	}
}

func TestVisionConditionIsEmpty(t *testing.T) {
	cond := &VisionCondition{
		Label: []string{},
		Face: &VisionFaceCondition{
			BlurredLikelihood: "",
		},
		Text: []string{},
	}
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
	face := &VisionFaceCondition{
		AngerLikelihood: "LIKELY|VERY_LIKELY",
	}
	flag, err := matchFace(as, face)
	if err != nil {
		t.Fatal(err)
	}
	if !flag {
		t.Fatalf("%v expeted but %v found", true, flag)
	}
}
