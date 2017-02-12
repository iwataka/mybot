package mybot

import (
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iwataka/mybot/models"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Cache interface {
	GetLatestTweetID(screenName string) (int64, error)
	SetLatestTweetID(screenName string, id int64) error
	GetLatestFavoriteID(screenName string) (int64, error)
	SetLatestFavoriteID(screenName string, id int64) error
	GetLatestDMID() (int64, error)
	SetLatestDMID(id int64) error
	GetTweetAction(tweetID int64) (*TweetAction, error)
	SetTweetAction(tweetID int64, action *TweetAction) error
	GetLatestImages(num int) ([]ImageCacheData, error)
	SetImage(data ImageCacheData) error
	Save() error
}

// FileCache is a cache data stored in the specified file.
type FileCache struct {
	LatestTweetID    map[string]int64 `json:"latest_tweet_id" toml:"latest_tweet_id"`
	LatestFavoriteID map[string]int64 `json:"latest_favorite_id" toml:"lates_favorite_id"`
	LatestDMID       int64            `json:"latest_dm_id" toml:"latest_dm_id"`
	// map[int64]interface{} can't be converted to json by go1.6 or older
	Tweet2Action map[string]*TweetAction `json:"tweet_to_action" toml:"tweet_to_action"`
	Images       []ImageCacheData        `json:"images" toml:"images"`
	File         string                  `json:"-" toml:"-"`
}

type ImageCacheData struct {
	models.VisionCacheProperties
}

