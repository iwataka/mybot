package mybot

import (
	"reflect"
	"testing"
)

func TestFileOAuthSetGetCreds(t *testing.T) {
	a, err := NewFileOAuthCreds("")
	if err != nil {
		t.Fatal(err)
	}
	testOAuthSetGetCreds(t, a)
}

func TestDBOAuthSetGetCreds(t *testing.T) {
	t.Skip("You must write mocking test for this")
	a, err := NewDBOAuthCreds(nil, "")
	if err != nil {
		t.Fatal(err)
	}
	testOAuthSetGetCreds(t, a)
}

func testOAuthSetGetCreds(t *testing.T, a OAuthCreds) {
	at := "foo"
	ats := "bar"
	a.SetCreds(at, ats)
	_at, _ats := a.GetCreds()
	if at != _at || ats != _ats {
		t.Fatalf("Inconsistent getter and setter of %s", reflect.TypeOf(a))
	}
}

func TestFileTwitterOAuthAppSetGetCreds(t *testing.T) {
	a, err := NewFileTwitterOAuthApp("")
	if err != nil {
		t.Fatal(err)
	}
	testOAuthAppSetGetCreds(t, a)
}

func TestDBTwitterOAuthAppSetGetCreds(t *testing.T) {
	t.Skip("You must write mocking test for this")
	a, err := NewDBTwitterOAuthApp(nil)
	if err != nil {
		t.Fatal(err)
	}
	testOAuthAppSetGetCreds(t, a)
}

func testOAuthAppSetGetCreds(t *testing.T, a OAuthApp) {
	at := "foo"
	ats := "bar"
	a.SetCreds(at, ats)
	_at, _ats := a.GetCreds()
	if at != _at || ats != _ats {
		t.Fatalf("Inconsistent getter and setter of %s", reflect.TypeOf(a))
	}
}
