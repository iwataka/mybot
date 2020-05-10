/*
Package oauth provides containers for OAuth data.
*/
package oauth

// TODO: Encrypt credential information

import (
	"os"

	"github.com/iwataka/anaconda"
	"github.com/iwataka/mybot/models"
	"github.com/iwataka/mybot/utils"
	"gopkg.in/mgo.v2/bson"
)

// OAuthCreds providves features to save/load and set/get user credential
// information for OAuth.
type OAuthCreds interface {
	utils.Savable
	utils.Loadable
	utils.Deletable
	SetCreds(at, ats string)
	GetCreds() (string, string)
}

// FileOAuthCreds is a user credential container for OAuth and associated with
// a specified file.
type FileOAuthCreds struct {
	OAuthCredsProps `yaml:",inline"`
	File            string `json:"-" toml:"-" bson:"-" yaml:"-"`
}

// Load retrieves credential information from a specified file and stores it.
// This does nothing and returns nil if the file doesn't exist.
func (a *FileOAuthCreds) Load() error {
	return utils.DecodeFile(a.File, a)
}

// Save saves the credential information in a into a.File.
func (a *FileOAuthCreds) Save() error {
	return utils.EncodeFile(a.File, a)
}

func (a *FileOAuthCreds) Delete() error {
	return os.RemoveAll(a.File)
}

// NewFileOAuthCreds returns a new FileOAuthCreds with file.
func NewFileOAuthCreds(file string) (*FileOAuthCreds, error) {
	a := &FileOAuthCreds{newOAuthCredsProps(), file}
	err := a.Load()
	return a, utils.WithStack(err)
}

// DBOAuthCreds is a user credential container for OAuth and associated with
// a specified database (currently only MongoDB supported).
type DBOAuthCreds struct {
	OAuthCredsProps `yaml:",inline"`
	col             models.MgoCollection
	ID              string `json:"id" toml:"id" bson:"id" yaml:"id"`
}

// NewDBOAuthCreds returns a new DBOAuthCreds with specified arguments.
func NewDBOAuthCreds(col models.MgoCollection, id string) (*DBOAuthCreds, error) {
	a := &DBOAuthCreds{newOAuthCredsProps(), col, id}
	err := a.Load()
	return a, utils.WithStack(err)
}

// Load retrieves credential information from a specified database and stores
// it.
func (a *DBOAuthCreds) Load() error {
	query := a.col.Find(bson.M{"id": a.ID})
	count, err := query.Count()
	if err != nil {
		return utils.WithStack(err)
	}
	if count > 0 {
		tmpCol := a.col
		err := query.One(a)
		a.col = tmpCol
		return utils.WithStack(err)
	}
	return nil
}

// Save saves the credential information to a specified database.
func (a *DBOAuthCreds) Save() error {
	_, err := a.col.Upsert(bson.M{"id": a.ID}, a)
	return utils.WithStack(err)
}

func (a *DBOAuthCreds) Delete() error {
	_, err := a.col.RemoveAll(bson.M{"id": a.ID})
	return err
}

// OAuthCredsProps contains actual variables for user credential information.
type OAuthCredsProps struct {
	AccessToken       string `json:"access_token" toml:"access_token" bson:"access_token" yaml:"access_token"`
	AccessTokenSecret string `json:"access_token_secret" toml:"access_token_secret" bson:"access_token_secret" yaml:"access_token_secret"`
}

func newOAuthCredsProps() OAuthCredsProps {
	return OAuthCredsProps{"", ""}
}

// SetCreds sets user credential information to a.
func (a *OAuthCredsProps) SetCreds(at, ats string) {
	a.AccessToken = at
	a.AccessTokenSecret = ats
}

// GetCreds returns user credential information stored in a.
func (a *OAuthCredsProps) GetCreds() (string, string) {
	return a.AccessToken, a.AccessTokenSecret
}

// OAuthApp provides features to save/load and set/get application credential
// information for OAuth.
type OAuthApp interface {
	utils.Savable
	utils.Loadable
	utils.Deletable
	SetCreds(ck, cs string)
	GetCreds() (string, string)
}

// TwitterOAuthAppProps is a credential information container for Twitter
// application.
type TwitterOAuthAppProps struct {
}

func newTwitterOAuthAppProps() *TwitterOAuthAppProps {
	return &TwitterOAuthAppProps{}
}

// OAuthAppProps abstracts features to set/get application credential
// properties for OAuth.
type OAuthAppProps interface {
	SetCreds(ck, cs string)
	GetCreds() (string, string)
}

