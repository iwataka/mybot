package data

import (
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/utils"
	"gopkg.in/mgo.v2/bson"

	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// Cache provides a set of functions to access cache of this application.
type Cache interface {
	utils.Savable
	utils.Deletable
	GetLatestTweetID(screenName string) int64
	SetLatestTweetID(screenName string, id int64)
	GetLatestFavoriteID(screenName string) int64
	SetLatestFavoriteID(screenName string, id int64)
	GetLatestDMID() int64
	SetLatestDMID(id int64)
	GetTweetAction(tweetID int64) Action
	SetTweetAction(tweetID int64, action Action)
	GetLatestImages(num int) []models.ImageCacheData
	SetImage(data models.ImageCacheData)
}

// FileCache is a cache data associated with a specified file.
type FileCache struct {
	CacheProperties `yaml:",inline"`
	File            string `json:"-" toml:"-" bson:"-" yaml:"-"`
}

// NewFileCache returns a new FileCache.
// If no file specified, this returns an empty Cache instance, which has no
// data and can't save the content to any file.
func NewFileCache(path string) (*FileCache, error) {
	c := &FileCache{
		newCacheProperties(),
		path,
	}

	info, _ := os.Stat(path)
	if info != nil && !info.IsDir() {
		err := utils.DecodeFile(path, c)
		if err != nil {
			return nil, utils.WithStack(err)
		}
	}
	return c, nil
}

// Save saves the cache data to the specified file
func (c *FileCache) Save() error {
	c.CacheProperties.m.RLock()
	defer c.CacheProperties.m.RUnlock()
	err := os.MkdirAll(filepath.Dir(c.File), 0600)
	if err != nil {
		return utils.WithStack(err)
	}
	if c != nil {
		err := utils.EncodeFile(c.File, c)
		if err != nil {
			return utils.WithStack(err)
		}
	}
	return nil
}

func (c *FileCache) Delete() error {
	return os.RemoveAll(c.File)
}

type DBCache struct {
	CacheProperties `yaml:",inline"`
	col             models.MgoCollection
	ID              string `json:"id" toml:"id" bson:"id" yaml:"id"`
}

func NewDBCache(col models.MgoCollection, id string) (*DBCache, error) {
	c := &DBCache{newCacheProperties(), col, id}
	query := col.Find(bson.M{"id": c.ID})
	count, err := query.Count()
	if err != nil {
		return nil, utils.WithStack(err)
	}
	if count > 0 {
		// query.One overrides m so temporarily backup it.
		tmpMutex := c.m //nolint: vet
		err := query.One(c)
		if err != nil {
			return nil, utils.WithStack(err)
		}
		c.col, c.m = col, tmpMutex //nolint: vet
	}
	return c, utils.WithStack(err)
}

func (c *DBCache) Save() error {
	c.CacheProperties.m.RLock()
	defer c.CacheProperties.m.RUnlock()
	_, err := c.col.Upsert(bson.M{"id": c.ID}, c)
	return utils.WithStack(err)
}

func (c *DBCache) Delete() error {
	_, err := c.col.RemoveAll(bson.M{"id": c.ID})
	return err
}

// CacheProperties contains common actual cache variables and is intended to be
// embedded into other structs.
// All functions of this struct are thread-safe.
type CacheProperties struct {
	m sync.RWMutex
	// LatestTweetID associates Twitter screen name with the latest tweet
	// ID in timeline.
	LatestTweetID map[string]int64 `json:"latest_tweet_id" toml:"latest_tweet_id" bson:"latest_tweet_id" yaml:"latest_tweet_id"`
	// LatestFavoriteID associates Twitter screen name with the latest
	// tweet ID in favorite list.
	LatestFavoriteID map[string]int64 `json:"latest_favorite_id" toml:"lates_favorite_id" bson:"latest_favorite_id" yaml:"latest_favorite_id"`
	// LatestDMID is latest direct message ID of the authenticated user
	// with the latest direct message ID.
	LatestDMID int64 `json:"latest_dm_id" toml:"latest_dm_id" bson:"latest_dm_id" yaml:"latest_dm_id"`
	// Tweet2Action associates tweet ID with Mybot action.
	// This is not an instance of map[int64]Action because it can't be
	// converted to json when Go Runtime is v1.6 or older
	Tweet2Action map[string]Action `json:"tweet_to_action" toml:"tweet_to_action" bson:"tweet_to_action" yaml:"tweet_to_action"`
	// Images is cache data of images analyzed by some API or method.
	Images []models.ImageCacheData `json:"images" toml:"images" bson:"images" yaml:"images"`
}

func newCacheProperties() CacheProperties {
	return CacheProperties{
		LatestTweetID:    make(map[string]int64),
		LatestFavoriteID: make(map[string]int64),
		LatestDMID:       0,
		Tweet2Action:     make(map[string]Action),
		Images:           []models.ImageCacheData{},
	}
}

// GetLatestTweetID returns the latest tweet ID associated with screenName in
// timeline. If there is no ID of screenName in c , this returns 0 (tweet ID
// can't be 0, which is known by Twitter API specification).
func (c *CacheProperties) GetLatestTweetID(screenName string) int64 {
	c.m.RLock()
	defer c.m.RUnlock()
	id, exists := c.LatestTweetID[screenName]
	if exists {
		return id
	}
	return 0
}

// SetLatestTweetID stores id as the latest tweet ID and associates it with
// screenName.
func (c *CacheProperties) SetLatestTweetID(screenName string, id int64) {
	c.m.Lock()
	defer c.m.Unlock()
	c.LatestTweetID[screenName] = id
}

// GetLatestFavoriteID returns the latest favorite tweet ID of screenName.
// If there is no ID of screenName in c, this returns 0 (tweet ID can't be 0,
// which is known by Twitter API specification).
func (c *CacheProperties) GetLatestFavoriteID(screenName string) int64 {
	c.m.RLock()
	defer c.m.RUnlock()
	id, exists := c.LatestFavoriteID[screenName]
	if exists {
		return id
	}
	return 0
}

// SetLatestFavoriteID stores id as the latest favorite tweet ID and associates
// it with screeName.
func (c *CacheProperties) SetLatestFavoriteID(screenName string, id int64) {
	c.m.Lock()
	defer c.m.Unlock()
	c.LatestFavoriteID[screenName] = id
}

// GetLatestDMID returns latest direct message ID of the authenticated user.
func (c *CacheProperties) GetLatestDMID() int64 {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.LatestDMID
}

// SetLatestDMID sets id as the latest direct message ID of the authenticated
// user.
func (c *CacheProperties) SetLatestDMID(id int64) {
	c.m.Lock()
	defer c.m.Unlock()
	c.LatestDMID = id
}

// GetTweetAction returns Mybot action associated with tweetID.
func (c *CacheProperties) GetTweetAction(tweetID int64) Action {
	c.m.RLock()
	defer c.m.RUnlock()
	// Do not use string(tweetID) because it returns broken characters if
	// tweetID is too large
	key := strconv.FormatInt(tweetID, 10)
	action := c.Tweet2Action[key]
	return action
}

// SetTweetAction associates action with tweetID.
func (c *CacheProperties) SetTweetAction(tweetID int64, action Action) {
	c.m.Lock()
	defer c.m.Unlock()
	// Do not use string(tweetID) because it returns broken characters if
	// tweetID is enough large
	key := strconv.FormatInt(tweetID, 10)
	c.Tweet2Action[key] = action
}

// GetLatestImages returns the num latest pieces of cache image data.
func (c *CacheProperties) GetLatestImages(num int) []models.ImageCacheData {
	c.m.RLock()
	defer c.m.RUnlock()
	if len(c.Images) >= num && num > 0 {
		return c.Images[len(c.Images)-num:]
	} else {
		return c.Images[:]
	}
}

// SetImage sets data as the latest image cache data.
func (c *CacheProperties) SetImage(data models.ImageCacheData) {
	c.m.Lock()
	defer c.m.Unlock()
	c.Images = append(c.Images, data)
}
