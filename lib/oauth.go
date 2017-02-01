package mybot

// OAuthCredentials contains values required for Twitter's user authentication.
type OAuthCredentials struct {
	AccessToken       string `json:"access_token" toml:"access_token"`
	AccessTokenSecret string `json:"access_token_secret" toml:"access_token_secret"`
	File              string `json:"-" toml:"-"`
}

type OAuthApp struct {
	ConsumerKey    string `json:"consumer_key" toml:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret" toml:"consumer_secret"`
	File           string `json:"-" toml:"-"`
}

// Decode does nothing and returns nil if the specified file doesn't exist.
func (a *OAuthCredentials) Decode(file string) error {
	a.File = file
	return DecodeFile(file, a)
}

func (a *OAuthCredentials) Encode() error {
	return EncodeFile(a.File, a)
}

// Decode does nothing and returns nil if the specified file doesn't exist.
func (a *OAuthApp) Decode(file string) error {
	a.File = file
	return DecodeFile(file, a)
}

func (a *OAuthApp) Encode() error {
	return EncodeFile(a.File, a)
}
