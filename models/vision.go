package models

import (
	vision "google.golang.org/api/vision/v1"
)

// VisionCondition is a condition to check whether images match or not by using
// Google Vision API.
type VisionCondition struct {
	Label    []string             `json:"label,omitempty" toml:"label,omitempty" bson:"label,omitempty"`
	Face     *VisionFaceCondition `json:"face,omitempty" toml:"face,omitempty" bson:"face,omitempty"`
	Text     []string             `json:"text,omitempty" toml:"text,omitempty" bson:"text,omitempty"`
	Landmark []string             `json:"landmark,omitempty" toml:"landmark,omitempty" bson:"landmark,omitempty"`
	Logo     []string             `json:"logo,omitempty" toml:"logo,omitempty" bson:"logo,omitempty"`
}

func NewVisionCondition() *VisionCondition {
	return &VisionCondition{
		Face: &VisionFaceCondition{},
	}
}

func (c *VisionCondition) IsEmpty() bool {
	return (c.Label == nil || len(c.Label) == 0) &&
		(c.Face == nil || c.Face.IsEmpty()) &&
		(c.Text == nil || len(c.Text) == 0) &&
		(c.Landmark == nil || len(c.Landmark) == 0) &&
		(c.Logo == nil || len(c.Logo) == 0)
}

func (cond *VisionCondition) VisionFeatures() []*vision.Feature {
	features := []*vision.Feature{}
	if cond.Label != nil && len(cond.Label) != 0 {
		f := &vision.Feature{
			Type:       "LABEL_DETECTION",
			MaxResults: 10,
		}
		features = append(features, f)
	}
	if cond.Face != nil && !cond.Face.IsEmpty() {
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

type VisionFaceCondition struct {
	AngerLikelihood    string `json:"anger_likelihood,omitempty" toml:"anger_likelihood,omitempty" bson:"anger_likelihood,omitempty"`
	BlurredLikelihood  string `json:"blurred_likelihood,omitempty" toml:"blurred_likelihood,omitempty" bson:"blurred_likelihood,omitempty"`
	HeadwearLikelihood string `json:"headwear_likelihood,omitempty" toml:"headwear_likelihood,omitempty" bson:"headwear_likelihood,omitempty"`
	JoyLikelihood      string `json:"joy_likelihood,omitempty" toml:"joy_likelihood,omitempty" bson:"joy_likelihood,omitempty"`
}

func (c *VisionFaceCondition) IsEmpty() bool {
	return len(c.AngerLikelihood) == 0 &&
		len(c.BlurredLikelihood) == 0 &&
		len(c.HeadwearLikelihood) == 0 &&
		len(c.JoyLikelihood) == 0
}
