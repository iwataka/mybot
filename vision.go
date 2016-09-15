package main

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/vision/v1"
)

// VisionAPI is a wrapper of vision.Service.
type VisionAPI struct {
	api   *vision.Service
	cache *MybotCache
}

// NewVisionAPI takes a path of a user's google-cloud credential file and cache
// and returns a VisionAPI instance for that user.
func NewVisionAPI(cache *MybotCache) (*VisionAPI, error) {
	c, err := google.DefaultClient(context.Background(), vision.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	a, err := vision.New(c)
	if err != nil {
		return nil, err
	}
	return &VisionAPI{a, cache}, nil
}

// VisionCondition is a condition to check whether images match or not by using
// Google Vision API.
type VisionCondition struct {
	Label    []string          `toml:"label"`
	Face     map[string]string `toml:"face"`
	Text     []string          `toml:"text"`
	Landmark []string          `toml:"landmark"`
	Logo     []string          `toml:"logo"`
}

// MatchImages takes image URLs and a Vision condition and returns whether the
// specified images match or not.
func (a *VisionAPI) MatchImages(urls []string, cond *VisionCondition) (bool, error) {
	// No image never match any conditions
	if len(urls) == 0 {
		return false, nil
	}

	imgData := make([][]byte, len(urls))
	for i, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return false, err
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}
		imgData[i] = data
		err = resp.Body.Close()
		if err != nil {
			return false, err
		}
	}

	features := getFeatures(cond)
	if len(features) == 0 {
		return true, nil
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
		return false, err
	}

	for i, r := range res.Responses {
		result, err := r.MarshalJSON()
		if err != nil {
			return false, err
		}
		cache.ImageURL = urls[i]
		cache.ImageAnalysisResult = string(result)
		cache.ImageAnalysisDate = time.Now().String()

		match := true
		if match && r.LabelAnnotations != nil && len(r.LabelAnnotations) != 0 {
			m, err := matchEntity(r.LabelAnnotations, cond.Label)
			if err != nil {
				return false, err
			}
			match = match && m
		}
		if match && r.FaceAnnotations != nil && len(r.FaceAnnotations) != 0 {
			m, err := matchFace(r.FaceAnnotations, cond.Face)
			if err != nil {
				return false, err
			}
			match = match && m
		}
		if match && r.TextAnnotations != nil && len(r.TextAnnotations) != 0 {
			m, err := matchEntity(r.TextAnnotations, cond.Text)
			if err != nil {
				return false, err
			}
			match = match && m
		}
		if match && r.LandmarkAnnotations != nil && len(r.LandmarkAnnotations) != 0 {
			m, err := matchEntity(r.LandmarkAnnotations, cond.Landmark)
			if err != nil {
				return false, err
			}
			match = match && m
		}
		if match && r.LogoAnnotations != nil && len(r.LogoAnnotations) != 0 {
			m, err := matchEntity(r.LogoAnnotations, cond.Logo)
			if err != nil {
				return false, err
			}
			match = match && m
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

func getFeatures(cond *VisionCondition) []*vision.Feature {
	features := []*vision.Feature{}
	if cond.Label != nil && len(cond.Label) != 0 {
		f := &vision.Feature{
			Type:       "LABEL_DETECTION",
			MaxResults: 10,
		}
		features = append(features, f)
	}
	if cond.Face != nil && len(cond.Face) != 0 {
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

func matchFace(as []*vision.FaceAnnotation, face map[string]string) (bool, error) {
	for key, val := range face {
		for _, a := range as {
			match := false
			var err error
			if key == "anger" {
				match, err = regexp.MatchString(val, a.AngerLikelihood)
				if err != nil {
					return false, err
				}
			} else if key == "blurred" {
				match, err = regexp.MatchString(val, a.BlurredLikelihood)
				if err != nil {
					return false, err
				}
			} else if key == "headwear" {
				match, err = regexp.MatchString(val, a.HeadwearLikelihood)
				if err != nil {
					return false, err
				}
			} else if key == "joy" {
				match, err = regexp.MatchString(val, a.JoyLikelihood)
				if err != nil {
					return false, err
				}
			}
			if !match {
				return false, nil
			}
		}
	}
	return true, nil
}
