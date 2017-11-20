package models

import (
	"net/url"

	"github.com/iwataka/anaconda"
	"github.com/jinzhu/gorm"
)

type TwitterAPI interface {
	VerifyCredentials() (bool, error)
	PostDMToScreenName(string, string) (anaconda.DirectMessage, error)
	GetCollectionListByUserId(int64, url.Values) (anaconda.CollectionListResult, error)
	PostTweet(string, url.Values) (anaconda.Tweet, error)
	GetFriendsList(url.Values) (anaconda.UserCursor, error)
	GetSelf(url.Values) (anaconda.User, error)
	GetUserTimeline(url.Values) ([]anaconda.Tweet, error)
	GetFavorites(url.Values) ([]anaconda.Tweet, error)
	GetSearch(string, url.Values) (anaconda.SearchResponse, error)
	Retweet(int64, bool) (anaconda.Tweet, error)
	Favorite(int64) (anaconda.Tweet, error)
	CreateCollection(string, url.Values) (anaconda.CollectionShowResult, error)
	AddEntryToCollection(string, int64, url.Values) (anaconda.CollectionEntryAddResult, error)
	GetUsersLookup(string, url.Values) ([]anaconda.User, error)
	PublicStreamFilter(url.Values) *anaconda.Stream
	UserStream(url.Values) *anaconda.Stream
	GetUsersShow(string, url.Values) (anaconda.User, error)
}

type TwitterAction struct {
	gorm.Model
	TwitterActionProperties
	Collections []TwitterCollection
}

func (a TwitterAction) GetCollections() []string {
	result := []string{}
	for _, c := range a.Collections {
		result = append(result, c.Name)
	}
	return result
}

func (a TwitterAction) SetCollections(cols []string) {
	a.Collections = []TwitterCollection{}
	for _, col := range cols {
		c := TwitterCollection{
			TwitterActionID: a.ID,
			Name:            col,
		}
		a.Collections = append(a.Collections, c)
	}
}

type TwitterActionProperties struct {
	Tweet    bool `json:"tweet" toml:"tweet" bson:"tweet"`
	Retweet  bool `json:"retweet" toml:"retweet" bson:"retweet"`
	Favorite bool `json:"favorite" toml:"favorite" bson:"favorite"`
}
