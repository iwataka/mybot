package mybot

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/iwataka/mybot/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Cache interface {
	Savable
	GetLatestTweetID(screenName string) int64
	SetLatestTweetID(screenName string, id int64)
	GetLatestFavoriteID(screenName string) int64
	SetLatestFavoriteID(screenName string, id int64)
	GetLatestDMID() int64
	SetLatestDMID(id int64)
	GetTweetAction(tweetID int64) Action
	SetTweetAction(tweetID int64, action Action)
	GetLatestImages(num int) []ImageCacheData
	SetImage(data ImageCacheData)
}

type CacheProperties struct {
	LatestTweetID    map[string]int64 `json:"latest_tweet_id" toml:"latest_tweet_id" bson:"latest_tweet_id"`
	LatestFavoriteID map[string]int64 `json:"latest_favorite_id" toml:"lates_favorite_id" bson:"latest_favorite_id"`
	LatestDMID       int64            `json:"latest_dm_id" toml:"latest_dm_id" bson:"latest_dm_id"`
	// map[int64]interface{} can't be converted to json by go1.6 or older
	Tweet2Action map[string]Action `json:"tweet_to_action" toml:"tweet_to_action" bson:"tweet_to_action"`
	Images       []ImageCacheData  `json:"images" toml:"images" bson:"images"`
}

func newCacheProperties() CacheProperties {
	return CacheProperties{
		make(map[string]int64),
		make(map[string]int64),
		0,
		make(map[string]Action),
		[]ImageCacheData{},
	}
}

// FileCache is a cache data stored in the specified file.
type FileCache struct {
	CacheProperties
	File string `json:"-" toml:"-" bson:"-"`
}

type ImageCacheData struct {
	models.VisionCacheProperties
}

// NewFileCache creates a Cache instance by using the specified file.
// If no file specified, this returns an emtpy Cache instance which doesn't
// have read/write features.
func NewFileCache(path string) (*FileCache, error) {
	c := &FileCache{
		newCacheProperties(),
		path,
	}

	info, _ := os.Stat(path)
	if info != nil && !info.IsDir() {
		err := DecodeFile(path, c)
		if err != nil {
			return nil, WithStack(err)
		}
	}
	return c, nil
}

func (c *CacheProperties) GetLatestTweetID(screenName string) int64 {
	id, exists := c.LatestTweetID[screenName]
	if exists {
		return id
	}
	return 0
}

func (c *CacheProperties) SetLatestTweetID(screenName string, id int64) {
	c.LatestTweetID[screenName] = id
}

func (c *CacheProperties) GetLatestFavoriteID(screenName string) int64 {
	id, exists := c.LatestFavoriteID[screenName]
	if exists {
		return id
	}
	return 0
}

func (c *CacheProperties) SetLatestFavoriteID(screenName string, id int64) {
	c.LatestFavoriteID[screenName] = id
}

func (c *CacheProperties) GetLatestDMID() int64 {
	return c.LatestDMID
}

func (c *CacheProperties) SetLatestDMID(id int64) {
	c.LatestDMID = id
}

func (c *CacheProperties) GetTweetAction(tweetID int64) Action {
	// Do not use string(tweetID) because it returns broken characters if
	// tweetID is enough large
	key := strconv.FormatInt(tweetID, 10)
	action, _ := c.Tweet2Action[key]
	return action
}

func (c *CacheProperties) SetTweetAction(tweetID int64, action Action) {
	// Do not use string(tweetID) because it returns broken characters if
	// tweetID is enough large
	key := strconv.FormatInt(tweetID, 10)
	c.Tweet2Action[key] = action
}

func (c *CacheProperties) GetLatestImages(num int) []ImageCacheData {
	if len(c.Images) >= num {
		return c.Images[len(c.Images)-num:]
	} else {
		return c.Images
	}
}

func (c *CacheProperties) SetImage(data ImageCacheData) {
	c.Images = append(c.Images, data)
}

// Save saves the cache data to the specified file
func (c *FileCache) Save() error {
	err := os.MkdirAll(filepath.Dir(c.File), 0600)
	if err != nil {
		return WithStack(err)
	}
	if c != nil {
		err := EncodeFile(c.File, c)
		if err != nil {
			return WithStack(err)
		}
	}
	return nil
}

type DBCache struct {
	CacheProperties
	col *mgo.Collection `json:"-" toml:"-" bson:"-"`
	ID  string          `json:"id" toml:"id" bson:"id"`
}

func NewDBCache(col *mgo.Collection, id string) (*DBCache, error) {
	c := &DBCache{newCacheProperties(), col, id}
	query := col.Find(bson.M{"id": c.ID})
	count, err := query.Count()
	if err != nil {
		return nil, WithStack(err)
	}
	if count > 0 {
		err := query.One(c)
		if err != nil {
			return nil, WithStack(err)
		}
		c.col = col
	}
	return c, WithStack(err)
}

func (c *DBCache) Save() error {
	_, err := c.col.Upsert(bson.M{"id": c.ID}, c)
	return WithStack(err)
}
