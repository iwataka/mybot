package models

import (
	"net/http"

	"github.com/markbates/goth"
)

type Authenticator interface {
	SetProvider(req *http.Request, name string)
	InitProvider(host, name string)
	CompleteUserAuth(provider string, w http.ResponseWriter, r *http.Request) (goth.User, error)
	Logout(provider string, w http.ResponseWriter, r *http.Request) error
}
