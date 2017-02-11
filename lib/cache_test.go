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
	name := "foo"
	var favoriteID int64
	favoriteID = 1
	err = c.SetLatestFavoriteID(name, favoriteID)
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
	var dmID int64
	dmID = 1
	err = c.SetLatestDMID(dmID)
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
	var tweetID int64 = 1
	action := &TweetAction{
		Twitter: &TwitterAction{
			Retweet:     true,
			Favorite:    false,
			Follow:      false,
			Collections: []string{"Collection"},
		},
		Slack: NewSlackAction(),
	}
	err = c.SetTweetAction(tweetID, action)
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
	is, err := c.GetLatestImages(1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(img, is[0]) {
		t.Fatalf("%v expected but %v found", img, is[0])
	}
}

func TestDBCache(t *testing.T) {
	file := "test.db"
	if info, _ := os.Stat(file); info != nil && !info.IsDir() {
		os.Remove(file)
	}
	cache, err := NewDBCache("sqlite3", file)
	if err != nil {
		t.Fatal(err)
	}
	if cache.client.Error != nil {
		t.Fatal(cache.client.Error)
	}

	action := &TweetAction{
		Twitter: &TwitterAction{
			Retweet:     true,
			Collections: []string{"foo"},
		},
		Slack: &SlackAction{
			Channels: []string{"bar"},
		},
	}
	err = cache.SetTweetAction(7, action)
	if err != nil {
		t.Fatal(err)
	}
	ac, err := cache.GetTweetAction(7)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(action, ac) {
		t.Fatalf("Tweet Action %v is not set", action)
	}

	image := ImageCacheData{
		URL:            "url",
		Src:            "src",
		AnalysisResult: "result",
		AnalysisDate:   "date",
	}
	err = cache.SetImage(image)
	if err != nil {
		t.Fatal(err)
	}
	imgs, err := cache.GetLatestImages(1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(image, imgs[0]) {
		t.Fatalf("Image %v is not set", image)
	}

	err = cache.SetLatestFavoriteID("golang", 13)
	if err != nil {
		t.Fatal(err)
	}
	favID, err := cache.GetLatestFavoriteID("golang")
	if err != nil {
		t.Fatal(err)
	}
	if favID != 13 {
		t.Fatalf("Latest Favorite ID %v is not set", 13)
	}

	err = cache.SetLatestTweetID("golang", 17)
	if err != nil {
		t.Fatal(err)
	}
	twID, err := cache.GetLatestTweetID("golang")
	if err != nil {
		t.Fatal(err)
	}
	if twID != 17 {
		t.Fatalf("Latest Tweet ID %v is not set", 17)
	}

	err = cache.SetLatestDMID(1)
	if err != nil {
		t.Fatal(err)
	}
	dmID, err := cache.GetLatestDMID()
	if err != nil {
		t.Fatal(err)
	}
	if dmID != 1 {
		t.Fatalf("Latest DM ID %v is not set", 1)
	}
}
