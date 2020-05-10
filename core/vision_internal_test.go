package core

import (
	"testing"

	"github.com/iwataka/mybot/models"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/vision/v1"
)

func Test_matchEntity(t *testing.T) {
	as := []*vision.EntityAnnotation{
		{
			Description: "foo",
		},
		{
			Description: "bar",
		},
	}
	ds := []string{
		"foo",
		"bar|any",
	}

	flag, err := matchEntity(as, ds)
	require.NoError(t, err)
	require.True(t, flag)
}

func Test_matchFace(t *testing.T) {
	as := []*vision.FaceAnnotation{
		{
			AngerLikelihood:   "LIKELY",
			BlurredLikelihood: "VERY_UNLIKELY",
		},
	}
	face := models.VisionFaceCondition{}
	face.AngerLikelihood = "LIKELY|VERY_LIKELY"

	flag, err := matchFace(as, face)
	require.NoError(t, err)
	require.True(t, flag)
}

func TestVisionAPI_Enabled(t *testing.T) {
	a := &VisionAPI{}
	require.False(t, a.Enabled())
}

func TestVisionAPI_retrieveAnnotateImageResposes(t *testing.T) {
	a := &VisionAPI{}
	urls := []string{"url"}
	imgCache := models.ImageCacheData{}
	imgCache.URL = "url"
	imgCache.AnalysisResult = "{}"
	imgCaches := []models.ImageCacheData{imgCache}

	reses, err := a.retrieveaAnnotateImageResponses(urls, imgCaches, nil)
	require.NoError(t, err)
	require.Len(t, reses, 1)
}