// NewFileCache creates a Cache instance by using the specified file.
// If no file specified, this returns an emtpy Cache instance which doesn't
// have read/write features.
func NewFileCache(path string) (*FileCache, error) {
	c := &FileCache{
		make(map[string]int64),
		make(map[string]int64),
		0,
		make(map[string]*TweetAction),
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

func (c *FileCache) GetLatestTweetID(screenName string) (int64, error) {
	id, exists := c.LatestTweetID[screenName]
	if exists {
		return id, nil
	}
	return 0, nil
}

func (c *FileCache) SetLatestTweetID(screenName string, id int64) error {
	c.LatestTweetID[screenName] = id
	return nil
}

func (c *FileCache) GetLatestFavoriteID(screenName string) (int64, error) {
	id, exists := c.LatestFavoriteID[screenName]
	if exists {
		return id, nil
	}
	return 0, nil
}

func (c *FileCache) SetLatestFavoriteID(screenName string, id int64) error {
	c.LatestFavoriteID[screenName] = id
	return nil
}

func (c *FileCache) GetLatestDMID() (int64, error) {
	return c.LatestDMID, nil
}

func (c *FileCache) SetLatestDMID(id int64) error {
	c.LatestDMID = id
	return nil
}

func (c *FileCache) GetTweetAction(tweetID int64) (*TweetAction, error) {
	action, _ := c.Tweet2Action[string(tweetID)]
	return action, nil
}

func (c *FileCache) SetTweetAction(tweetID int64, action *TweetAction) error {
	c.Tweet2Action[string(tweetID)] = action
	return nil
}

func (c *FileCache) GetLatestImages(num int) ([]ImageCacheData, error) {
	if len(c.Images) >= num {
		return c.Images[len(c.Images)-num:], nil
	} else {
		return c.Images, nil
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

type DBCache struct {
	client *gorm.DB
}

func NewDBCache(driverName, dataSource string) (*DBCache, error) {
	db, err := gorm.Open(driverName, dataSource)
	return &DBCache{db}, err
}

func (c *DBCache) GetLatestTweetID(screenName string) (int64, error) {
	record := &models.TwitterUserCache{}
	c.client.AutoMigrate(record)
	c.client.Where("screen_name = ?", screenName).First(record)
	if record == nil {
		return 0, c.client.Error
	}
	return record.LatestTweetID, c.client.Error
}

func (c *DBCache) SetLatestTweetID(screenName string, id int64) error {
	cache := &models.Cache{}
	c.client.AutoMigrate(cache)
	c.client.FirstOrCreate(cache)

	record := &models.TwitterUserCache{}
	c.client.AutoMigrate(record)
	c.client.Where("screen_name = ?", screenName).FirstOrCreate(record)
	record.ScreenName = screenName
	record.LatestTweetID = id
	record.CacheID = cache.ID

	c.client.Save(record)
	return c.client.Error
}

func (c *DBCache) GetLatestFavoriteID(screenName string) (int64, error) {
	record := &models.TwitterUserCache{}
	c.client.AutoMigrate(record)
	c.client.Where("screen_name = ?", screenName).First(record)
	if record == nil {
		return 0, c.client.Error
	}
	return record.LatestFavoriteID, c.client.Error
}

func (c *DBCache) SetLatestFavoriteID(screenName string, id int64) error {
	cache := &models.Cache{}
	c.client.AutoMigrate(cache)
	c.client.FirstOrCreate(cache)

	record := &models.TwitterUserCache{}
	c.client.AutoMigrate(record)
	c.client.Where("screen_name = ?", screenName).FirstOrCreate(record)
	record.ScreenName = screenName
	record.LatestFavoriteID = id
	record.CacheID = cache.ID

	c.client.Save(record)
	return c.client.Error
}

func (c *DBCache) GetLatestDMID() (int64, error) {
	record := &models.Cache{}
	c.client.AutoMigrate(record)
	c.client.First(record)
	if record == nil {
		return 0, c.client.Error
	}
	return record.LatestDMID, c.client.Error
}

func (c *DBCache) SetLatestDMID(id int64) error {
	record := &models.Cache{}
	c.client.AutoMigrate(record)
	c.client.FirstOrCreate(record)
	record.LatestDMID = id
	c.client.Save(record)
	return c.client.Error
}

func (c *DBCache) GetTweetAction(tweetID int64) (*TweetAction, error) {
	record := &models.TweetActionCache{}
	c.client.AutoMigrate(record)
	c.client.Where("tweet_id = ?", tweetID).First(record)
	if record == nil {
		return nil, c.client.Error
	}

	twitter, err := c.getTwitterAction(record.TwitterActionID)
	if err != nil {
		return nil, err
	}
	record.TwitterAction = *twitter
	slack, err := c.getSlackAction(record.SlackActionID)
	if err != nil {
		return nil, err
	}
	record.SlackAction = *slack

	result := NewTweetAction()
	result.Twitter.TwitterActionProperties = record.TwitterAction.TwitterActionProperties
	result.Twitter.Collections = record.TwitterAction.GetCollections()
	result.Slack.Channels = record.SlackAction.GetChannels()
	return result, c.client.Error
}

func (c *DBCache) getTwitterAction(id uint) (*models.TwitterAction, error) {
	record := &models.TwitterAction{}
	c.client.AutoMigrate(record)
	c.client.First(record, id)
	cols, err := c.getTwitterCollections(id)
	if err != nil {
		return nil, err
	}
	record.Collections = cols
	return record, c.client.Error
}

func (c *DBCache) getTwitterCollections(id uint) ([]models.TwitterCollection, error) {
	records := []models.TwitterCollection{}
	c.client.AutoMigrate(&models.TwitterCollection{})
	c.client.Where("twitter_action_id = ?", id).Find(&records)
	return records, c.client.Error
}

func (c *DBCache) getSlackAction(id uint) (*models.SlackAction, error) {
	record := &models.SlackAction{}
	c.client.AutoMigrate(record)
	c.client.First(record, id)
	chs, err := c.getSlackChannels(id)
	if err != nil {
		return nil, err
	}
	record.Channels = chs
	return record, c.client.Error
}

func (c *DBCache) getSlackChannels(id uint) ([]models.SlackChannel, error) {
	records := []models.SlackChannel{}
	c.client.AutoMigrate(&models.SlackChannel{})
	c.client.Where("slack_action_id = ?", id).Find(&records)
	return records, c.client.Error
}

func (c *DBCache) SetTweetAction(tweetID int64, action *TweetAction) error {
	cache := &models.Cache{}
	c.client.AutoMigrate(cache)
	c.client.FirstOrCreate(cache)

	c.client.AutoMigrate(&models.TwitterAction{})
	c.client.AutoMigrate(&models.SlackAction{})
	c.client.AutoMigrate(&models.TwitterCollection{})
	c.client.AutoMigrate(&models.SlackChannel{})
	record := &models.TweetActionCache{}
	c.client.AutoMigrate(record)

	c.client.Where("tweet_id = ?", tweetID).FirstOrCreate(record)
	record.CacheID = cache.ID
	record.TweetID = tweetID
	record.TwitterAction.Retweet = action.Twitter.Retweet
	record.TwitterAction.Favorite = action.Twitter.Favorite
	record.TwitterAction.Follow = action.Twitter.Follow
	record.TwitterAction.SetCollections(action.Twitter.Collections)
	record.SlackAction.SetChannels(action.Slack.Channels)

	c.client.Save(record)
	for _, col := range record.TwitterAction.Collections {
		col.TwitterActionID = record.TwitterAction.ID
		c.client.Save(col)
	}
	for _, ch := range record.SlackAction.Channels {
		ch.SlackActionID = record.SlackAction.ID
		c.client.Save(ch)
	}
	return c.client.Error
}

func (c *DBCache) GetLatestImages(num int) ([]ImageCacheData, error) {
	record := &models.VisionCache{}
	c.client.AutoMigrate(record)
	records := []models.VisionCache{}
	c.client.Order("id desc").Limit(num).Find(&records)
	results := []ImageCacheData{}
	for _, r := range records {
		img := ImageCacheData{}
		img.URL = r.URL
		img.Src = r.Src
		img.AnalysisResult = r.AnalysisResult
		img.AnalysisDate = r.AnalysisDate
		results = append(results, img)
	}

	return results, c.client.Error
}

func (c *DBCache) SetImage(data ImageCacheData) error {
	record := &models.VisionCache{}
	c.client.AutoMigrate(record)
	record.URL = data.URL
	record.Src = data.Src
	record.AnalysisResult = data.AnalysisResult
	record.AnalysisDate = data.AnalysisDate
	c.client.Save(record)

	return c.client.Error
}

func (c *DBCache) Save() error {
	c.client.Commit()
	return c.client.Error
}
