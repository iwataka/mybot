package core

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/utils"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/vision/v1"
)

// VisionAPI is a wrapper of vision.Service.
type VisionAPI struct {
	api *vision.Service
}

// NewVisionMatcher takes a path of a user's google-cloud credential file
// and returns a VisionAPI instance for that user.
func NewVisionMatcher(file string) (VisionMatcher, error) {
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" && len(file) != 0 {
		err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", file)
		if err != nil {
			return nil, utils.WithStack(err)
		}
	}
	c, err := google.DefaultClient(context.Background(), vision.CloudPlatformScope)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	a, err := vision.New(c)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	return &VisionAPI{a}, nil
}

type VisionMatcher interface {
	MatchImages([]string, models.VisionCondition, []models.ImageCacheData) ([]string, []bool, error)
	Enabled() bool
}

// MatchImages takes image URLs and a Vision condition and returns whether the
// specified images match or not.
func (a *VisionAPI) MatchImages(
	urls []string,
	cond models.VisionCondition,
	imgCaches []models.ImageCacheData,
) ([]string, []bool, error) {
	// No image never match any conditions
	if len(urls) == 0 {
		return []string{}, []bool{}, nil
	}

	features := cond.VisionFeatures()
	if len(features) == 0 {
		results := make([]string, len(urls))
		matches := []bool{}
		for i := range matches {
			matches[i] = true
		}
		return results, matches, nil
	}

	responses, err := a.retrieveaAnnotateImageResponses(urls, imgCaches, features)
	if err != nil {
		return nil, nil, utils.WithStack(err)
	}

	results := []string{}
	matches := []bool{}
	for _, r := range responses {
		result, err := r.MarshalJSON()
		if err != nil {
			return nil, nil, utils.WithStack(err)
		}
		results = append(results, string(result))

		match := true
		if match && r.LabelAnnotations != nil && len(r.LabelAnnotations) != 0 {
			m, err := matchEntity(r.LabelAnnotations, cond.Label)
			if err != nil {
				return nil, nil, utils.WithStack(err)
			}
			match = match && m
		}
		if match && r.FaceAnnotations != nil && len(r.FaceAnnotations) != 0 {
			m, err := matchFace(r.FaceAnnotations, cond.Face)
			if err != nil {
				return nil, nil, utils.WithStack(err)
			}
			match = match && m
		}
		if match && r.TextAnnotations != nil && len(r.TextAnnotations) != 0 {
			m, err := matchEntity(r.TextAnnotations, cond.Text)
			if err != nil {
				return nil, nil, utils.WithStack(err)
			}
			match = match && m
		}
		if match && r.LandmarkAnnotations != nil && len(r.LandmarkAnnotations) != 0 {
			m, err := matchEntity(r.LandmarkAnnotations, cond.Landmark)
			if err != nil {
				return nil, nil, utils.WithStack(err)
			}
			match = match && m
		}
		if match && r.LogoAnnotations != nil && len(r.LogoAnnotations) != 0 {
			m, err := matchEntity(r.LogoAnnotations, cond.Logo)
			if err != nil {
				return nil, nil, utils.WithStack(err)
			}
			match = match && m
		}
		matches = append(matches, match)
	}
	return results, matches, nil
}

func (a *VisionAPI) retrieveaAnnotateImageResponses(urls []string, caches []models.ImageCacheData, features []*vision.Feature) ([]*vision.AnnotateImageResponse, error) {
	uncachedUrls := []string{}
	reses := []*vision.AnnotateImageResponse{}
	url2res := map[string]*vision.AnnotateImageResponse{}

	if caches == nil {
		uncachedUrls = urls
	} else {
		for _, url := range urls {
			var exists bool
			for _, cache := range caches {
				if cache.URL == url {
					res := &vision.AnnotateImageResponse{}
					err := json.Unmarshal([]byte(cache.AnalysisResult), res)
					if err != nil {
						return nil, utils.WithStack(err)
					}
					exists = true
					url2res[url] = res
					break
				}
			}
			if !exists {
				uncachedUrls = append(uncachedUrls, url)
			}
		}
	}

	uncachedReses, err := a.retrieveaAnnotateImageResponsesThroughAPI(uncachedUrls, features)
	if err != nil {
		return nil, utils.WithStack(err)
	}
	for i, url := range uncachedUrls {
		url2res[url] = uncachedReses[i]
	}

	for _, url := range urls {
		if res, exists := url2res[url]; exists {
			reses = append(reses, res)
		} else {
			return nil, errors.New("unexpected result")
		}
	}

	return reses, nil
}

func (a *VisionAPI) retrieveaAnnotateImageResponsesThroughAPI(urls []string, features []*vision.Feature) ([]*vision.AnnotateImageResponse, error) {
	if len(urls) == 0 {
		return nil, nil
	}

	imgData := make([][]byte, len(urls))
	for i, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			return nil, utils.WithStack(err)
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, utils.WithStack(err)
		}
		imgData[i] = data
		err = resp.Body.Close()
		if err != nil {
			return nil, utils.WithStack(err)
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
		return nil, utils.WithStack(err)
	}

	return res.Responses, nil
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
				return false, utils.WithStack(err)
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
	if len(face.AngerLikelihood) > 0 {
		var match bool
		for _, a := range as {
			m, err := regexp.MatchString(face.AngerLikelihood, a.AngerLikelihood)
			if err != nil {
				return false, utils.WithStack(err)
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

	if len(face.BlurredLikelihood) > 0 {
		var match bool
		for _, a := range as {
			m, err := regexp.MatchString(face.BlurredLikelihood, a.BlurredLikelihood)
			if err != nil {
				return false, utils.WithStack(err)
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

	if len(face.HeadwearLikelihood) > 0 {
		var match bool
		for _, a := range as {
			m, err := regexp.MatchString(face.HeadwearLikelihood, a.HeadwearLikelihood)
			if err != nil {
				return false, utils.WithStack(err)
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

	if len(face.JoyLikelihood) > 0 {
		var match bool
		for _, a := range as {
			m, err := regexp.MatchString(face.JoyLikelihood, a.JoyLikelihood)
			if err != nil {
				return false, utils.WithStack(err)
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
