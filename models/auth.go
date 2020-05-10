package models

import (
	"net/http"

	"github.com/markbates/goth"
)

type Authenticator interface {
	SetProvider(name string, r *http.Request)
	InitProvider(provider, callback, consumerKey, consumerSecret string)
	CompleteUserAuth(provider string, w http.ResponseWriter, r *http.Request) (goth.User, error)
	Login(user goth.User, w http.ResponseWriter, r *http.Request) error
	GetLoginUser(r *http.Request) (goth.User, error)
	Logout(w http.ResponseWriter, r *http.Request) error
}
