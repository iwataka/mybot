package data_test

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/mybot/data"
	"github.com/iwataka/mybot/mocks"
	"github.com/iwataka/mybot/models"
	"github.com/stretchr/testify/require"

	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewFileCache_TakesInvalidJson(t *testing.T) {
	_, err := data.NewFileCache("testdata/cache_invalid.json")
	require.Error(t, err)
}

func TestFileCache_Save(t *testing.T) {
	dir := os.TempDir()
	fname := "cache.json"
	fpath := filepath.Join(dir, fname)
	c, err := data.NewFileCache(fpath)
	require.NoError(t, err)
	require.NoError(t, c.Save())
	defer func() { require.NoError(t, c.Delete()) }()
	_, err = os.Stat(fpath)
	require.NoError(t, err)

	invalidPath := filepath.Join(fpath, fname)
	invalidCache, err := data.NewFileCache(invalidPath)
	require.NoError(t, err)
	require.Error(t, invalidCache.Save())
}

func TestFileCache_Save_withUnwritablePath(t *testing.T) {
	dir := os.TempDir()
	unwritableFpath := filepath.Join(dir, "unwritable.json")
	_, err := os.Create(unwritableFpath)
	require.NoError(t, err)
	require.NoError(t, os.Chmod(unwritableFpath, 0555))
	unwritableCache, err := data.NewFileCache(unwritableFpath)
	require.NoError(t, err)
	defer func() { require.NoError(t, unwritableCache.Delete()) }()
	defer func() { require.NoError(t, os.Chmod(unwritableFpath, 0755)) }()
	require.Error(t, unwritableCache.Save())
}

func TestDBCache_Save_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	col := mocks.NewMockMgoCollection(ctrl)
	query := mocks.NewMockMgoQuery(ctrl)
	query.EXPECT().Count().Return(0, nil)
	col.EXPECT().Find(gomock.Any()).Return(query)
	c, err := data.NewDBCache(col, "foo")
	require.NoError(t, err)
	col.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil, nil)
	require.NoError(t, c.Save())
	col.EXPECT().RemoveAll(gomock.Any()).Return(nil, nil)
	require.NoError(t, c.Delete())
}

func TestNewDBCache_withExistingData(t *testing.T) {
	ctrl := gomock.NewController(t)
	col := mocks.NewMockMgoCollection(ctrl)
	query := mocks.NewMockMgoQuery(ctrl)
	query.EXPECT().Count().Return(1, nil)
	query.EXPECT().One(gomock.Any()).Return(nil)
	col.EXPECT().Find(gomock.Any()).Return(query)
	_, err := data.NewDBCache(col, "foo")
	require.NoError(t, err)
}

func TestNewDBCache_withCountError(t *testing.T) {
	ctrl := gomock.NewController(t)
	col := mocks.NewMockMgoCollection(ctrl)
	query := mocks.NewMockMgoQuery(ctrl)
	query.EXPECT().Count().Return(0, errors.New("error"))
	col.EXPECT().Find(gomock.Any()).Return(query)
	_, err := data.NewDBCache(col, "foo")
	require.Error(t, err)
}

func TestNewDBCache_withQueryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	col := mocks.NewMockMgoCollection(ctrl)
	query := mocks.NewMockMgoQuery(ctrl)
	query.EXPECT().Count().Return(1, nil)
	query.EXPECT().One(gomock.Any()).Return(errors.New("error"))
	col.EXPECT().Find(gomock.Any()).Return(query)
	_, err := data.NewDBCache(col, "foo")
	require.Error(t, err)
}

func TestCache_LatestTweetID(t *testing.T) {
	c, err := data.NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheLatestTweetID(t, c)
}

func testCacheLatestTweetID(t *testing.T, c data.Cache) {
	name := "foo"
	var tweetID int64 = 1
	c.SetLatestTweetID(name, tweetID)
	require.Equal(t, tweetID, c.GetLatestTweetID(name))
	require.EqualValues(t, 0, c.GetLatestTweetID("bar"))
}

func TestCache_LatestFavoriteID(t *testing.T) {
	c, err := data.NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheLatestFavoriteID(t, c)
}

func testCacheLatestFavoriteID(t *testing.T, c data.Cache) {
	name := "foo"
	var favoriteID int64 = 1
	c.SetLatestFavoriteID(name, favoriteID)
	require.Equal(t, favoriteID, c.GetLatestFavoriteID(name))
	require.EqualValues(t, 0, c.GetLatestFavoriteID("bar"))
}

func TestCache_LatestDMID(t *testing.T) {
	c, err := data.NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheLatestDMID(t, c)
}

func testCacheLatestDMID(t *testing.T, c data.Cache) {
	var dmID int64 = 1
	c.SetLatestDMID(dmID)
	id := c.GetLatestDMID()
	require.Equal(t, dmID, id)
}

func TestCache_TweetAction(t *testing.T) {
	c, err := data.NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheTweetAction(t, c)
}

func testCacheTweetAction(t *testing.T, c data.Cache) {
	var tweetID int64 = 1
	action := data.Action{
		Twitter: data.TwitterAction{
			Collections: []string{"foo"},
		},
		Slack: data.SlackAction{
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

func TestCache_Image(t *testing.T) {
	c, err := data.NewFileCache("cache.json")
	if err != nil {
		t.Fatal(err)
	}
	testCacheImage(t, c)
}

func testCacheImage(t *testing.T, c data.Cache) {
	img := models.ImageCacheData{}
	c.SetImage(img)
	expected := []models.ImageCacheData{img}
	require.True(t, reflect.DeepEqual(expected, c.GetLatestImages(1)))
	require.True(t, reflect.DeepEqual(expected, c.GetLatestImages(2)))
}
