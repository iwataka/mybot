package main

import (
	"encoding/base64"
	"io/ioutil"
	"regexp"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/vision/v1"
)

type VisionAPI struct {
	*vision.Service
	ProjectID string
}

func NewVisionAPI(path string) (*VisionAPI, error) {
	cred, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg, err := google.JWTConfigFromJSON(cred, vision.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	projectID := strings.Split(cfg.Email, "@")[0]
	c := cfg.Client(context.Background())
	a, err := vision.New(c)
	if err != nil {
		return nil, err
	}
	return &VisionAPI{a, projectID}, nil
}

func (a *VisionAPI) MatchImageDescription(imgData [][]byte, ds []string) (bool, error) {
	// No image never match any description
	if len(imgData) == 0 {
		return false, nil
	}

	imgs := make([]*vision.Image, len(imgData))
	for i, d := range imgData {
		enc := base64.StdEncoding.EncodeToString(d)
		imgs[i] = &vision.Image{Content: enc}
	}

	feature := &vision.Feature{
		Type: "LABEL_DETECTION",
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

	for _, r := range res.Responses {
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
