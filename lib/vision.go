package mybot

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/vision/v1"
)

// VisionAPI is a wrapper of vision.Service.
type VisionAPI struct {
	api    *vision.Service
	cache  *Cache
	config *Config
	File   string
}

// NewVisionAPI takes a path of a user's google-cloud credential file and cache
// and returns a VisionAPI instance for that user.
func NewVisionAPI(cache *Cache, config *Config, file string) (*VisionAPI, error) {
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
	return &VisionAPI{a, cache, config, file}, nil
}

// VisionCondition is a condition to check whether images match or not by using
// Google Vision API.
type VisionCondition struct {
	Label    []string             `toml:"label,omitempty"`
	Face     *VisionFaceCondition `toml:"face,omitempty"`
	Text     []string             `toml:"text,omitempty"`
	Landmark []string             `toml:"landmark,omitempty"`
	Logo     []string             `toml:"logo,omitempty"`
}

func (c *VisionCondition) isEmpty() bool {
	return (c.Label == nil || len(c.Label) == 0) &&
		(c.Face == nil || c.Face.isEmpty()) &&
		(c.Text == nil || len(c.Text) == 0) &&
		(c.Landmark == nil || len(c.Landmark) == 0) &&
		(c.Logo == nil || len(c.Logo) == 0)
}

type VisionFaceCondition struct {
	AngerLikelihood    string `toml:"anger_likelihood,omitempty"`
	BlurredLikelihood  string `toml:"blurred_likelihood,omitempty"`
	HeadwearLikelihood string `toml:"headwear_likelihood,omitempty"`
	JoyLikelihood      string `toml:"joy_likelihood,omitempty"`
}

func (c *VisionFaceCondition) isEmpty() bool {
	return len(c.AngerLikelihood) == 0 &&
		len(c.BlurredLikelihood) == 0 &&
		len(c.HeadwearLikelihood) == 0 &&
		len(c.JoyLikelihood) == 0
}

// MatchImages takes image URLs and a Vision condition and returns whether the
// specified images match or not.
func (a *VisionAPI) MatchImages(
	urls []string,
	cond *VisionCondition,
) ([]string, []bool, error) {
	// No image never match any conditions
	if len(urls) == 0 {
		return []string{}, []bool{}, nil
	}

	features := VisionFeatures(cond)
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

func VisionFeatures(cond *VisionCondition) []*vision.Feature {
	features := []*vision.Feature{}
	if cond.Label != nil && len(cond.Label) != 0 {
		f := &vision.Feature{
			Type:       "LABEL_DETECTION",
			MaxResults: 10,
		}
		features = append(features, f)
	}
	if cond.Face != nil && !cond.Face.isEmpty() {
		f := &vision.Feature{
			Type:       "FACE_DETECTION",
			MaxResults: 10,
		}
		features = append(features, f)
	}
	if cond.Text != nil && len(cond.Text) != 0 {
		f := &vision.Feature{
			Type:       "TEXT_DETECTION",
			MaxResults: 10,
		}
		features = append(features, f)
	}
	if cond.Landmark != nil && len(cond.Landmark) != 0 {
		f := &vision.Feature{
			Type:       "LANDMARK_DETECTION",
			MaxResults: 10,
		}
		features = append(features, f)
	}
	if cond.Logo != nil && len(cond.Logo) != 0 {
		f := &vision.Feature{
			Type:       "LOGO_DETECTION",
			MaxResults: 10,
		}
		features = append(features, f)
	}
	return features
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

func matchFace(as []*vision.FaceAnnotation, face *VisionFaceCondition) (bool, error) {
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
