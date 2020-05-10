package oauth_test

import (
	gomock "github.com/golang/mock/gomock"
	"github.com/iwataka/mybot/mocks"
	"github.com/iwataka/mybot/oauth"
	"github.com/iwataka/mybot/utils"
	"github.com/stretchr/testify/require"

	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFileOAuth_SetGetCreds(t *testing.T) {
	a, err := oauth.NewFileOAuthCreds("")
	require.NoError(t, err)
	testOAuthSetGetCreds(t, a)
}

func testOAuthSetGetCreds(t *testing.T, a oauth.OAuthAppProps) {
	at := "foo"
	ats := "bar"
	a.SetCreds(at, ats)
	_at, _ats := a.GetCreds()
	require.Equal(t, at, _at)
	require.Equal(t, ats, _ats)
}

func TestNewDBOAuthCreds(t *testing.T) {
	ctrl := gomock.NewController(t)
	col := mocks.NewMockMgoCollection(ctrl)
	query := mocks.NewMockMgoQuery(ctrl)
	col.EXPECT().Find(gomock.Any()).Return(query)
	query.EXPECT().Count().Return(0, nil)
	c, err := oauth.NewDBOAuthCreds(col, "foo")
	require.NoError(t, err)
	col.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil, nil)
	require.NoError(t, c.Save())
	col.EXPECT().RemoveAll(gomock.Any()).Return(nil, nil)
	require.NoError(t, c.Delete())
}

func TestNewDBOAuthCreds_withCountError(t *testing.T) {
	ctrl := gomock.NewController(t)
	col := mocks.NewMockMgoCollection(ctrl)
	query := mocks.NewMockMgoQuery(ctrl)
	col.EXPECT().Find(gomock.Any()).Return(query)
	query.EXPECT().Count().Return(0, errors.New("error"))
	_, err := oauth.NewDBOAuthCreds(col, "foo")
	require.Error(t, err)
}

func TestNewDBOAuthCreds_withExistingData(t *testing.T) {
	ctrl := gomock.NewController(t)
	col := mocks.NewMockMgoCollection(ctrl)
	query := mocks.NewMockMgoQuery(ctrl)
	col.EXPECT().Find(gomock.Any()).Return(query)
	query.EXPECT().Count().Return(1, nil)
	query.EXPECT().One(gomock.Any()).Return(nil)
	_, err := oauth.NewDBOAuthCreds(col, "foo")
	require.NoError(t, err)
}

func TestFileOAuthApp_SetGetCreds(t *testing.T) {
	a, err := oauth.NewFileOAuthApp("")
	require.NoError(t, err)
	testOAuthAppSetGetCreds(t, a)
}

func TestFileTwitterOAuthApp_SetGetCreds(t *testing.T) {
	a, err := oauth.NewFileTwitterOAuthApp("")
	require.NoError(t, err)
	testOAuthAppSetGetCreds(t, a)
}

func testOAuthAppSetGetCreds(t *testing.T, a oauth.OAuthApp) {
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
	a, err := oauth.NewFileOAuthCreds(file)
	require.NoError(t, err)
	defer func() { require.NoError(t, a.Delete()) }()
	testFileSave(t, a, file)
}

func TestFileOAuthApp_Save(t *testing.T) {
	var err error
	dir := os.TempDir()
	file := filepath.Join(dir, "creds.json")
	a, err := oauth.NewFileOAuthApp(file)
	require.NoError(t, err)
	defer func() { require.NoError(t, a.Delete()) }()
	testFileSave(t, a, file)
}

func testFileSave(t *testing.T, s utils.Savable, file string) {
	require.NoError(t, s.Save())
	info, _ := os.Stat(file)
	require.NotNil(t, info)
	require.False(t, info.IsDir())
}

func TestNewDBOAuthApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	col := mocks.NewMockMgoCollection(ctrl)
	query := mocks.NewMockMgoQuery(ctrl)
	col.EXPECT().Find(gomock.Any()).Return(query)
	query.EXPECT().Count().Return(0, nil)
	c, err := oauth.NewDBOAuthApp(col)
	require.NoError(t, err)
	col.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil, nil)
	require.NoError(t, c.Save())
	col.EXPECT().RemoveAll(gomock.Any()).Return(nil, nil)
	require.NoError(t, c.Delete())
}

func TestNewDBOAuthApp_withCountError(t *testing.T) {
	ctrl := gomock.NewController(t)
	col := mocks.NewMockMgoCollection(ctrl)
	query := mocks.NewMockMgoQuery(ctrl)
	col.EXPECT().Find(gomock.Any()).Return(query)
	query.EXPECT().Count().Return(0, errors.New("error"))
	_, err := oauth.NewDBOAuthApp(col)
	require.Error(t, err)
}

func TestNewDBOAuthApp_withExistingData(t *testing.T) {
	ctrl := gomock.NewController(t)
	col := mocks.NewMockMgoCollection(ctrl)
	query := mocks.NewMockMgoQuery(ctrl)
	col.EXPECT().Find(gomock.Any()).Return(query)
	query.EXPECT().Count().Return(1, nil)
	query.EXPECT().One(gomock.Any()).Return(errors.New("error"))
	_, err := oauth.NewDBOAuthApp(col)
	require.Error(t, err)
}
