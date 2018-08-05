package oauth_test

import (
	. "github.com/iwataka/mybot/oauth"
	"github.com/iwataka/mybot/utils"
	"github.com/stretchr/testify/require"

	"os"
	"path/filepath"
	"testing"
)

func TestFileOAuth_SetGetCreds(t *testing.T) {
	a, err := NewFileOAuthCreds("")
	require.NoError(t, err)
	testOAuthSetGetCreds(t, a)
}

func TestDBOAuth_SetGetCreds(t *testing.T) {
	t.Skip("You must write mocking test for this")
	a, err := NewDBOAuthCreds(nil, "")
	require.NoError(t, err)
	testOAuthSetGetCreds(t, a)
}

func testOAuthSetGetCreds(t *testing.T, a OAuthCreds) {
	at := "foo"
	ats := "bar"
	a.SetCreds(at, ats)
	_at, _ats := a.GetCreds()
	require.Equal(t, at, _at)
	require.Equal(t, ats, _ats)
}

func TestFileOAuthApp_SetGetCreds(t *testing.T) {
	a, err := NewFileOAuthApp("")
	require.NoError(t, err)
	testOAuthAppSetGetCreds(t, a)
}

func TestFileTwitterOAuthApp_SetGetCreds(t *testing.T) {
	a, err := NewFileTwitterOAuthApp("")
	require.NoError(t, err)
	testOAuthAppSetGetCreds(t, a)
}

func TestDBTwitterOAuthApp_SetGetCreds(t *testing.T) {
	t.Skip("You must write mocking test for this")
	a, err := NewDBTwitterOAuthApp(nil)
	require.NoError(t, err)
	testOAuthAppSetGetCreds(t, a)
}

func testOAuthAppSetGetCreds(t *testing.T, a OAuthApp) {
	at := "foo"
	ats := "bar"
	a.SetCreds(at, ats)
	_at, _ats := a.GetCreds()
	require.Equal(t, at, _at)
	require.Equal(t, ats, _ats)
}

func TestFileOAuthCreds_Save(t *testing.T) {
	var err error
	dir := os.TempDir()
	file := filepath.Join(dir, "creds.json")
	defer os.Remove(file)

	a, err := NewFileOAuthCreds(file)
	require.NoError(t, err)
	testFileSave(t, a, file)
}

func TestFileOAuthApp_Save(t *testing.T) {
	var err error
	dir := os.TempDir()
	file := filepath.Join(dir, "creds.json")
	defer os.Remove(file)

	a, err := NewFileOAuthApp(file)
	require.NoError(t, err)
	testFileSave(t, a, file)
}

func testFileSave(t *testing.T, s utils.Savable, file string) {
	require.NoError(t, s.Save())
	info, _ := os.Stat(file)
	require.NotNil(t, info)
	require.False(t, info.IsDir())
}
