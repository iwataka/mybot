package models

type VisionFaceConditionProperties struct {
	AngerLikelihood    string `json:"anger_likelihood,omitempty" toml:"anger_likelihood,omitempty"`
	BlurredLikelihood  string `json:"blurred_likelihood,omitempty" toml:"blurred_likelihood,omitempty"`
	HeadwearLikelihood string `json:"headwear_likelihood,omitempty" toml:"headwear_likelihood,omitempty"`
	JoyLikelihood      string `json:"joy_likelihood,omitempty" toml:"joy_likelihood,omitempty"`
}
