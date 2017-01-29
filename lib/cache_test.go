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
	c.Save()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("%s expected to exist but not", path)
	}
}

func TestFileCacheLatestTweetID(t *testing.T) {
	c, err := NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	name := "foo"
	var tweetID int64
	tweetID = 1
	err = c.SetLatestTweetID(name, tweetID)
	if err != nil {
		t.Fatal(err)
	}
	id, exists := c.GetLatestTweetID(name)
	if !exists {
		t.Fatalf("Tweet ID does not exist")
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
	name := "foo"
	var favoriteID int64
	favoriteID = 1
	err = c.SetLatestFavoriteID(name, favoriteID)
	if err != nil {
		t.Fatal(err)
	}
	id, exists := c.GetLatestFavoriteID(name)
	if !exists {
		t.Fatalf("Favorite ID does not exist")
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
	var dmID int64
	dmID = 1
	err = c.SetLatestDMID(dmID)
	if err != nil {
		t.Fatal(err)
	}
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
	var tweetID string
	tweetID = "1"
	action := &TwitterAction{
		Retweet:     true,
		Favorite:    false,
		Follow:      false,
		Collections: []string{"Collection"},
	}
	err = c.SetTweetAction(tweetID, action)
	if err != nil {
		t.Fatal(err)
	}
	a, exists := c.GetTweetAction(tweetID)
	if !exists {
		t.Fatalf("Tweet ID does not exist")
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
	img := ImageCacheData{
		URL:            "",
		Src:            "",
		AnalysisResult: "",
		AnalysisDate:   "",
	}
	err = c.SetImage(img)
	if err != nil {
		t.Fatal(err)
	}
	i := c.GetLatestImages(1)[0]
	if !reflect.DeepEqual(img, i) {
		t.Fatalf("%v expected but %v found", img, i)
	}
}
