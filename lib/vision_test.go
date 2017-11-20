package mybot

import (
	"testing"

	"github.com/iwataka/mybot/models"
	"google.golang.org/api/vision/v1"
)

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
	face := models.VisionFaceCondition{}
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
