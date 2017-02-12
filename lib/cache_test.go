package mybot

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFileCacheSave(t *testing.T) {
	dir := os.TempDir()
	fname := "cache.json"
	path := filepath.Join(dir, fname)
	c, err := NewFileCache(path)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Save()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("%s expected to exist but not", path)
	}
}

func TestFileCacheLatestTweetID(t *testing.T) {
	cache, err := NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheLatestTweetID(t, cache)
}

func TestDBCacheLatestTweetID(t *testing.T) {
	c, err := NewDBCache("sqlite3", "test.db")
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
	err := c.SetLatestTweetID(name, tweetID)
	if err != nil {
		t.Fatal(err)
	}
	id, err := c.GetLatestTweetID(name)
	if err != nil {
		t.Fatal(err)
	}
	if id != tweetID {
		t.Fatalf("%v expected but %v found", tweetID, id)
	}
}

func TestFileCacheLatestFavoriteID(t *testing.T) {
	c, err := NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheLatestFavoriteID(t, c)
}

func TesDBCacheLatestFavoriteID(t *testing.T) {
	c, err := NewDBCache("sqlite3", "test.db")
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
	err := c.SetLatestFavoriteID(name, favoriteID)
	if err != nil {
		t.Fatal(err)
	}
	id, err := c.GetLatestFavoriteID(name)
	if err != nil {
		t.Fatal(err)
	}
	if id != favoriteID {
		t.Fatalf("%v expected but %v found", favoriteID, id)
	}
}

func TestFileCacheLatestDMID(t *testing.T) {
	c, err := NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheLatestDMID(t, c)
}

func TestDBCacheLatestDMID(t *testing.T) {
	c, err := NewDBCache("sqlite3", "test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.db")
	testCacheLatestDMID(t, c)
}

func testCacheLatestDMID(t *testing.T, c Cache) {
	var dmID int64
	dmID = 1
	err := c.SetLatestDMID(dmID)
	if err != nil {
		t.Fatal(err)
	}
	id, err := c.GetLatestDMID()
	if err != nil {
		t.Fatal(err)
	}
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
	c, err := NewDBCache("sqlite3", "test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.db")
	testCacheTweetAction(t, c)
}

func testCacheTweetAction(t *testing.T, c Cache) {
	var tweetID int64 = 1
	action := &TweetAction{
		Twitter: &TwitterAction{
			Collections: []string{"foo"},
		},
		Slack: NewSlackAction(),
	}
	action.Twitter.Retweet = true
	err := c.SetTweetAction(tweetID, action)
	if err != nil {
		t.Fatal(err)
	}
	a, err := c.GetTweetAction(tweetID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(action, a) {
		t.Fatalf("%v expected but %v found", action, a)
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
	c, err := NewDBCache("sqlite3", "test.db")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.db")
	testCacheImage(t, c)
}

func testCacheImage(t *testing.T, c Cache) {
	img := ImageCacheData{}
	err := c.SetImage(img)
	if err != nil {
		t.Fatal(err)
	}
	is, err := c.GetLatestImages(1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(img, is[0]) {
		t.Fatalf("%v expected but %v found", img, is[0])
	}
}
