package mybot

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/iwataka/mybot/models"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/vision/v1"
)

// VisionAPI is a wrapper of vision.Service.
type VisionAPI struct {
	api *vision.Service
}

// NewVisionAPI takes a path of a user's google-cloud credential file
// and returns a VisionAPI instance for that user.
func NewVisionAPI(file string) (*VisionAPI, error) {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" && len(file) != 0 {
		err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", file)
		if err != nil {
			return nil, err
		}
	}
	c, err := google.DefaultClient(context.Background(), vision.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	a, err := vision.New(c)
	if err != nil {
		return nil, err
	}
	return &VisionAPI{a}, nil
}

type VisionMatcher interface {
	MatchImages([]string, models.VisionCondition) ([]string, []bool, error)
	Enabled() bool
}

// MatchImages takes image URLs and a Vision condition and returns whether the
// specified images match or not.
func (a *VisionAPI) MatchImages(
	urls []string,
	cond models.VisionCondition,
) ([]string, []bool, error) {
	// No image never match any conditions
	if len(urls) == 0 {
		return []string{}, []bool{}, nil
	}

	features := cond.VisionFeatures()
	if len(features) == 0 {
		results := make([]string, len(urls), len(urls))
		matches := []bool{}
		for i, _ := range matches {
			matches[i] = true
		}
		return results, matches, nil
	}

	imgData := make([][]byte, len(urls))
	for i, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return nil, nil, err
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, err
		}
		imgData[i] = data
		err = resp.Body.Close()
		if err != nil {
			return nil, nil, err
		}
	}

	imgs := make([]*vision.Image, len(imgData))
	for i, d := range imgData {
		enc := base64.StdEncoding.EncodeToString(d)
		imgs[i] = &vision.Image{Content: enc}
	}

	reqs := make([]*vision.AnnotateImageRequest, len(imgs))
	for i, img := range imgs {
		reqs[i] = &vision.AnnotateImageRequest{
			Image:    img,
			Features: features,
		}
	}

	batch := &vision.BatchAnnotateImagesRequest{
		Requests: reqs,
	}

	res, err := a.api.Images.Annotate(batch).Do()
	if err != nil {
		return nil, nil, err
	}

	results := []string{}
	matches := []bool{}
	for _, r := range res.Responses {
		result, err := r.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		results = append(results, string(result))

		match := true
		if match && r.LabelAnnotations != nil && len(r.LabelAnnotations) != 0 {
			m, err := matchEntity(r.LabelAnnotations, cond.Label)
			if err != nil {
				return nil, nil, err
			}
			match = match && m
		}
		if match && r.FaceAnnotations != nil && len(r.FaceAnnotations) != 0 {
			m, err := matchFace(r.FaceAnnotations, cond.Face)
			if err != nil {
				return nil, nil, err
			}
			match = match && m
		}
		if match && r.TextAnnotations != nil && len(r.TextAnnotations) != 0 {
			m, err := matchEntity(r.TextAnnotations, cond.Text)
			if err != nil {
				return nil, nil, err
			}
			match = match && m
		}
		if match && r.LandmarkAnnotations != nil && len(r.LandmarkAnnotations) != 0 {
			m, err := matchEntity(r.LandmarkAnnotations, cond.Landmark)
			if err != nil {
				return nil, nil, err
			}
			match = match && m
		}
		if match && r.LogoAnnotations != nil && len(r.LogoAnnotations) != 0 {
			m, err := matchEntity(r.LogoAnnotations, cond.Logo)
			if err != nil {
				return nil, nil, err
			}
			match = match && m
		}
		matches = append(matches, match)
	}
	return results, matches, nil
}

func (a *VisionAPI) Enabled() bool {
	return a.api != nil
}

func matchEntity(as []*vision.EntityAnnotation, ds []string) (bool, error) {
	for _, d := range ds {
		match := false
		for _, a := range as {
			m, err := regexp.MatchString(d, a.Description)
			if err != nil {
				return false, err
			}
			if m {
				match = true
				break
			}
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

func matchFace(as []*vision.FaceAnnotation, face models.VisionFaceCondition) (bool, error) {
	for _, a := range as {
		match := false
		var err error

		match, err = regexp.MatchString(face.AngerLikelihood, a.AngerLikelihood)
		if err != nil {
			return false, err
		}

		match, err = regexp.MatchString(face.BlurredLikelihood, a.BlurredLikelihood)
		if err != nil {
			return false, err
		}

		match, err = regexp.MatchString(face.HeadwearLikelihood, a.HeadwearLikelihood)
		if err != nil {
			return false, err
		}

		match, err = regexp.MatchString(face.JoyLikelihood, a.JoyLikelihood)
		if err != nil {
			return false, err
		}

		if !match {
			return false, nil
		}
	}
	return true, nil
}
