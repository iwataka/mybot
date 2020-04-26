package models

type ImageCacheData struct {
	VisionCacheProperties
}

type VisionCacheProperties struct {
	URL            string `json:"url" toml:"url" bson:"url" yaml:"url"`
	Src            string `json:"src" toml:"src" bson:"src" yaml:"src"`
	AnalysisResult string `json:"analysis_result" toml:"analysis_result" bson:"analysis_result" yaml:"analysis_result"`
	AnalysisDate   string `json:"analysis_date" toml:"analysis_date" bson:"analysis_date" yaml:"analysis_date"`
}
