package models

import (
	"github.com/jinzhu/gorm"
)

// Cache is a cache data stored in the specified file.
type Cache struct {
	gorm.Model
	Tweets       []TweetActionCache
	TwitterUsers []TwitterUserCache
	LatestDMID   int64 `gorm:"latest_dm_id"`
	Visions      []VisionCache
}

type VisionCache struct {
	gorm.Model
	CacheID        uint
	URL            string
	Src            string
	AnalysisResult string
	AnalysisDate   string
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

type TwitterAction struct {
	gorm.Model
	Retweet     bool
	Favorite    bool
	Follow      bool
	Collections []TwitterCollection
}

func (a *TwitterAction) GetCollections() []string {
	result := []string{}
	for _, c := range a.Collections {
		result = append(result, c.Name)
	}
	return result
}

func (a *TwitterAction) SetCollections(cols []string) {
	a.Collections = []TwitterCollection{}
	for _, col := range cols {
		c := TwitterCollection{
			TwitterActionID: a.ID,
			Name:            col,
		}
		a.Collections = append(a.Collections, c)
	}
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

type SlackChannel struct {
	gorm.Model
	SlackActionID uint
	Name          string
}
