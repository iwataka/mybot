package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// MybotCache is a cache data of mybot
type MybotCache struct {
	LatestCommitSHA     map[string]map[string]string
	LatestTweetID       map[string]int64
	LatestFavoriteID    map[string]int64
	LatestDMID          int64 `json:"-"`
	Tweet2Action        map[string]*TwitterAction
	ImageURL            string
	ImageSource         string
	ImageAnalysisResult string
	ImageAnalysisDate   string
}

// NewMybotCache creates an instance of MybotCache
func NewMybotCache(path string) (*MybotCache, error) {
	c := &MybotCache{
		make(map[string]map[string]string),
		make(map[string]int64),
		make(map[string]int64),
		0,
		make(map[string]*TwitterAction),
		"",
		"",
		"",
		"",
	}

	info, _ := os.Stat(path)
	if info != nil && !info.IsDir() {
		data, err := ioutil.ReadFile(path)

		// If the specified file is empty, returns empty cache
		if len(data) == 0 {
			return c, nil
		}

		if err != nil {
			return c, err
		}
		err = json.Unmarshal(data, c)
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

// Save saves the cache data to the specified file
func (c *MybotCache) Save(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0600)
	if err != nil {
		return err
	}
	if c != nil {
		data, err := json.Marshal(c)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, data, 0600)
		if err != nil {
			return err
		}
	}
	return nil
}
