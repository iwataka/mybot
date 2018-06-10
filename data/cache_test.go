package data_test

import (
	. "github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewFileCache_TakesInvalidJson(t *testing.T) {
	_, err := NewFileCache("testdata/cache_invalid.json")
	assert.Error(t, err)
}

func TestFileCacheSave(t *testing.T) {
	var err error
	var path string
	var c Cache
	var fname string
	dir := os.TempDir()

	fname = "cache.json"
	path = filepath.Join(dir, fname)
	c, err = NewFileCache(path)
	assert.NoError(t, err)
	assert.NoError(t, c.Save())
	defer os.Remove(path)
	_, err = os.Stat(path)
	assert.NoError(t, err)

	// Invalid path
	path = filepath.Join(path, fname)
	c, err = NewFileCache(path)
	assert.NoError(t, err)
	assert.Error(t, c.Save())
	defer os.Remove(path)

	// Unwritable file
	fname = "unwritable.json"
	path = filepath.Join(dir, fname)
	_, err = os.Create(path)
	require.NoError(t, err)
	defer os.Remove(path)
	require.NoError(t, os.Chmod(path, 0555))
	defer os.Chmod(path, 0755)
	c, err = NewFileCache(path)
	assert.NoError(t, err)
	assert.Error(t, c.Save())
}

func TestFileCacheLatestTweetID(t *testing.T) {
	c, err := NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheLatestTweetID(t, c)
}

func TestDBCacheLatestTweetID(t *testing.T) {
	t.Skip("You must write mocking test for this")
	c, err := NewDBCache(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.db")
	testCacheLatestTweetID(t, c)
}

func testCacheLatestTweetID(t *testing.T, c Cache) {
	name := "foo"
	var tweetID int64
	tweetID = 1
	c.SetLatestTweetID(name, tweetID)
	assert.Equal(t, tweetID, c.GetLatestTweetID(name))
	assert.EqualValues(t, 0, c.GetLatestTweetID("bar"))
}

func TestFileCacheLatestFavoriteID(t *testing.T) {
	c, err := NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheLatestFavoriteID(t, c)
}

func TestDBCacheLatestFavoriteID(t *testing.T) {
	t.Skip("You must write mocking test for this")
	c, err := NewDBCache(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.db")
	testCacheLatestFavoriteID(t, c)
}

func testCacheLatestFavoriteID(t *testing.T, c Cache) {
	name := "foo"
	var favoriteID int64
	favoriteID = 1
	c.SetLatestFavoriteID(name, favoriteID)
	assert.Equal(t, favoriteID, c.GetLatestFavoriteID(name))
	assert.EqualValues(t, 0, c.GetLatestFavoriteID("bar"))
}

func TestFileCacheLatestDMID(t *testing.T) {
	c, err := NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheLatestDMID(t, c)
}

func TestDBCacheLatestDMID(t *testing.T) {
	t.Skip("You must write mocking test for this")
	c, err := NewDBCache(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.db")
	testCacheLatestDMID(t, c)
}

func testCacheLatestDMID(t *testing.T, c Cache) {
	var dmID int64
	dmID = 1
	c.SetLatestDMID(dmID)
	id := c.GetLatestDMID()
	if id != dmID {
		t.Fatalf("%v expected but %v found", dmID, id)
	}
}

func TestFileCacheTweetAction(t *testing.T) {
	c, err := NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheTweetAction(t, c)
}

func TestDBCacheTweetAction(t *testing.T) {
	t.Skip("You must write mocking test for this")
	c, err := NewDBCache(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.db")
	testCacheTweetAction(t, c)
}

func testCacheTweetAction(t *testing.T, c Cache) {
	var tweetID int64 = 1
	action := Action{
		Twitter: TwitterAction{
			Collections: []string{"foo"},
		},
		Slack: SlackAction{
			Reactions: []string{"smile"},
			Channels:  []string{"bar"},
		},
	}
	action.Twitter.Retweet = true
	action.Slack.Pin = true
	c.SetTweetAction(tweetID, action)
	a := c.GetTweetAction(tweetID)
	if !reflect.DeepEqual(action.Twitter, a.Twitter) {
		t.Fatalf("%v expected but %v found", action.Twitter, a.Twitter)
	}
	if !reflect.DeepEqual(action.Slack, a.Slack) {
		t.Fatalf("%v expected but %v found", action.Slack, a.Slack)
	}
}

func TestFileCacheImage(t *testing.T) {
	c, err := NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheImage(t, c)
}

func TestDBCacheImage(t *testing.T) {
	t.Skip("You must write mocking test for this")
	c, err := NewDBCache(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.db")
	testCacheImage(t, c)
}

func testCacheImage(t *testing.T, c Cache) {
	img := models.ImageCacheData{}
	c.SetImage(img)
	expected := []models.ImageCacheData{img}
	assert.True(t, reflect.DeepEqual(expected, c.GetLatestImages(1)))
	assert.True(t, reflect.DeepEqual(expected, c.GetLatestImages(2)))
}
