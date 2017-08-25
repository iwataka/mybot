package mybot

import (
	"github.com/iwataka/anaconda"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type OAuthCreds interface {
	SetCreds(at, ats string)
	GetCreds() (string, string)
	Load() error
	Save() error
}

// FileOAuthCreds contains values required for authenticated user
type FileOAuthCreds struct {
	*OAuthCredsProps
	File string `json:"-" toml:"-" bson:"-"`
}

// Load does nothing and returns nil if the specified file doesn't exist.
func (a *FileOAuthCreds) Load() error {
	return DecodeFile(a.File, a)
}

func (a *FileOAuthCreds) Save() error {
	return EncodeFile(a.File, a)
}

func NewFileOAuthCreds(file string) (*FileOAuthCreds, error) {
	a := &FileOAuthCreds{newOAuthCredsProps(), file}
	err := a.Load()
	return a, err
}

type DBOAuthCreds struct {
	*OAuthCredsProps
	col *mgo.Collection `json:"-" toml:"-" bson:"-"`
	ID  string          `json:"id" toml:"id" bson:"id"`
}

func NewDBOAuthCreds(col *mgo.Collection, id string) (*DBOAuthCreds, error) {
	a := &DBOAuthCreds{newOAuthCredsProps(), col, id}
	err := a.Load()
	return a, err
}

func (a *DBOAuthCreds) Load() error {
	query := a.col.Find(bson.M{"id": a.ID})
	count, err := query.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		tmpCol := a.col
		err := query.One(a)
		a.col = tmpCol
		return err
	}
	return nil
}

func (a *DBOAuthCreds) Save() error {
	_, err := a.col.Upsert(bson.M{"id": a.ID}, a)
	return err
}

type OAuthCredsProps struct {
	AccessToken       string `json:"access_token" toml:"access_token" bson:"access_token"`
	AccessTokenSecret string `json:"access_token_secret" toml:"access_token_secret" bson:"access_token_secret"`
}

func newOAuthCredsProps() *OAuthCredsProps {
	return &OAuthCredsProps{"", ""}
}

func (a *OAuthCredsProps) SetCreds(at, ats string) {
	a.AccessToken = at
	a.AccessTokenSecret = ats
}

func (a *OAuthCredsProps) GetCreds() (string, string) {
	return a.AccessToken, a.AccessTokenSecret
}

type OAuthApp interface {
	SetCreds(ck, cs string)
	GetCreds() (string, string)
	Load() error
	Save() error
}

type TwitterOAuthAppProps struct {
}

func newTwitterOAuthAppProps() *TwitterOAuthAppProps {
	return &TwitterOAuthAppProps{}
}

type OAuthAppProps interface {
	SetCreds(ck, cs string)
	GetCreds() (string, string)
}

type DefaultOAuthAppProps struct {
	ConsumerKey    string `json:"consumer_key" toml:"consumer_key" bson:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret" toml:"consumer_secret" bson:"consumer_secret"`
}

func (a *DefaultOAuthAppProps) SetCreds(ck, cs string) {
	a.ConsumerKey = ck
	a.ConsumerSecret = cs
}

func (a *DefaultOAuthAppProps) GetCreds() (string, string) {
	return a.ConsumerKey, a.ConsumerSecret
}

func newOAuthAppProps() *DefaultOAuthAppProps {
	return &DefaultOAuthAppProps{"", ""}
}

func (a *TwitterOAuthAppProps) SetCreds(ck, cs string) {
	anaconda.SetConsumerKey(ck)
	anaconda.SetConsumerSecret(cs)
}

func (a *TwitterOAuthAppProps) GetCreds() (string, string) {
	return anaconda.GetConsumerKey(), anaconda.GetConsumerSecret()
}

type FileOAuthApp struct {
	OAuthAppProps
	File string `json:"-" toml:"-" bson:"-"`
}

func NewFileOAuthApp(file string) (*FileOAuthApp, error) {
	a := &FileOAuthApp{&DefaultOAuthAppProps{}, file}
	err := a.Load()
	return a, err
}

func NewFileTwitterOAuthApp(file string) (*FileOAuthApp, error) {
	a := &FileOAuthApp{newTwitterOAuthAppProps(), file}
	err := a.Load()
	return a, err
}

// Decode does nothing and returns nil if the specified file doesn't exist.
func (a *FileOAuthApp) Load() error {
	tmp := &DefaultOAuthAppProps{}
	err := DecodeFile(a.File, tmp)
	a.SetCreds(tmp.ConsumerKey, tmp.ConsumerSecret)
	return err
}

func (a *FileOAuthApp) Save() error {
	ck, cs := a.GetCreds()
	tmp := &DefaultOAuthAppProps{ck, cs}
	return EncodeFile(a.File, tmp)
}

type DBOAuthApp struct {
	OAuthAppProps
	col *mgo.Collection `json:"-" toml:"-" bson:"-"`
}

func NewDBOAuthApp(col *mgo.Collection) (*DBOAuthApp, error) {
	a := &DBOAuthApp{&DefaultOAuthAppProps{}, col}
	err := a.Load()
	return a, err
}

func NewDBTwitterOAuthApp(col *mgo.Collection) (*DBOAuthApp, error) {
	a := &DBOAuthApp{newTwitterOAuthAppProps(), col}
	err := a.Load()
	return a, err
}

func (a *DBOAuthApp) Load() error {
	tmp := &DefaultOAuthAppProps{}
	query := a.col.Find(nil)
	count, err := query.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		query.One(tmp)
	}
	a.SetCreds(tmp.ConsumerKey, tmp.ConsumerSecret)
	return err
}

func (a *DBOAuthApp) Save() error {
	ck, cs := a.GetCreds()
	tmp := &DefaultOAuthAppProps{ck, cs}
	_, err := a.col.Upsert(nil, tmp)
	return err
}
