package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	defaultCachePath = os.ExpandEnv("$HOME/.cache/mybot/cache.json")
	cache            *mybotCache
)

type mybotCache struct {
	LatestCommitSHA map[string]map[string]string
	LatestTweetId   map[string]int64
	LatestDM        map[string]int64
}

func unmarshalCache(path string) error {
	if path == "" {
		path = defaultCachePath
	}

	if cache == nil {
		cache = &mybotCache{
			make(map[string]map[string]string),
			make(map[string]int64),
			make(map[string]int64),
		}
	}

	info, _ := os.Stat(path)
	if info != nil && !info.IsDir() {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, cache)
		if err != nil {
			return err
		}
	}
	return nil
}

func marshalCache(path string) error {
	var err error
	if path == "" {
		path = defaultCachePath
	}
	err = os.MkdirAll(filepath.Dir(path), 0600)
	if err != nil {
		return err
	}
	if cache != nil {
		data, err := json.Marshal(cache)
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
