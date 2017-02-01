package mybot

import (
	"os"
	"path/filepath"
)

type Cache interface {
	GetLatestTweetID(screenName string) (int64, bool)
	SetLatestTweetID(screenName string, id int64) error
	GetLatestFavoriteID(screenName string) (int64, bool)
	SetLatestFavoriteID(screenName string, id int64) error
	GetLatestDMID() int64
	SetLatestDMID(id int64) error
	GetTweetAction(tweetID string) (*TwitterAction, bool)
	SetTweetAction(tweetID string, action *TwitterAction) error
	GetLatestImages(num int) []ImageCacheData
	SetImage(data ImageCacheData) error
	Save() error
}

// FileCache is a cache data of mybot
// TODO: This should be stored in DB, not .json file.
type FileCache struct {
	LatestTweetID    map[string]int64          `json:"latest_tweet_id" toml:"latest_tweet_id"`
	LatestFavoriteID map[string]int64          `json:"latest_favorite_id" toml:"lates_favorite_id"`
	LatestDMID       int64                     `json:"latest_dm_id" toml:"latest_dm_id"`
	Tweet2Action     map[string]*TwitterAction `json:"tweet_to_action" toml:"tweet_to_action"`
	Images           []ImageCacheData          `json:"images" toml:"images"`
	File             string                    `json:"-" toml:"-"`
}

type ImageCacheData struct {
	URL            string `json:"url" toml:"url"`
	Src            string `json:"src" toml:"src"`
	AnalysisResult string `json:"analysis_result" toml:"analysis_result"`
	AnalysisDate   string `json:"analysis_date" toml:"analysis_date"`
}

// NewFileCache creates an instance of Cache
func NewFileCache(path string) (*FileCache, error) {
	c := &FileCache{
		make(map[string]int64),
		make(map[string]int64),
		0,
		make(map[string]*TwitterAction),
		[]ImageCacheData{},
		path,
	}

	info, _ := os.Stat(path)
	if info != nil && !info.IsDir() {
		err := DecodeFile(path, c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *FileCache) GetLatestTweetID(screenName string) (int64, bool) {
	id, exists := c.LatestTweetID[screenName]
	return id, exists
}

func (c *FileCache) SetLatestTweetID(screenName string, id int64) error {
	c.LatestTweetID[screenName] = id
	return nil
}

func (c *FileCache) GetLatestFavoriteID(screenName string) (int64, bool) {
	id, exists := c.LatestFavoriteID[screenName]
	return id, exists
}

func (c *FileCache) SetLatestFavoriteID(screenName string, id int64) error {
	c.LatestFavoriteID[screenName] = id
	return nil
}

func (c *FileCache) GetLatestDMID() int64 {
	return c.LatestDMID
}

func (c *FileCache) SetLatestDMID(id int64) error {
	c.LatestDMID = id
	return nil
}

func (c *FileCache) GetTweetAction(tweetID string) (*TwitterAction, bool) {
	action, exists := c.Tweet2Action[tweetID]
	return action, exists
}

func (c *FileCache) SetTweetAction(tweetID string, action *TwitterAction) error {
	c.Tweet2Action[tweetID] = action
	return nil
}

func (c *FileCache) GetLatestImages(num int) []ImageCacheData {
	if len(c.Images) >= num {
		return c.Images[len(c.Images)-num:]
	} else {
		return c.Images
	}
}

func (c *FileCache) SetImage(data ImageCacheData) error {
	c.Images = append(c.Images, data)
	return nil
}

// Save saves the cache data to the specified file
func (c *FileCache) Save() error {
	err := os.MkdirAll(filepath.Dir(c.File), 0600)
	if err != nil {
		return err
	}
	if c != nil {
		err := EncodeFile(c.File, c)
		if err != nil {
			return err
		}
	}
	return nil
}
