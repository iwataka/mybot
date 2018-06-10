package oauth_test

import (
	. "github.com/iwataka/mybot/oauth"
	"github.com/iwataka/mybot/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFileOAuthSetGetCreds(t *testing.T) {
	a, err := NewFileOAuthCreds("")
	assert.NoError(t, err)
	testOAuthSetGetCreds(t, a)
}

func TestDBOAuthSetGetCreds(t *testing.T) {
	t.Skip("You must write mocking test for this")
	a, err := NewDBOAuthCreds(nil, "")
	assert.NoError(t, err)
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

func TestFileOAuthAppSetGetCreds(t *testing.T) {
	a, err := NewFileOAuthApp("")
	assert.NoError(t, err)
	testOAuthAppSetGetCreds(t, a)
}

func TestFileTwitterOAuthAppSetGetCreds(t *testing.T) {
	a, err := NewFileTwitterOAuthApp("")
	assert.NoError(t, err)
	testOAuthAppSetGetCreds(t, a)
}

func TestDBTwitterOAuthAppSetGetCreds(t *testing.T) {
	t.Skip("You must write mocking test for this")
	a, err := NewDBTwitterOAuthApp(nil)
	assert.NoError(t, err)
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

func TestFileOAuthCreds_Save(t *testing.T) {
	var err error
	dir := os.TempDir()
	file := filepath.Join(dir, "creds.json")
	defer os.Remove(file)

	a, err := NewFileOAuthCreds(file)
	assert.NoError(t, err)
	testFileSave(t, a, file)
}

func TestFileOAuthApp_Save(t *testing.T) {
	var err error
	dir := os.TempDir()
	file := filepath.Join(dir, "creds.json")
	defer os.Remove(file)

	a, err := NewFileOAuthApp(file)
	assert.NoError(t, err)
	testFileSave(t, a, file)
}

func testFileSave(t *testing.T, s utils.Savable, file string) {
	assert.NoError(t, s.Save())
	info, _ := os.Stat(file)
	require.NotNil(t, info)
	assert.False(t, info.IsDir())
}
