package mybot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// OAuthCredentials contains values required for Twitter's user authentication.
type OAuthCredentials struct {
	ConsumerKey       string `json:"consumer_key" toml:"consumer_key"`
	ConsumerSecret    string `json:"consumer_secret" toml:"consumer_secret"`
	AccessToken       string `json:"access_token" toml:"access_token"`
	AccessTokenSecret string `json:"access_token_secret" toml:"access_token_secret"`
	File              string `json:"-" toml:"-"`
}

// Read does nothing and returns nil if the specified file doesn't exist.
func (a *OAuthCredentials) Read(file string) error {
	a.File = file
	ext := filepath.Ext(file)
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	switch ext {
	case ".json":
		err = json.Unmarshal(bs, a)
		if err != nil {
			return err
		}
	case ".toml":
		md, err := toml.Decode(string(bs), a)
		if err != nil {
			return err
		}
		if len(md.Undecoded()) != 0 {
			return fmt.Errorf("%v undecoded in %s", md.Undecoded(), file)
		}
	default:
		return fmt.Errorf("%s is invalid extension", ext)
	}
	return nil
}

func (a *OAuthCredentials) Write() error {
	ext := filepath.Ext(a.File)
	var bs []byte
	var err error
	switch ext {
	case ".json":
		bs, err = json.Marshal(a)
		if err != nil {
			return err
		}
	case ".toml":
		buf := new(bytes.Buffer)
		enc := toml.NewEncoder(buf)
		err = enc.Encode(a)
		if err != nil {
			return err
		}
		bs = buf.Bytes()
	default:
		return fmt.Errorf("%s is invalid extension", ext)
	}
	err = ioutil.WriteFile(a.File, bs, 0640)
	if err != nil {
		return err
	}
	return nil
}
