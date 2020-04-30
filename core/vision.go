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
		match, err := matchImage(r, cond)
		if err != nil {
			return nil, nil, utils.WithStack(err)
		}
		matches = append(matches, match)
	}
	return results, matches, nil
}

func matchImage(r *vision.AnnotateImageResponse, cond models.VisionCondition) (bool, error) {
	matchLabel, err := matchEntity(r.LabelAnnotations, cond.Label)
	if err != nil {
		return false, utils.WithStack(err)
	}
	if !matchLabel {
		return false, nil
	}

	matchFace, err := matchFace(r.FaceAnnotations, cond.Face)
	if err != nil {
		return false, utils.WithStack(err)
	}
	if !matchFace {
		return false, nil
	}

	matchText, err := matchEntity(r.TextAnnotations, cond.Text)
	if err != nil {
		return false, utils.WithStack(err)
	}
	if !matchText {
		return false, nil
	}

	matchLandmark, err := matchEntity(r.LandmarkAnnotations, cond.Landmark)
	if err != nil {
		return false, utils.WithStack(err)
	}
	if !matchLandmark {
		return false, nil
	}

	matchLogo, err := matchEntity(r.LogoAnnotations, cond.Logo)
	if err != nil {
		return false, utils.WithStack(err)
	}
	if !matchLogo {
		return false, nil
	}

	return true, nil
}

func (a *VisionAPI) retrieveaAnnotateImageResponses(urls []string, caches []models.ImageCacheData, features []*vision.Feature) ([]*vision.AnnotateImageResponse, error) {
	uncachedUrls := []string{}
	reses := []*vision.AnnotateImageResponse{}
	url2res := map[string]*vision.AnnotateImageResponse{}

	if caches == nil {
		uncachedUrls = urls
	} else {
		for _, url := range urls {
			res, err := getAnnotateImageResponseFromCache(url, caches)
			if err != nil {
				return nil, err
			}
			if res != nil {
				url2res[url] = res
			} else {
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

func getAnnotateImageResponseFromCache(url string, caches []models.ImageCacheData) (*vision.AnnotateImageResponse, error) {
	for _, cache := range caches {
		if cache.URL == url {
			res := &vision.AnnotateImageResponse{}
			err := json.Unmarshal([]byte(cache.AnalysisResult), res)
			if err != nil {
				return nil, utils.WithStack(err)
			}
			return res, nil
		}
	}
	return nil, nil
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
	if len(as) == 0 {
		return true, nil
	}

	for _, d := range ds {
		m, err := matchEntityAnnotationBySinglePattern(as, d)
		if err != nil {
			return false, utils.WithStack(err)
		}
		if !m {
			return false, nil
		}
	}
	return true, nil
}

func matchEntityAnnotationBySinglePattern(as []*vision.EntityAnnotation, d string) (bool, error) {
	for _, a := range as {
		m, err := regexp.MatchString(d, a.Description)
		if err != nil {
			return false, utils.WithStack(err)
		}
		if m {
			return true, nil
		}
	}
	return false, nil
}

func matchFace(as []*vision.FaceAnnotation, face models.VisionFaceCondition) (bool, error) {
	if len(as) == 0 {
		return true, nil
	}

	getAnger := func(a *vision.FaceAnnotation) string { return a.AngerLikelihood }
	matchAnger, err := matchEachFaceAnnotation(as, face.AngerLikelihood, getAnger)
	if err != nil {
		return false, err
	}
	if !matchAnger {
		return false, nil
	}

	getBlurred := func(a *vision.FaceAnnotation) string { return a.BlurredLikelihood }
	matchBlurred, err := matchEachFaceAnnotation(as, face.BlurredLikelihood, getBlurred)
	if err != nil {
		return false, err
	}
	if !matchBlurred {
		return false, nil
	}

	getHeadwear := func(a *vision.FaceAnnotation) string { return a.HeadwearLikelihood }
	matchHeadwear, err := matchEachFaceAnnotation(as, face.HeadwearLikelihood, getHeadwear)
	if err != nil {
		return false, err
	}
	if !matchHeadwear {
		return false, nil
	}

	getJoy := func(a *vision.FaceAnnotation) string { return a.JoyLikelihood }
	matchJoy, err := matchEachFaceAnnotation(as, face.JoyLikelihood, getJoy)
	if err != nil {
		return false, err
	}
	if !matchJoy {
		return false, nil
	}

	return true, nil
}

func matchEachFaceAnnotation(as []*vision.FaceAnnotation, pattern string, toText func(*vision.FaceAnnotation) string) (bool, error) {
	for _, a := range as {
		m, err := regexp.MatchString(pattern, toText(a))
		if err != nil {
			return false, utils.WithStack(err)
		}
		if m {
			return true, nil
		}
	}
	return false, nil
}
