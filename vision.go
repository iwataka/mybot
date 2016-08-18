package main

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/antonholmquist/jason"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/vision/v1"
)

type VisionAPI struct {
	*vision.Service
	ProjectID string
	cache     *MybotCache
}

func NewVisionAPI(path string, cache *MybotCache) (*VisionAPI, error) {
	cred, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg, err := google.JWTConfigFromJSON(cred, vision.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	v, err := jason.NewObjectFromBytes(cred)
	projectID, err := v.GetString("project_id")
	if err != nil {
		return nil, err
	}
	c := cfg.Client(context.Background())
	a, err := vision.New(c)
	if err != nil {
		return nil, err
	}
	return &VisionAPI{a, projectID, cache}, nil
}

type VisionCondition struct {
	Label    []string
	Face     map[string]string
	Text     []string
	Landmark []string
	Logo     []string
}

func (a *VisionAPI) MatchImage(urls []string, cond *VisionCondition) (bool, error) {
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
		resp.Body.Close()
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

	res, err := a.Images.Annotate(batch).Do()
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
