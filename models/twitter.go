package models

import (
	"github.com/jinzhu/gorm"
)

type TwitterAction struct {
	gorm.Model
	TwitterActionProperties
	Collections []TwitterCollection
}

func (a *TwitterAction) GetCollections() []string {
	result := []string{}
	for _, c := range a.Collections {
		result = append(result, c.Name)
	}
	return result
}

func (a *TwitterAction) SetCollections(cols []string) {
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
	Retweet  bool `json:"retweet" toml:"retweet"`
	Favorite bool `json:"favorite" toml:"favorite"`
	Follow   bool `json:"follow" toml:"follow"`
}
