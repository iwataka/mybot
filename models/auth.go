package models

import (
	"net/http"

	"github.com/markbates/goth"
)

type Authenticator interface {
	SetProvider(name string, r *http.Request)
	InitProvider(host, name, callback string)
	CompleteUserAuth(provider string, w http.ResponseWriter, r *http.Request) (goth.User, error)
	Logout(provider string, w http.ResponseWriter, r *http.Request) error
}
