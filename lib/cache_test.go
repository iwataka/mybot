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
	id := c.GetLatestTweetID(name)
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
	id := c.GetLatestFavoriteID(name)
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
	img := ImageCacheData{}
	c.SetImage(img)
	is := c.GetLatestImages(1)
	if !reflect.DeepEqual(img, is[0]) {
		t.Fatalf("%v expected but %v found", img, is[0])
	}
}
