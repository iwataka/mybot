package models

import (
	"gopkg.in/mgo.v2"
)

type MgoCollection interface {
	Find(query interface{}) MgoQuery
	RemoveAll(selector interface{}) (*mgo.ChangeInfo, error)
	Upsert(selector interface{}, update interface{}) (*mgo.ChangeInfo, error)
}

type DefaultMgoCollection struct {
	col *mgo.Collection
}

func NewMgoCollection(c *mgo.Collection) MgoCollection {
	return &DefaultMgoCollection{c}
}

func (c *DefaultMgoCollection) Find(query interface{}) MgoQuery {
	return c.col.Find(query)
}

func (c *DefaultMgoCollection) RemoveAll(selector interface{}) (*mgo.ChangeInfo, error) {
	return c.col.RemoveAll(selector)
}

func (c *DefaultMgoCollection) Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return c.col.Upsert(selector, update)
}

type MgoQuery interface {
	Count() (int, error)
	One(result interface{}) error
}