// DefaultOAuthAppProps contains variables for general OAuth application usage.
type DefaultOAuthAppProps struct {
	ConsumerKey    string `json:"consumer_key" toml:"consumer_key" bson:"consumer_key" yaml:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret" toml:"consumer_secret" bson:"consumer_secret" yaml:"consumer_secret"`
}

// SetCreds sets application credential information to a.
func (a *DefaultOAuthAppProps) SetCreds(ck, cs string) {
	a.ConsumerKey = ck
	a.ConsumerSecret = cs
}

// GetCreds returns application credential information stored in a.
func (a *DefaultOAuthAppProps) GetCreds() (string, string) {
	return a.ConsumerKey, a.ConsumerSecret
}

// SetCreds sets application credential information to a.
func (a *TwitterOAuthAppProps) SetCreds(ck, cs string) {
	anaconda.SetConsumerKey(ck)
	anaconda.SetConsumerSecret(cs)
}

// GetCreds returns application credential information stored in a.
func (a *TwitterOAuthAppProps) GetCreds() (string, string) {
	return anaconda.GetConsumerKey(), anaconda.GetConsumerSecret()
}

// FileOAuthApp is OAuthApp associated with a specified file.
type FileOAuthApp struct {
	OAuthAppProps
	File string `json:"-" toml:"-" bson:"-" yaml:"-"`
}

// NewFileOAuthApp returns a new FileOAuthApp with file.
func NewFileOAuthApp(file string) (*FileOAuthApp, error) {
	a := &FileOAuthApp{&DefaultOAuthAppProps{}, file}
	err := a.Load()
	return a, utils.WithStack(err)
}

// NewFileTwitterOAuthApp returns a new FileTwitterOAuthApp with file.
func NewFileTwitterOAuthApp(file string) (*FileOAuthApp, error) {
	a := &FileOAuthApp{newTwitterOAuthAppProps(), file}
	err := a.Load()
	return a, utils.WithStack(err)
}

// Load retrieves credential information from a.File and stores it.
// This method does nothing and returns nil if the specified file doesn't exist.
func (a *FileOAuthApp) Load() error {
	tmp := &DefaultOAuthAppProps{}
	err := utils.DecodeFile(a.File, tmp)
	a.SetCreds(tmp.ConsumerKey, tmp.ConsumerSecret)
	return utils.WithStack(err)
}

// Save saves the credential information to a.File.
func (a *FileOAuthApp) Save() error {
	ck, cs := a.GetCreds()
	tmp := &DefaultOAuthAppProps{ck, cs}
	return utils.EncodeFile(a.File, tmp)
}

func (a *FileOAuthApp) Delete() error {
	return os.RemoveAll(a.File)
}

// DBOAuthApp is OAuthApp associated with a specified database.
type DBOAuthApp struct {
	OAuthAppProps
	col models.MgoCollection
}

// NewDBOAuthApp returns a new DBOAuthApp with a specified MongoDB collection.
// Currently only supported database is MongoDB.
func NewDBOAuthApp(col models.MgoCollection) (*DBOAuthApp, error) {
	a := &DBOAuthApp{&DefaultOAuthAppProps{}, col}
	err := a.Load()
	return a, utils.WithStack(err)
}

// NewDBTwitterOAuthApp returns a new DBTwitterOAuthApp with a specified
// MongoDB collection.
func NewDBTwitterOAuthApp(col models.MgoCollection) (*DBOAuthApp, error) {
	a := &DBOAuthApp{newTwitterOAuthAppProps(), col}
	err := a.Load()
	return a, utils.WithStack(err)
}

// Load retrieves credential information from a specified database and stores
// it.
func (a *DBOAuthApp) Load() error {
	tmp := &DefaultOAuthAppProps{}
	query := a.col.Find(nil)
	count, err := query.Count()
	if err != nil {
		return utils.WithStack(err)
	}
	if count > 0 {
		err = query.One(tmp)
		if err != nil {
			return utils.WithStack(err)
		}
	}
	a.SetCreds(tmp.ConsumerKey, tmp.ConsumerSecret)
	return nil
}

// Save saves the credential information to a specified database.
func (a *DBOAuthApp) Save() error {
	ck, cs := a.GetCreds()
	tmp := &DefaultOAuthAppProps{ck, cs}
	_, err := a.col.Upsert(nil, tmp)
	return utils.WithStack(err)
}

func (a *DBOAuthApp) Delete() error {
	_, err := a.col.RemoveAll(nil)
	return err
}
