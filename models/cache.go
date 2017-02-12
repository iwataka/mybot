package models

import (
	"github.com/jinzhu/gorm"
)

// Cache is a cache data stored in the specified file.
type Cache struct {
	gorm.Model
	Tweets       []TweetActionCache
	TwitterUsers []TwitterUserCache
	LatestDMID   int64
	Visions      []VisionCache
}

type VisionCache struct {
	gorm.Model
	VisionCacheProperties
	CacheID uint
}

type VisionCacheProperties struct {
	URL            string `json:"url" toml:"url"`
	Src            string `json:"src" toml:"src"`
	AnalysisResult string `json:"analysis_result" toml:"analysis_result"`
	AnalysisDate   string `json:"analysis_date" toml:"analysis_date"`
}

type TwitterUserCache struct {
	gorm.Model
	CacheID          uint
	ScreenName       string
	LatestTweetID    int64
	LatestFavoriteID int64
}

type TweetActionCache struct {
	gorm.Model
	CacheID         uint
	TweetID         int64
	TwitterActionID uint
	TwitterAction   TwitterAction
	SlackActionID   uint
	SlackAction     SlackAction
}

type SlackAction struct {
	gorm.Model
	Channels []SlackChannel
}

func (a *SlackAction) GetChannels() []string {
	result := []string{}
	for _, c := range a.Channels {
		result = append(result, c.Name)
	}
	return result
}

func (a *SlackAction) SetChannels(chs []string) {
	a.Channels = []SlackChannel{}
	for _, ch := range chs {
		c := SlackChannel{
			SlackActionID: a.ID,
			Name:          ch,
		}
		a.Channels = append(a.Channels, c)
	}
}

type TwitterCollection struct {
	gorm.Model
	TwitterActionID uint
	Name            string
}
