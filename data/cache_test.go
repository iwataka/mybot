package data

import (
	"github.com/iwataka/mybot/models"
	"github.com/stretchr/testify/require"

	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewFileCache_TakesInvalidJson(t *testing.T) {
	_, err := NewFileCache("testdata/cache_invalid.json")
	require.Error(t, err)
}

func TestFileCache_Save(t *testing.T) {
	var err error
	var path string
	var c Cache
	var fname string
	dir := os.TempDir()

	fname = "cache.json"
	path = filepath.Join(dir, fname)
	c, err = NewFileCache(path)
	require.NoError(t, err)
	require.NoError(t, c.Save())
	defer os.Remove(path)
	_, err = os.Stat(path)
	require.NoError(t, err)

	// Invalid path
	path = filepath.Join(path, fname)
	c, err = NewFileCache(path)
	require.NoError(t, err)
	require.Error(t, c.Save())
	defer os.Remove(path)

	// Unwritable file
	fname = "unwritable.json"
	path = filepath.Join(dir, fname)
	_, err = os.Create(path)
	require.NoError(t, err)
	defer os.Remove(path)
	require.NoError(t, os.Chmod(path, 0555))
	defer func() { _ = os.Chmod(path, 0755) }()
	c, err = NewFileCache(path)
	require.NoError(t, err)
	require.Error(t, c.Save())
}

func TestFileCache_LatestTweetID(t *testing.T) {
	c, err := NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheLatestTweetID(t, c)
}

func TestDBCache_LatestTweetID(t *testing.T) {
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
	var tweetID int64 = 1
	c.SetLatestTweetID(name, tweetID)
	require.Equal(t, tweetID, c.GetLatestTweetID(name))
	require.EqualValues(t, 0, c.GetLatestTweetID("bar"))
}

func TestFileCache_LatestFavoriteID(t *testing.T) {
	c, err := NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheLatestFavoriteID(t, c)
}

func TestDBCache_LatestFavoriteID(t *testing.T) {
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
	var favoriteID int64 = 1
	c.SetLatestFavoriteID(name, favoriteID)
	require.Equal(t, favoriteID, c.GetLatestFavoriteID(name))
	require.EqualValues(t, 0, c.GetLatestFavoriteID("bar"))
}

func TestFileCache_LatestDMID(t *testing.T) {
	c, err := NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheLatestDMID(t, c)
}

func TestDBCache_LatestDMID(t *testing.T) {
	t.Skip("You must write mocking test for this")
	c, err := NewDBCache(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("test.db")
	testCacheLatestDMID(t, c)
}

func testCacheLatestDMID(t *testing.T, c Cache) {
	var dmID int64 = 1
	c.SetLatestDMID(dmID)
	id := c.GetLatestDMID()
	require.Equal(t, dmID, id)
}

func TestFileCache_TweetAction(t *testing.T) {
	c, err := NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheTweetAction(t, c)
}

func TestDBCache_TweetAction(t *testing.T) {
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
		t.Fatalf("%#v expected but %#v found", action.Twitter, a.Twitter)
	}
	if !reflect.DeepEqual(action.Slack, a.Slack) {
		t.Fatalf("%#v expected but %#v found", action.Slack, a.Slack)
	}
}

func TestFileCache_Image(t *testing.T) {
	c, err := NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheImage(t, c)
}

func TestDBCache_Image(t *testing.T) {
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
	require.True(t, reflect.DeepEqual(expected, c.GetLatestImages(1)))
	require.True(t, reflect.DeepEqual(expected, c.GetLatestImages(2)))
}
