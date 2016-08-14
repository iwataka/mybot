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

func (a *VisionAPI) MatchImageDescription(urls []string, ds []string) (bool, error) {
	// No image never match any description
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

	imgs := make([]*vision.Image, len(imgData))
	for i, d := range imgData {
		enc := base64.StdEncoding.EncodeToString(d)
		imgs[i] = &vision.Image{Content: enc}
	}

	feature := &vision.Feature{
		Type:       "LABEL_DETECTION",
		MaxResults: 10,
	}

	reqs := make([]*vision.AnnotateImageRequest, len(imgs))
	for i, img := range imgs {
		reqs[i] = &vision.AnnotateImageRequest{
			Image:    img,
			Features: []*vision.Feature{feature},
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

		match, err := a.matchDescription(r.LabelAnnotations, ds)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

func (a *VisionAPI) matchDescription(as []*vision.EntityAnnotation, ds []string) (bool, error) {
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
