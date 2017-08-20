package mybot

import (
	"github.com/iwataka/anaconda"
	"gopkg.in/mgo.v2"
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
}

func NewDBOAuthCreds(col *mgo.Collection) (*DBOAuthCreds, error) {
	a := &DBOAuthCreds{newOAuthCredsProps(), col}
	err := a.Load()
	return a, err
}

func (a *DBOAuthCreds) Load() error {
	query := a.col.Find(nil)
	count, err := query.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return query.One(a.OAuthCredsProps)
	}
	return nil
}

func (a *DBOAuthCreds) Save() error {
	_, err := a.col.Upsert(nil, a.OAuthCredsProps)
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

type OAuthAppProps struct {
	ConsumerKey    string `json:"consumer_key" toml:"consumer_key" bson:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret" toml:"consumer_secret" bson:"consumer_secret"`
}

func newOAuthAppProps() *OAuthAppProps {
	return &OAuthAppProps{"", ""}
}

func (a *TwitterOAuthAppProps) SetCreds(ck, cs string) {
	anaconda.SetConsumerKey(ck)
	anaconda.SetConsumerSecret(cs)
}

func (a *TwitterOAuthAppProps) GetCreds() (string, string) {
	return anaconda.GetConsumerKey(), anaconda.GetConsumerSecret()
}

type FileTwitterOAuthApp struct {
	*TwitterOAuthAppProps
	File string `json:"-" toml:"-" bson:"-"`
}

func NewFileTwitterOAuthApp(file string) (*FileTwitterOAuthApp, error) {
	a := &FileTwitterOAuthApp{newTwitterOAuthAppProps(), file}
	err := a.Load()
	return a, err
}

// Decode does nothing and returns nil if the specified file doesn't exist.
func (a *FileTwitterOAuthApp) Load() error {
	tmp := &OAuthAppProps{}
	err := DecodeFile(a.File, tmp)
	a.SetCreds(tmp.ConsumerKey, tmp.ConsumerSecret)
	return err
}

func (a *FileTwitterOAuthApp) Save() error {
	ck, cs := a.GetCreds()
	tmp := &OAuthAppProps{ck, cs}
	return EncodeFile(a.File, tmp)
}

type DBTwitterOAuthApp struct {
	*TwitterOAuthAppProps
	col *mgo.Collection `json:"-" toml:"-" bson:"-"`
}

func NewDBTwitterOAuthApp(col *mgo.Collection) (*DBTwitterOAuthApp, error) {
	a := &DBTwitterOAuthApp{newTwitterOAuthAppProps(), col}
	err := a.Load()
	return a, err
}

func (a *DBTwitterOAuthApp) Load() error {
	tmp := &OAuthAppProps{}
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

func (a *DBTwitterOAuthApp) Save() error {
	ck, cs := a.GetCreds()
	tmp := &OAuthAppProps{ck, cs}
	_, err := a.col.Upsert(nil, tmp)
	return err
}
